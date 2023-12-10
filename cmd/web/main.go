package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/imrcht/bed-n-breakfast/internals/config"
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

	err := run()
	if err != nil {
		log.Fatal(err)
	}

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

func run() error {
	// * Adding custom var type to session
	gob.Register(models.Reservation{})

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

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal(err)
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false
	app.InfoLog = infoLog
	app.ErrorLog = errorLog

	render.SetApp(&app)

	repo := handlers.NewHandler(app)
	handlers.NewRepo(repo)

	helpers.NewHelpers(&app)

	return nil
}
