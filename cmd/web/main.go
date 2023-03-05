package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/imrcht/bed-n-breakfast/internals/config"
	"github.com/imrcht/bed-n-breakfast/internals/handlers"
	"github.com/imrcht/bed-n-breakfast/internals/models"
	"github.com/imrcht/bed-n-breakfast/internals/render"
)

const portNumber = ":7000"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	// Adding custom var type to session
	gob.Register(models.Reservation{})

	// change this to true when in production
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
	}

	app.TemplateCache = tc
	app.UseCache = false

	render.SetApp(&app)

	repo := handlers.NewHandler(app)
	handlers.NewRepo(repo)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)

	log.Println("Server listening on port: 7000")

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
