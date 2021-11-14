package main

import (
	"html/template"
	"path/filepath"
	"time"

	"kerseeeHuang.com/snippetbox/pkg/forms"
	"kerseeeHuang.com/snippetbox/pkg/models"
)

// templateData store snippets that we want to render with html templates.
type templateData struct {
	CurrentYear		int
	Flash 			string
	Form 			*forms.Form
	IsAuthenticated	bool
	Snippet 		*models.Snippet
	Snippets 		[]*models.Snippet
}

// humanDate return a nicely formatted string of time.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// functions store the custom functions used in templates.
// Template functions should only return one value, or one value and an error.
var functions = template.FuncMap {
	"humanDate": humanDate,
}

// newTemplateCache create the cache of tamplates with pages in given dir.
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize cache.
	cache := map[string]*template.Template{}

	// Get all files with suffix ".page.tmpl"
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Parse each page and store their cache.
	for _, page := range pages {
		// Extract the file name.
		name := filepath.Base(page)
		
		// New a set of templates, register the custom functions used in templates, 
		// and parse the files to the templates.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add layout pages into template set.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add partial pages into template set.
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add this template set to caches.
		cache[name] = ts
	}
	
	return cache, nil
}