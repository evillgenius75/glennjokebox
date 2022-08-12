package main

import (
	"errors"
	"fmt"
	"github.com/evillgenius75/glennjokebox/internal/models"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	jokes, err := app.jokes.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, joke := range jokes {
		fmt.Fprintf(w, "%v+\n", joke)
	}

	//files := []string{
	//	"./ui/html/base.tmpl",
	//	"./ui/html/pages/home.tmpl",
	//	"./ui/html/partials/nav.tmpl",
	//}
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	app.serverError(w, err)
	//	return
	//}
	//
	//err = ts.ExecuteTemplate(w, "base", nil)
	//if err != nil {
	//	app.serverError(w, err)
	//}
}

func (app *application) jokeView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/pages/views.tmpl",
		"./ui/html/partials/nav.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Joke: joke,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}

	//w.Write([]byte("Display a random Joke from DB..."))
	//fmt.Fprintf(w, "%+v", joke)
}

func (app *application) jokeCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	uuid := "sadfe87viha"
	joke := "I spent a whole bunch of money childproofing my house. It doesn't seem to be working because they still keep coming in!"

	id, err := app.jokes.Insert(uuid, joke)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/joke/view?id=%d", id), http.StatusSeeOther)

	//w.Write([]byte("Create a new joke..."))
}
