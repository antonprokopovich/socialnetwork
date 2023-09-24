package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
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

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))

		return
	}

	buff := new(bytes.Buffer)

	if err := ts.Execute(buff, td); err != nil {
		app.serverError(w, err)

		return
	}

	if _, err := buff.WriteTo(w); err != nil {
		app.serverError(w, err)

		return
	}
}