package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
)

type Config struct {
	addr      string
	staticDir string
}

type Application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	config   *Config
}

func NewApplication(infoLog *log.Logger, errorLog *log.Logger, config *Config) *Application {
	return &Application{infoLog: infoLog, errorLog: errorLog, config: config}
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
