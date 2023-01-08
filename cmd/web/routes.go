package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// The method returns http.Handler - middleware is used at the end
func (app *application) routes() http.Handler {
	// Initialize the router.
	router := httprouter.New()

	// Create a handler function which wraps our client error function
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.clientError(w, http.StatusNotFound)
	})

	// Create a file server which serves files from "./ui/static" directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// For matching paths, we strip the "/static" prefix before the request reaches the file server
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// routes using the appropriate methods, patterns and handlers
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/players", app.playersView)
	router.HandlerFunc(http.MethodGet, "/players/help", app.playersHelp)
	router.HandlerFunc(http.MethodGet, "/player/view/:pid", app.playerView)
	router.HandlerFunc(http.MethodGet, "/player/create", app.playerForm)
	router.HandlerFunc(http.MethodPost, "/player/create", app.playerCreate)
	router.HandlerFunc(http.MethodGet, "/tournaments", app.tournaments)

	// Create a middleware chain containing our 'standard' middleware
	// which will be used for every request our application receives.
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// using middleware chain: recoverPanic -> logRequest -> secureHeaders -> servemux -> app handler
	//return app.recoverPanic(app.logRequest(secureHeaders(mux)))
	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(router)
}
