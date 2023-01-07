package main

import "net/http"

// The method returns http.Handler - middleware is used at the end
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Create a file server which serves files from "./ui/static" directory
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Register the two new handler functions and corresponding URL patterns with
	// the servemux, in exactly the same way that we did before.
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/players", app.playersView)
	mux.HandleFunc("/players/help", app.playersHelp)
	mux.HandleFunc("/player/view", app.playerView)
	mux.HandleFunc("/player/create", app.playerCreate)
	mux.HandleFunc("/tournaments", app.tournaments)

	// using middleware chain: recoverPanic -> logRequest -> secureHeaders -> servemux -> app handler
	return app.recoverPanic(app.logRequest(secureHeaders(mux)))
}
