package models

import "github.com/imrcht/bed-n-breakfast/internals/forms"

// TemplateData holds data set to templates from handlers
type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	BoolMap         map[string]bool
	Data            map[string]interface{}
	CSRFToken       string //Cross site request forgery token
	Flash           string
	Warning         string
	Error           string
	Form            *forms.Form
	IsAuthenticated bool
	UserAccessLevel int
}
