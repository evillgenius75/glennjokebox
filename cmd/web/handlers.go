package main

import (
	"errors"
	"fmt"
	"github.com/evillgenius75/glennjokebox/internal/models"
	"github.com/evillgenius75/glennjokebox/internal/validator"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type jokeCreateForm struct {
	UserName string
	Content  string
	Explicit int
	validator.Validator
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	jokes, err := app.jokes.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Jokes = jokes

	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) jokeView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	joke, err := app.jokes.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Joke = joke
	app.render(w, http.StatusOK, "views.tmpl", data)
}

func (app *application) jokeCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = jokeCreateForm{
		Explicit: 0,
	}

	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) jokeCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	explicit, err := strconv.Atoi(r.PostForm.Get("explicit"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := jokeCreateForm{
		UserName: r.PostForm.Get("username"),
		Content:  r.PostForm.Get("content"),
		Explicit: explicit,
	}

	form.CheckField(validator.NotBlank(form.UserName), "username", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.UserName, 100), "username", "This field cannot be more that 100 characters long")

	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.jokes.Insert(form.UserName, form.Content, form.Explicit)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/v1/joke/view/%d", id), http.StatusSeeOther)

	//w.Write([]byte("Create a new joke..."))
}
