package main

import (
	"errors"
	"fmt"
	"net/http"
	"social-network/internal/models"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)

		return
	}

	users, err := app.users.Latest()
	if err != nil {
		app.serverError(w, err)

		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{Users: users})
}

func (app *application) showUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)

		return
	}

	u, err := app.users.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{User: u})
}

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", 405)

		return
	}

	// Dummy data
	firstName := "Alice"
	lastName := "Doe"
	interests := "sports, reading, science"
	city := "Moscow"
	gender := models.GenderFemale
	age := uint32(28)

	id, err := app.users.Insert(firstName, lastName, interests, city, gender, age)
	if err != nil {
		app.serverError(w, err)

		return
	}

	http.Redirect(w, r, fmt.Sprintf("/user?id=%d", id), http.StatusSeeOther)
}
