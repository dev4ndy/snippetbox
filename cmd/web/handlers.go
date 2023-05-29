package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dev4ndy/snippetbox/internal/models"
	"github.com/dev4ndy/snippetbox/internal/validator"
	"github.com/julienschmidt/httprouter"
)

func (app *Application) HomeView(rw http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()

	if err != nil {
		app.ServerError(rw, err)
	}

	data := app.NewTemplateData(r)
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

	data := app.NewTemplateData(r)
	data.Snippet = snippet

	app.Render(rw, http.StatusOK, "view.html", data)

}

type SnippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *Application) SnippetCreate(rw http.ResponseWriter, r *http.Request) {
	var form SnippetCreateForm

	err := app.DecodePostForm(r, &form)

	if err != nil {
		app.ClientError(rw, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.Required(form.Title), "title", "This field is required")
	form.CheckField(validator.MaxLength(form.Title, 100), "title", "This field can not be more than 100 characters long")
	form.CheckField(validator.Required(form.Content), "content", "This field is required")
	form.CheckField(validator.AllowedValues(form.Expires, 1, 7, 365), "expires", "The expiration should be between 1 day and 1 year")

	if form.Invalid() {
		data := app.NewTemplateData(r)
		data.Form = form
		app.Render(rw, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.ServerError(rw, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(rw, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}

func (app *Application) SnippetCreateView(rw http.ResponseWriter, r *http.Request) {
	data := app.NewTemplateData(r)
	data.Form = SnippetCreateForm{
		Expires: 365,
	}
	app.Render(rw, http.StatusOK, "create.html", data)
}
