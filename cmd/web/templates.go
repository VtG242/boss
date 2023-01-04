package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/VtG242/boss/internal/models"
)

// Struct for any dynamic data passed to templates.
type templateData struct {
	CurrentYear int
	Player      *models.Player
	Players     []*models.Player
}

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func humanDate(t time.Time) string {
	return t.Format("02.03. 2006")
}

// Initialize a template.FuncMap object and store it in a global variable. This is
// essentially a string-keyed map which acts as a lookup between the names of our
// custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// path and pattern for pages templates
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	// create teplate set for each page and add it to cache
	for _, page := range pages {

		// name of template in cache map
		name := filepath.Base(page)

		// The template.FuncMap must be registered with the template set before you
		// call the ParseFiles() method. This means we have to use template.New() to
		// create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse base layout template.
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/layout.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() on previously created ts to add any parts templates like menu
		//ts, err = ts.ParseGlob("./ui/html/parts/*.html")
		//if err != nil {
		//	return nil, err
		//}

		// Finally add the page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map, using the name of the page (like 'home.html') as the key.
		cache[name] = ts
	}

	// Return the map.
	return cache, nil
}
