package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/VtG242/boss/internal/models"
	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	// Call the newTemplateData() helper to get a templateData struct containing the 'default' data
	data := app.newTemplateData(r)

	// Use the new render helper.
	app.render(w, http.StatusOK, "home.html", data)

}

func (app *application) playersHelp(w http.ResponseWriter, r *http.Request) {
	// Call the newTemplateData() helper to get a templateData struct containing the 'default' data
	data := app.newTemplateData(r)

	// Use the new render helper.
	app.render(w, http.StatusOK, "players-help.html", data)

}

func (app *application) playersView(w http.ResponseWriter, r *http.Request) {
	
	//time.Sleep(30 * time.Second)
	players, err := app.players.All()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Call the newTemplateData() helper to get a templateData struct containing the 'default' data
	data := app.newTemplateData(r)
	data.Players = players

	// Use the new render helper.
	app.render(w, http.StatusOK, "players.html", data)
}

// Add player handler function.
func (app *application) playerView(w http.ResponseWriter, r *http.Request) {

	// When httprouter is parsing a request, the values of any named parameters will be stored in the request context
	params := httprouter.ParamsFromContext(r.Context())

	// get the value of pid
	id, err := strconv.Atoi(params.ByName("pid"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusNotFound)
		return
	}

	// Use the BossModel object's Get method to retrieve the data for player id
	// If no matching record is found, return a 404 Not Found response.
	player, err := app.players.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Call the newTemplateData() helper to get a templateData struct containing the 'default' data
	data := app.newTemplateData(r)
	data.Player = player

	// Use the new render helper.
	app.render(w, http.StatusOK, "player.html", data)

}

// show player form for create/edit action
func (app *application) playerForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new player..."))
}

// Add player Create handler function.
func (app *application) playerCreate(w http.ResponseWriter, r *http.Request) {

	surname := "Testovic"
	firstname := "Test"
	sex := "M"
	birthdate, _ := time.Parse("2006-01-02", "1981-02-01")
	town := "Mesto"
	country := "CZ"
	nickname := ""
	hash := "abcdef123"

	// Pass the data to the Players.Insert() method, receiving the ID of the new record back.
	id, err := app.players.Insert(surname, firstname, sex, birthdate, town, country, nickname, hash)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect to page displaying newly created player
	http.Redirect(w, r, fmt.Sprintf("/player/view/%d", id), http.StatusSeeOther)

}

func (app *application) tournaments(w http.ResponseWriter, r *http.Request) {
	// Call the newTemplateData() helper to get a templateData struct containing the 'default' data
	data := app.newTemplateData(r)

	// Use the new render helper.
	app.render(w, http.StatusOK, "tournaments.html", data)

}
