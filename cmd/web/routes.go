package main

import (
	"github.com/evillgenius75/glennjokebox/ui"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/v1/joke/view/:id", dynamic.ThenFunc(app.jokeView))
	router.Handler(http.MethodGet, "/v1/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/v1/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/v1/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/v1/user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/v1/joke/create", protected.ThenFunc(app.jokeCreate))
	router.Handler(http.MethodPost, "/v1/joke/create", protected.ThenFunc(app.jokeCreatePost))
	router.Handler(http.MethodPost, "/v1/user/logout", protected.ThenFunc(app.userLogoutPost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
