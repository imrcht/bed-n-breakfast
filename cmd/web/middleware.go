package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// WriteToConsole writes the request in every hit
// func WriteToConsole(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("Hitting the page")
// 		next.ServeHTTP(w, r)
// 	})
// }

// NoSurf adds CSRF protection to every post request
func NoSurf(next http.Handler) http.Handler {
	csrfHanlder := nosurf.New(next)

	csrfHanlder.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Secure:   app.InProduction,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHanlder
}

// SessionLoad loads and save the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
