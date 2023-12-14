package render

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/imrcht/bed-n-breakfast/internals/config"
	"github.com/imrcht/bed-n-breakfast/internals/helpers"
	"github.com/imrcht/bed-n-breakfast/internals/models"
	"github.com/justinas/nosurf"
)

var App *config.AppConfig

func NewRenderer(a *config.AppConfig) {
	App = a
}

// DefaultTempData
func defaultTempData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = App.Session.PopString(r.Context(), "flash")
	td.Error = App.Session.PopString(r.Context(), "error")
	td.Warning = App.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

func Template(w http.ResponseWriter, r *http.Request, temp string, td *models.TemplateData) {

	var tc map[string]*template.Template

	if App.UseCache {
		// * Get the template cache from app config
		tc = App.TemplateCache
	} else {

		// * Create template cache and store
		newTc, errInCreatingTc := CreateTemplateCache()
		if errInCreatingTc != nil {
			helpers.ServerError(w, errInCreatingTc)
			return
		}
		tc = newTc
	}

	// * Get requested template from cache
	t, ok := tc[temp]

	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = defaultTempData(td, r)

	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// * Render template to response writer
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

	// * Old approach
	/* parsedFile, _ := template.ParseFiles("./templates/"+temp, "./templates/base.layout.tmpl")
	err := parsedFile.Execute(w, nil)

	if err != nil {
		fmt.Fprint(w, "Error parsing template", err)
		return
	}
	*/
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// * Get all pages of templates starting with *.page.tmpl;
	pages, err := filepath.Glob("./templates/*.page.tmpl")

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		t := filepath.Base(page)

		ts, err := template.New(t).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[t] = ts
	}

	return myCache, nil
}

// var tc = make(map[string]*template.Template)

// func RenderHtml(w http.ResponseWriter, t string) {
// 	var tmpl *template.Template
// 	var err error

// 	_, inMap := tc[t]

// 	if !inMap {
// 		// create a template
// 		log.Println("Creating template and adding to cache")
// 		err = createTemplate(t)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	} else {
// 		log.Println("Using cached tc")
// 	}

// 	tmpl = tc[t]

// 	err = tmpl.Execute(w, nil)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

// func createTemplate(t string) error {
// 	templates := []string{
// 		fmt.Sprintf("./templates/%s", t),
// 		"./templates/base.layout.tmpl",
// 	}

// 	tmpl, err := template.ParseFiles(templates...)

// 	if err != nil {
// 		return err
// 	}

// 	tc[t] = tmpl

// 	return nil
// }
