package main

import (
	"html/template"
	"path/filepath"

	"snippetbox.mlodev.net/internal/models"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	//initialize a new map to act as a cache
	cache := map[string]*template.Template{}

	//gets a slice of all filepaths that match the pattern ./ui/html/pages/*.tmpl
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	//loop through the page filepaths one-by-one
	for _, page := range pages {

		//extract the filename from the full filepath
		name := filepath.Base(page)

		//create a slice containing the filepaths for our base template
		files := []string{
			"./ui/html/base.tmpl",
			".ui/html/partials/nav.tmpl",
			page,
		}

		//parse the files into a template set
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	//Return the map
	return cache, nil

}
