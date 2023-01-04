package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/VtG242/boss/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.clientError(w, http.StatusNotFound)
		return
	}

	//panic("oops! something went wrong") // Deliberate panic

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

func (app *application) players(w http.ResponseWriter, r *http.Request) {
	players, err := app.db.Latest()
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
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.clientError(w, http.StatusNotFound)
		return
	}

	// Use the BossModel object's Get method to retrieve the data for player id
	// If no matching record is found, return a 404 Not Found response.
	player, err := app.db.Get(id)
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

// Add player Create handler function.
func (app *application) playerCreate(w http.ResponseWriter, r *http.Request) {
	// Use r.Method to check whether the request is using POST or not.
	if r.Method != "POST" {
		// suggest allowed method
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	surname := "Testovic"
	firstname := "Test"
	sex := "M"
	birthdate, _ := time.Parse("2006-01-02", "1981-02-01")
	town := "Mesto"
	country := "CZ"
	nickname := ""
	hash := "abcdef123"

	// Pass the data to the Players.Insert() method, receiving the ID of the new record back.
	id, err := app.db.Insert(surname, firstname, sex, birthdate, town, country, nickname, hash)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the player.
	http.Redirect(w, r, fmt.Sprintf("/player/view?id=%d", id), http.StatusSeeOther)

}

func (app *application) tournaments(w http.ResponseWriter, r *http.Request) {
	// Call the newTemplateData() helper to get a templateData struct containing the 'default' data
	data := app.newTemplateData(r)

	// Use the new render helper.
	app.render(w, http.StatusOK, "tournaments.html", data)

}
