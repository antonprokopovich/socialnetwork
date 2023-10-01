package main

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"social-network/internal/models"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	users, err := app.users.Latest()
	if err != nil {
		app.serverError(w, err)

		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{Users: users})
}

func (app *application) showUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
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
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	firstName := r.PostForm.Get("first_name")
	lastName := r.PostForm.Get("last_name")
	interests := r.PostForm.Get("interests")
	gender := r.PostForm.Get("gender")
	city := r.PostForm.Get("city")
	ageVal := r.PostForm.Get("age")

	age, err := strconv.Atoi(ageVal)
	if err != nil {
		err = errors.Wrap(err, "gender value should be of int type")

		app.serverError(w, err)

		return
	}

	id, err := app.users.Insert(firstName, lastName, interests, city, models.Gender(gender), uint32(age))
	if err != nil {
		app.serverError(w, err)

		return
	}

	http.Redirect(w, r, fmt.Sprintf("/user/%d", id), http.StatusSeeOther)
}

func (app *application) registerUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", nil)
}
