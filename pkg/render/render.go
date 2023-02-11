package render

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/imrcht/bed-n-breakfast/pkg/config"
	"github.com/imrcht/bed-n-breakfast/pkg/models"
)

var App *config.AppConfig

func SetApp(a *config.AppConfig) {
	App = a
}

// DefaultTempData
func defaultTempData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderHtml(w http.ResponseWriter, temp string, td *models.TemplateData) {

	var tc map[string]*template.Template

	if App.UseCache {
		// get the template cache from app config
		tc = App.TemplateCache
	} else {

		// Create template cache and store
		tc, _ = CreateTemplateCache()
		// if err != nil {
		// 	log.Fatal(err)
		// }
	}

	// get requested template from cache
	t, ok := tc[temp]

	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = defaultTempData(td)

	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}

	// render
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}

	// parsedFile, _ := template.ParseFiles("./templates/"+temp, "./templates/base.layout.tmpl")

	// err := parsedFile.Execute(w, nil)

	// if err != nil {
	// 	fmt.Fprint(w, "Error parsing teamplate", err)
	// 	return
	// }
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all pages of templates starting with *.page.tmpl;
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
