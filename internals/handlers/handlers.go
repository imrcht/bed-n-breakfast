package handlers

import (
	"encoding/json"
	"errors"
	"html"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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
	// * Fetch the room id of generals room page from database and pass it to template data as a string map
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

type jsonResponse struct {
	Ok        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomId    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON: handles request for availability and send response in json format
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	isAvailable, err := m.DB.SearchAvailabilityByDatesByRoomId(startDate, endDate, roomId)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	respStruct := jsonResponse{
		Ok:        isAvailable,
		Message:   " ",
		RoomId:    strconv.Itoa(roomId),
		StartDate: sd,
		EndDate:   ed,
	}

	respJson, err := json.MarshalIndent(respStruct, "", "  ")

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJson)
}

// PostAvailability: handles request for availability and redirects to choose-room page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	// * 2001-01-01 --> 01/02 03:04:05PM '06 -0700
	// * Layout is the format of the date
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	availableRooms, err := m.DB.SearchAvailabilityForAllRoomsByDates(startDate, endDate)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(availableRooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No available rooms")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = availableRooms

	// * Store reservation details in session so that it can be used in next page
	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// SelectRoom: takes room id as url param, stores it in session and redirects to make-reservation page
func (m *Repository) SelectRoom(w http.ResponseWriter, r *http.Request) {
	roomId, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// * Get reservation details from session
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		m.App.ErrorLog.Println("Cannot get reservation from session")
		m.App.Session.Put(r.Context(), "error", "Cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.RoomId = roomId
	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// BookRoom: takes room id, start date and end date as query params, stores it in session and redirects to make-reservation page
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// * Fetch room details from DB
	room, err := m.DB.GetRoomById(roomId)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
		RoomId:    roomId,
		Room:      room,
	}
	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// Reservation: takes reservations details from session and renders make-reservation page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	// sendForm := forms.Form{}
	// * Get reservation details from session
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		m.App.ErrorLog.Println("Cannot get reservation from session")
		helpers.ServerError(w, errors.New("cannot get reservation from session"))
		return
	}

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")
	roomId := res.RoomId

	// * Fetch room details from DB
	room, err := m.DB.GetRoomById(roomId)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room = room
	// * Store reservation details in session so that it can be used in next page
	m.App.Session.Put(r.Context(), "reservation", res)

	strmap := make(map[string]string)
	strmap["start_date"] = sd
	strmap["end_date"] = ed
	strmap["room_name"] = room.RoomName

	// * interface is a type that can hold any type of data where `{}` represents empty data for any type
	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: strmap,
	})
}

// PostReservation: takes reservation details from form, checks whether room is available or not, validates form, add the record to reservation and room restriction table and redirects to reservation-summary page
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Cannot get reservation from session")
		m.App.Session.Put(r.Context(), "error", "cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// sd := r.Form.Get("start_date")

	// // * 2001-01-01 --> 01/02 03:04:05PM '06 -0700
	// // * Layout is the format of the date
	// layout := "2006-01-02"
	// startDate, err := time.Parse(layout, sd)

	// * Check if room is available
	roomAvailabilityStatus, err := m.DB.SearchAvailabilityByDatesByRoomId(reservation.StartDate, reservation.EndDate, reservation.RoomId)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if !roomAvailabilityStatus {
		m.App.Session.Put(r.Context(), "error", "Selected room is not available")
		http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

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

	reservationId, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	reservation.ID = reservationId

	roomRestriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		ReservationId: reservationId,
		RoomID:        reservation.RoomId,
		RestrictionID: 3,
	}

	err = m.DB.InsertRoomRestriction(roomRestriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// * Send notification to guest
	htmlMsg := `
		<strong>Reservation Confirmation</strong><br>
		Dear ` + html.EscapeString(reservation.FirstName) + `, <br>
		This is to confirm your reservation from ` + reservation.StartDate.Format("2006-01-02") + ` to ` + reservation.EndDate.Format("2006-01-02") + `.
	`
	msg := models.MailData{
		To:      reservation.Email,
		From:    "Bed N'Breakfast <no-reply@bnb.com>",
		Subject: "Reservation Confirmation",
		Content: htmlMsg,
	}
	m.App.MailChan <- msg

	// * Send notification to property owner
	htmlMsg = `
		<strong>Reservation Notification</strong><br>
		A reservation has been made for ` + reservation.FirstName + ` ` + reservation.LastName + ` from ` + reservation.StartDate.Format("2006-01-02") + ` to ` + reservation.EndDate.Format("2006-01-02") + `.
	`
	msg = models.MailData{
		To:      "propertyowner@bnb.com",
		From:    "Bed N'Breakfast <no-reply@bnb.com>",
		Subject: "Reservation Notification",
		Content: htmlMsg,
	}
	m.App.MailChan <- msg

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// ReservationSummary: takes reservation details from session, clear the session and renders reservation-summary page
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

	stringMap := make(map[string]string)
	stringMap["start_date"] = reservation.StartDate.Format("2006-01-02")
	stringMap["end_date"] = reservation.EndDate.Format("2006-01-02")

	data["reservation"] = reservation
	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}
