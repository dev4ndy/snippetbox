package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

type Home struct {
	appLogger *ApplicationLogger
}

func NewHome(appLogger *ApplicationLogger) *Home {
	return &Home{appLogger}
}

func (h *Home) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(rw, r)
		return
	}
	h.View(rw, r)
}

func (h *Home) View(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"../../ui/html/main.tmpl",
		"../../ui/html/partials/nav.tmpl",
		"../../ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		h.appLogger.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		h.appLogger.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Write([]byte("Hello World!"))
}

type Snippet struct {
	appLogger *ApplicationLogger
}

func NewSnippet(appLogger *ApplicationLogger) *Snippet {
	return &Snippet{appLogger}
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
	http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func (s *Snippet) View(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (s *Snippet) Create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new snippet"))
}
