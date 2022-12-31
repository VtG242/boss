package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// The serverError helper writes an error message and stack trace to the log,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.log.Error(trace)
	http.Error(w, "BOSS: "+http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
// to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, "BOSS: "+http.StatusText(status), status)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// Retrieve the appropriate template set from the cache based on the page name (like 'home.html').
	// If no entry exists in the cache, create a new error and call the serverError() helper
	ts, ok := app.templates[page]
	if !ok {
		err := fmt.Errorf("template %s doesn't exist", page)
		app.serverError(w, err)
		return
	}

	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	// Write the template to the buffer instead of http.ResponseWriter
	// if there is any error we call the the serverError() helper
	err := ts.ExecuteTemplate(buf, "layout", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// If the template is written to the buffer without any errors, we are safe
	// to go ahead and write the HTTP status code to http.ResponseWriter
	w.WriteHeader(status)

	// Write content of buffer to the http.ResponseWriter
	buf.WriteTo(w)
}

// Create an newTemplateData() helper, which returns a pointer to a templateData
// struct initialized with the current year.
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}
