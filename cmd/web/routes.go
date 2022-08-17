package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	fileServer := http.FileServer(http.Dir(cfg.staticDir))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/v1/joke/view/:id", dynamic.ThenFunc(app.jokeView))
	router.Handler(http.MethodGet, "/v1/joke/create", dynamic.ThenFunc(app.jokeCreate))
	router.Handler(http.MethodPost, "/v1/joke/create", dynamic.ThenFunc(app.jokeCreatePost))

	router.Handler(http.MethodGet, "/v1/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/v1/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/v1/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/v1/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/v1/user/logout", dynamic.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
