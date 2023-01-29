package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

type Home struct {
	app *Application
}

func NewHome(app *Application) *Home {
	return &Home{app}
}

func (h *Home) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(rw, r)
		return
	}
	h.View(rw, r)
}

func (h *Home) View(rw http.ResponseWriter, r *http.Request) {
	files := []string{
		"../../ui/html/main.tmpl",
		"../../ui/html/partials/nav.tmpl",
		"../../ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		h.app.ServerError(rw, err)
	}

	err = ts.ExecuteTemplate(rw, "base", nil)
	if err != nil {
		h.app.ServerError(rw, err)
	}
}

type Snippet struct {
	app *Application
}

func NewSnippet(app *Application) *Snippet {
	return &Snippet{app}
}

func (s *Snippet) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.View(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		s.Create(rw, r)
		return
	}
	rw.Header().Set("Allow", http.MethodPost)
	rw.Header().Set("Allow", http.MethodGet)
	s.app.ClientError(rw, http.StatusMethodNotAllowed)
}

func (s *Snippet) View(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		s.app.NotFound(rw)
		return
	}
	fmt.Fprintf(rw, "Display a specific snippet with ID %d...", id)
}

func (s *Snippet) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new snippet"))
}
