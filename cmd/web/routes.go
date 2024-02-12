package main

import "net/http"

func (app *Application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.RequirePermission(PermissionReadNews, app.news))
	mux.HandleFunc("/test", app.test)
	mux.HandleFunc("/test1", app.test1)
	mux.HandleFunc("/login", app.loginHandler)
	mux.HandleFunc("/login-submit", app.loginSubmitHandler)
	mux.HandleFunc("/register", app.registerHandler)
	mux.HandleFunc("/register-submit", app.registerSubmitHandler)
	mux.HandleFunc("/logout", app.logoutHandler)
	mux.HandleFunc("/create", app.RequirePermission(PermissionCreateNews, app.createNews))
	mux.HandleFunc("/accounts_dep", app.RequirePermission(PermissionReadAccDep, app.showAccDepList))
	mux.HandleFunc("/accounts_dep/create", app.RequirePermission(PermissionCreateAccDep, app.createAccDep))
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	return mux
}
