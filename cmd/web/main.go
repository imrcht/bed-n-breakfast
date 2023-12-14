package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/imrcht/bed-n-breakfast/internals/config"
	"github.com/imrcht/bed-n-breakfast/internals/driver"
	"github.com/imrcht/bed-n-breakfast/internals/handlers"
	"github.com/imrcht/bed-n-breakfast/internals/helpers"
	"github.com/imrcht/bed-n-breakfast/internals/models"
	"github.com/imrcht/bed-n-breakfast/internals/render"
)

const portNumber = ":8000"

// @start command for multiple go files: go run cmd/web/*.go

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	log.Printf("Server listening on port: %s", portNumber)

	// We can use http.Server instance to run
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()

	log.Fatal(err)

	// Or we can use directly ListenAndServer function to listen
	// _ = http.ListenAndServe(portNumber, routes(&app))
}

func run() (*driver.DB, error) {
	// * Adding custom var type to session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// * Change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// * Connect to database
	dsn := `host=localhost port=5432 dbname=bookings user=rachitgupta password=`
	db, err := driver.ConnectSql(dsn)

	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
		return db, err
	}

	log.Println("Connected to database!")

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal(err)
		return db, err
	}

	app.TemplateCache = tc
	app.UseCache = false
	app.InfoLog = infoLog
	app.ErrorLog = errorLog

	repo := handlers.NewHandler(&app, db)
	handlers.NewRepo(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
