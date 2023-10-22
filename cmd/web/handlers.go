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
	form.Required("first_name", "last_name", "gender", "age", "email", "password")
	form.MaxLength("first_name", 50)
	form.MaxLength("last_name", 50)
	form.MaxLength("city", 50)
	form.MaxLength("interests", 500)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	form.PermittedValues("gender", "male", "female")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})

		return
	}

	ageInput := form.Get("age")
	age, _ := strconv.Atoi(ageInput)

	_, err := app.users.Insert(
		form.Get("first_name"),
		form.Get("last_nae"),
		form.Get("interests"),
		form.Get("city"),
		form.Get("email"),
		form.Get("password"),
		models.Gender(form.Get("gender")),
		uint32(age),
	)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}

		return
	}

	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) registerUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.NewForm(nil),
	})
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.NewForm(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.NewForm(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}

		return
	}

	app.session.Put(r, "authenticatedUserID", id)

	http.Redirect(w, r, fmt.Sprintf("/user/%d", id), http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}
