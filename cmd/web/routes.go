package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := pat.New()

	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/user/register", http.HandlerFunc(app.registerUserForm))
	mux.Post("/user/register", http.HandlerFunc(app.registerUser))
	mux.Get("/user/:id", http.HandlerFunc(app.showUser))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
