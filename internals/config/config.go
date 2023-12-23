package config

// * Config is imported by other parts of the application but it should not import anything else from the application itself otherwise it may create import cycle problem

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
	"github.com/imrcht/bed-n-breakfast/internals/models"
)

// * AppConfig holds the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
