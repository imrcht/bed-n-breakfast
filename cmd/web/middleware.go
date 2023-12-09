package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

// * WriteToConsole writes the request in every hit
func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("************** Hitting the url: ", r.URL, "**************")
		next.ServeHTTP(w, r)
	})
}

// * NoSurf adds CSRF protection to every post request
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

// * SessionLoad loads and save the session data for the current request and communicates the session to and from the client in a cookie.
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
