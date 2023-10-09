package main

import (
	"html/template"
	"path/filepath"
	"social-network/internal/models"
	"social-network/pkg/forms"
	"time"
)

type templateData struct {
	CurrentYear int
	Flash       string
	Form        *forms.Form
	User        *models.User
	Users       []*models.User
}

// humanDate returns a nicely formatted human-readable string representation of time.Time.
func humanDate(t time.Time) string {
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

func seq(start, end int) []int {
	var res []int
	for i := start; i <= end; i++ {
		res = append(res, i)
	}

	return res
}

var functions = template.FuncMap{
	"humanDate": humanDate,
	"seq":       seq,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
