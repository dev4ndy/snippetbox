package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dev4ndy/snippetbox/internal/models"
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
	snippets, err := h.app.snippets.Latest()

	if err != nil {
		h.app.ServerError(rw, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(rw, "%+v\n", snippet)
	}
	// files := []string{
	// 	"../../ui/html/main.tmpl",
	// 	"../../ui/html/partials/nav.tmpl",
	// 	"../../ui/html/pages/home.tmpl",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	h.app.ServerError(rw, err)
	// }

	// err = ts.ExecuteTemplate(rw, "base", nil)
	// if err != nil {
	// 	h.app.ServerError(rw, err)
	// }
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

	snippet, err := s.app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			s.app.NotFound(rw)
		} else {
			s.app.ServerError(rw, err)
		}
		return
	}

	fmt.Fprintf(rw, "%+v", snippet)
}

func (s *Snippet) Create(w http.ResponseWriter, r *http.Request) {
	title := "Andres"
	content := "Content"
	expires := 7

	id, err := s.app.snippets.Insert(title, content, expires)
	if err != nil {
		s.app.ServerError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

}
