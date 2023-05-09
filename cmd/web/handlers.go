package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/dev4ndy/snippetbox/internal/models"
	"github.com/julienschmidt/httprouter"
)

func (app *Application) HomeView(rw http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()

	if err != nil {
		app.ServerError(rw, err)
	}

	data := app.NewTemplateData()
	data.Snippets = snippets

	app.Render(rw, http.StatusOK, "home.html", data)
}

func (app *Application) SnippetView(rw http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.NotFound(rw)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.NotFound(rw)
		} else {
			app.ServerError(rw, err)
		}
		return
	}

	data := app.NewTemplateData()
	data.Snippet = snippet

	app.Render(rw, http.StatusOK, "view.html", data)

}

type SnippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

func (app *Application) SnippetCreate(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ClientError(rw, http.StatusBadRequest)
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))

	if err != nil {
		app.ClientError(rw, http.StatusBadRequest)
	}

	form := SnippetCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{},
	}

	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field is required"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field can not be more than 100 characters long"
	}

	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field is required"
	}

	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "The expiration should be between 1 day and 1 year"
	}

	if len(form.FieldErrors) > 0 {
		data := app.NewTemplateData()
		data.Form = form
		app.Render(rw, http.StatusUnprocessableEntity, "create.html", data)
		fmt.Fprint(rw, form.FieldErrors)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.ServerError(rw, err)
		return
	}

	http.Redirect(rw, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}

func (app *Application) SnippetCreateView(rw http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData()
	data.Form = SnippetCreateForm{
		Expires: 365,
	}
	app.Render(rw, http.StatusOK, "create.html", data)
}
