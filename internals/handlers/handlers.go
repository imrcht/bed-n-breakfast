package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/imrcht/bed-n-breakfast/internals/config"
	"github.com/imrcht/bed-n-breakfast/internals/forms"
	"github.com/imrcht/bed-n-breakfast/internals/models"
	"github.com/imrcht/bed-n-breakfast/internals/render"
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

	render.RenderHtml(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	mapString := map[string]string{}
	mapString["test"] = "Hello again"

	// Fetching remoteIp from session which was stored in home page , return -> string (empty string if there's no value present for that key)
	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")

	mapString["remote_ip"] = remoteIp
	render.RenderHtml(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: mapString,
	})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderHtml(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderHtml(w, r, "majors.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderHtml(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	respStruct := jsonResponse{
		Ok:      true,
		Message: "Availability in json",
	}

	respJson, err := json.MarshalIndent(respStruct, "", "  ")

	if err != nil {
		log.Println("Error: ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("Start date is: %s, End date is %s: ", start, end)))
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderHtml(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	// sendForm := forms.Form{}
	emptyReservation := models.Reservation{}

	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderHtml(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Println("Error while parsing form data: ", err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	// JOI validations in GO
	// Method 1
	// fields := []string{"first_name", "last_name", "email", "phone"}
	// form.HasMany(fields, r)
	// Method 2
	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 4, r)
	form.IsValidEmail("email", r)

	if !form.Valid() {
		data := make(map[string]interface{})

		data["reservation"] = reservation

		render.RenderHtml(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {

}
