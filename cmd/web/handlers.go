package main

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"social-network/internal/models"
	"social-network/pkg/forms"
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

	app.render(w, r, "show.page.tmpl", &templateData{
		User: u,
	})
}

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form := forms.NewForm(r.PostForm)
	form.Required("first_name", "last_name", "gender", "age", "city", "interests")
	form.MaxLength("first_name", 50)
	form.MaxLength("last_name", 50)
	form.MaxLength("city", 50)
	form.MaxLength("interests", 500)
	form.PermittedValues("gender", "male", "female")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})

		return
	}

	ageInput := form.Get("age")
	age, _ := strconv.Atoi(ageInput)

	id, err := app.users.Insert(
		form.Get("first_name"),
		form.Get("last_nae"),
		form.Get("interests"),
		form.Get("city"),
		models.Gender(form.Get("gender")),
		uint32(age),
	)
	if err != nil {
		app.serverError(w, err)

		return
	}

	app.session.Put(r, "flash", "User successfully registered!")

	http.Redirect(w, r, fmt.Sprintf("/user/%d", id), http.StatusSeeOther)
}

func (app *application) registerUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.NewForm(nil),
	})
}
