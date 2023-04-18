package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/dev4ndy/snippetbox/internal/models"
)

type Config struct {
	addr      string
	staticDir string
	dsn       string
}

type Application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	config        *Config
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func NewApplication(infoLog *log.Logger, errorLog *log.Logger, config *Config, snippets *models.SnippetModel, templateCache map[string]*template.Template) *Application {
	return &Application{infoLog: infoLog, errorLog: errorLog, config: config, snippets: snippets, templateCache: templateCache}
}

func (app *Application) ServerError(rw http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) ClientError(rw http.ResponseWriter, code int) {
	http.Error(rw, http.StatusText(code), code)
}

func (app *Application) NotFound(rw http.ResponseWriter) {
	app.ClientError(rw, http.StatusNotFound)
}

func (app *Application) Routes() *http.ServeMux {
	nfs := NeuteredFileSystem{http.Dir(app.config.staticDir)}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(nfs)))

	home := NewHome(app)
	snippet := NewSnippet(app)

	mux.Handle("/", home)
	mux.Handle("/snippet", snippet)

	return mux
}

func (app *Application) Render(rw http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]

	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.ServerError(rw, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)

	if err != nil {
		app.ServerError(rw, err)
		return
	}

	rw.WriteHeader(status)

	buf.WriteTo(rw)
}

func (app *Application) NewTemplateData() *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

type NeuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs NeuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		return nil, os.ErrNotExist
	}

	return f, nil
}
