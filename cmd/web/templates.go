package main

import (
	"github.com/evillgenius75/glennjokebox/internal/models"
	"github.com/evillgenius75/glennjokebox/ui"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"
)

type templateData struct {
	CurrentYear     int
	Joke            *models.Joke
	Jokes           []*models.Joke
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		//files := []string{
		//	"./ui/html/base.tmpl",
		//	"./ui/html/partials/nav.tmpl",
		//	page,
		//}

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
