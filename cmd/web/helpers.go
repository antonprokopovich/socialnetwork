package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	if err = app.errorLog.Output(2, trace); err != nil {
		app.errorLog.Println("issue with printing error logs", err)
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) badRequest(w http.ResponseWriter) {
	app.clientError(w, http.StatusBadRequest)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CurrentYear = time.Now().Year()
	td.Flash = app.session.PopString(r, "flash")
	td.IsAuthenticated = app.isAuthenticated(r)
	td.AuthenticatedUserID = app.authenticatedUserID(r)

	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))

		return
	}

	buff := new(bytes.Buffer)

	if err := ts.Execute(buff, app.addDefaultData(td, r)); err != nil {
		app.serverError(w, err)

		return
	}

	if _, err := buff.WriteTo(w); err != nil {
		app.serverError(w, err)

		return
	}
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.session.Exists(r, "authenticatedUserID")
}

func (app *application) authenticatedUserID(r *http.Request) int {
	if !app.session.Exists(r, "authenticatedUserID") {
		return -1
	}

	return app.session.Get(r, "authenticatedUserID").(int)
}
