package main

import (
	db "DAKExDUCK/assignment_1/internal/db"
	pgsql "DAKExDUCK/assignment_1/pkg/models/pgsql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type application struct {
	errorLog   *log.Logger
	infoLog    *log.Logger
	newsModel  *pgsql.NewsModel
	accountDep *pgsql.AccountDep
}

func main() {
	addr := flag.String("addr", "localhost:8080", "HTTP network address")
	dsn := flag.String("dsn", "user:pass@host:port", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := db.OpenDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog:   errorLog,
		infoLog:    infoLog,
		newsModel:  &pgsql.NewsModel{DB: db},
		accountDep: &pgsql.AccountDep{DB: db},
	}
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	log.Fatal(err)
}
