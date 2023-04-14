package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/dev4ndy/snippetbox/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	var cfg Config

	// command-line flags
	// To change the Http network default address, use: `go run main.go -addr=":80"`
	flag.StringVar(&cfg.addr, "addr", ":8080", "Http network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "../../ui/static/", "Path to static assets")

	// Data Base
	flag.StringVar(&cfg.dsn, "dns", "snippet:snippet@tcp(mysql:3306)/snippetbox?parseTime=true", "MySQL data source name")
	// We need to use flag.Parse before to use the flags
	// otherwise it will take the default value
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(cfg.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	// Logging
	// To redirect the stdout and stderr to a file on-disk:
	// `go run main.go >>/tmp/info.log 2>>/tmp/error.log`
	app := NewApplication(
		infoLog,
		errorLog,
		&cfg,
		&models.SnippetModel{DB: db},
	)

	defer db.Close()

	server := http.Server{
		Addr:     cfg.addr,
		ErrorLog: app.errorLog,
		Handler:  app.Routes(),
	}

	app.infoLog.Printf("Starting Server on %s", cfg.addr)
	err = server.ListenAndServe()
	app.errorLog.Fatal(err)
}

func openDB(dns string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
