package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/imrcht/bed-n-breakfast/internals/config"
	"github.com/imrcht/bed-n-breakfast/internals/handlers"
)

// * Routes returns mux router
func routes(app *config.AppConfig) http.Handler {
	// * Using pat for routing
	/*mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))*/

	// * Using chi for routing
	mux := chi.NewMux()

	// * Using recoverer middleware comes along with chi
	mux.Use(middleware.Recoverer)
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)

	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Get("/choose-room/{id}", handlers.Repo.SelectRoom)
	mux.Get("/book-room", handlers.Repo.BookRoom)

	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Post("/user/login", handlers.Repo.PostShowLogin)
	mux.Post("/user/register", handlers.Repo.PostSignUpJson)
	mux.Get("/user/logout", handlers.Repo.Logout)

	mux.Route("/admin", func(r chi.Router) {
		// * Using middleware to check if user is authenticated
		r.Use(Auth)
		r.Get("/dashboard", handlers.Repo.AdminDashboard)
	})

	// Using static folder
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
