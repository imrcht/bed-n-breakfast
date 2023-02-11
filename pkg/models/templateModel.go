package models

// TemplateData holds data set to templates from handlers
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	BoolMap   map[string]bool
	Data      map[string]interface{}
	CSRFToken string //Cross site request forgery token
	Flash     string
	Warning   string
	Error     string
}
