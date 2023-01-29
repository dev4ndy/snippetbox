package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

type Config struct {
	addr      string
	staticDir string
}

type Application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func NewApplication(infoLog *log.Logger, errorLog *log.Logger) *Application {
	return &Application{infoLog: infoLog, errorLog: errorLog}
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
