package main

import (
	"html/template"

	"snippetbox.mlodev.net/internal/models"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	//initialize a new map to act as a cache
	cache := map[string]*template.Template{}

}
