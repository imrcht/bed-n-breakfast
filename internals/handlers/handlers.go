package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/imrcht/bed-n-breakfast/internals/config"
	"github.com/imrcht/bed-n-breakfast/internals/driver"
	"github.com/imrcht/bed-n-breakfast/internals/forms"
	"github.com/imrcht/bed-n-breakfast/internals/helpers"
	"github.com/imrcht/bed-n-breakfast/internals/models"
	"github.com/imrcht/bed-n-breakfast/internals/render"
	"github.com/imrcht/bed-n-breakfast/internals/repository"
	"github.com/imrcht/bed-n-breakfast/internals/repository/dbrepo"
)

// Repository
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Repo
var Repo *Repository

// NewHandler
func NewHandler(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgressDBRepo(a, db.SQL),
	}
}

// Create NewRepo
func NewRepo(r *Repository) {
	Repo = r
}

// * Receivers: Structs with functions
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	// * Fetching remote address from request and storing it in session
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// * `{}` represents empty data for any type and can be used to prepopulate data
	mapString := map[string]string{}
	mapString["test"] = "Hello again"

	// * Fetching remoteIp from session which was stored in home page , return -> string (empty string if there's no value present for that key)
	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")

	mapString["remote_ip"] = remoteIp
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: mapString,
	})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

// * AvailabilityJSON handles request for availability and send response in json format
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	respStruct := jsonResponse{
		Ok:      true,
		Message: "Availability in json",
	}

	respJson, err := json.MarshalIndent(respStruct, "", "  ")

	if err != nil {
		helpers.ServerError(w, err)
		return
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
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	// sendForm := forms.Form{}
	emptyReservation := models.Reservation{}

	// * interface is a type that can hold any type of data where `{}` represents empty data for any type
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
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
	form.IsValidPhone("phone", r)

	if !form.Valid() {
		data := make(map[string]interface{})

		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		m.App.ErrorLog.Println("Reservation not found in session")
		m.App.Session.Put(r.Context(), "flash", "Reservation not found ")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})

	data["reservation"] = reservation
	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
