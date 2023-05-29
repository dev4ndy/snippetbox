package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/dev4ndy/snippetbox/internal/models"
	"github.com/go-playground/form/v4"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
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
	formDecoder   *form.Decoder
}

func NewApplication(infoLog *log.Logger, errorLog *log.Logger, config *Config, snippets *models.SnippetModel, templateCache map[string]*template.Template, formDecoder *form.Decoder) *Application {
	return &Application{infoLog: infoLog, errorLog: errorLog, config: config, snippets: snippets, templateCache: templateCache, formDecoder: formDecoder}
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

func (app *Application) Routes() http.Handler {
	nfs := NeuteredFileSystem{http.Dir(app.config.staticDir)}

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		app.NotFound(rw)
	})

	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", http.FileServer(nfs)))

	router.HandlerFunc(http.MethodGet, "/", app.HomeView)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.SnippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.SnippetCreateView)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.SnippetCreate)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
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

func (app *Application) DecodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Call Decode() on our decoder instance, passing the target destination as
	// the first parameter.
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// If we try to use an invalid target destination, the Decode() method
		// will return an error with the type *form.InvalidDecoderError.We use
		// errors.As() to check for this and raise a panic rather than returning
		// the error.
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// For all other errors, we return them as normal.
		return err
	}

	return nil
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
