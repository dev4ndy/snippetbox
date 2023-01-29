package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type Config struct {
	addr      string
	staticDir string
}

type ApplicationLogger struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	var cfg Config

	// command-line flags
	// To change the Http network default address, use: `go run main.go -addr=":80"`
	flag.StringVar(&cfg.addr, "addr", ":8080", "Http network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "../../ui/static/", "Path to static assets")
	// We need to use flag.Parse before to use the flags
	// otherwise it will take the default value
	flag.Parse()

	// Logging
	// To redirect the stdout and stderr to a file on-disk:
	// `go run main.go >>/tmp/info.log 2>>/tmp/error.log`
	appLogger := ApplicationLogger{
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Llongfile),
	}

	nfs := neuteredFileSystem{http.Dir(cfg.staticDir)}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(nfs)))

	home := NewHome(&appLogger)
	snippet := NewSnippet(&appLogger)

	mux.Handle("/", home)
	mux.Handle("/snippet", snippet)

	server := http.Server{
		Addr:     cfg.addr,
		ErrorLog: appLogger.errorLog,
		Handler:  mux,
	}

	appLogger.infoLog.Printf("Starting Server on %s", cfg.addr)
	err := server.ListenAndServe()
	appLogger.errorLog.Fatal(err)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
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
