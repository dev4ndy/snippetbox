package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	nfs := neuteredFileSystem{http.Dir("../../ui/static/")}
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(nfs)))
	home := NewHome()

	mux.Handle("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Println("Starting Server on : 8080")
	err := http.ListenAndServe(":8080", mux)
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
