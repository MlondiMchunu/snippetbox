package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snippetbox.mlodev.net/internal/models"
	"snippetbox.mlodev.net/ui"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	//initialize a new map to act as a cache
	cache := map[string]*template.Template{}

	//gets a slice of all filepaths that match the pattern ./ui/html/pages/*.tmpl
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	//loop through the page filepaths one-by-one
	for _, page := range pages {

		//extract the filename from the full filepath
		name := filepath.Base(page)
		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		//create a slice containing the filepaths for our base template
		/*files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			page,
		}*/

		// Parse the base template file into a template set.
		/*ts, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}
		*/

		//add the template set to the mao as normal
		cache[name] = ts
	}
	//Return the map
	return cache, nil

}
