package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/imrcht/bed-n-breakfast/internals/config"
)

var app *config.AppConfig

// NewHelpers sets up config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.ErrorLog.Println("Client error with status of", http.StatusText(status))
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
