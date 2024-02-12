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

const (
	PermissionReadNews     = "read_news"
	PermissionCreateNews   = "create_news"
	PermissionReadAccDep   = "read_acc_dep"
	PermissionCreateAccDep = "create_acc_dep"
)

func (app *Application) RequirePermission(requiredPermission string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userEmail := app.GetUserMail(r)

		if !app.HasPermission(userEmail, requiredPermission) {
			http.Error(w, "Not Permitted", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func (app *Application) HasPermission(userEmail, requiredPermission string) bool {
	permissions, err := app.users.GetUserPermissions(userEmail)
	if err != nil {
		return false
	}

	return app.HavePermission(permissions, requiredPermission)
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

	app := &Application{
		errorLog:   errorLog,
		infoLog:    infoLog,
		newsModel:  &pgsql.NewsModel{DB: db},
		accountDep: &pgsql.AccountDep{DB: db},
		users:      &pgsql.UserDB{DB: db},
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
