package handlers

import (
	"net/http"

	"github.com/imrcht/bed-n-breakfast/pkg/config"
	"github.com/imrcht/bed-n-breakfast/pkg/models"
	"github.com/imrcht/bed-n-breakfast/pkg/render"
)

// Repository
type Repository struct {
	App config.AppConfig
}

// Repo
var Repo *Repository

// NewHandler
func NewHandler(a config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// Create NewRepo
func NewRepo(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	// Fetching remote address from request and storing it in session
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.RenderHtml(w, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	mapString := map[string]string{}
	mapString["test"] = "Hello again"

	// Fetching remoteIp from session which was stored in home page , return -> string (empty string if there's no value present for that key)
	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")

	mapString["remote_ip"] = remoteIp
	render.RenderHtml(w, "about.page.tmpl", &models.TemplateData{
		StringMap: mapString,
	})
}
