package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.news)
	mux.HandleFunc("/create", app.createNews)
	mux.HandleFunc("/accounts_dep", app.showAccDepList)
	mux.HandleFunc("/accounts_dep/create", app.createAccDep)
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
