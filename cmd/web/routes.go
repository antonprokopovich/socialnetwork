package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()

	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))

	mux.Get("/user/register", dynamicMiddleware.ThenFunc(app.registerUserForm))
	mux.Post("/user/register", dynamicMiddleware.ThenFunc(app.registerUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.ThenFunc(app.logoutUser))

	mux.Get("/user/:id", dynamicMiddleware.ThenFunc(app.showUser))
	mux.Post("/user/:id", dynamicMiddleware.ThenFunc(app.showUser))

	mux.Get("/users/search", dynamicMiddleware.ThenFunc(app.findUsers))

	mux.Post("/user/:id/add", dynamicMiddleware.ThenFunc(app.sendFriendRequest))
	mux.Post("/user/:id/accept", dynamicMiddleware.ThenFunc(app.acceptFriendRequest))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
