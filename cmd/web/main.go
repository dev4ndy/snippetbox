package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type config struct {
	addr      string
	staticDir string
}

func main() {

	var cfg config

	// command-line flags
	// To change the Http network default address, use: `go run main.go -addr=":80"`
	flag.StringVar(&cfg.addr, "addr", ":8080", "Http network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "../../ui/static/", "Path to static assets")
	// We need to use flag.Parse before to use the flags
	// otherwise it will take the default value
	flag.Parse()

	mux := http.NewServeMux()
	nfs := neuteredFileSystem{http.Dir(cfg.staticDir)}
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(nfs)))
	home := NewHome()

	mux.Handle("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Printf("Starting Server on %s", cfg.addr)
	err := http.ListenAndServe(cfg.addr, mux)
	log.Fatal(err)
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
