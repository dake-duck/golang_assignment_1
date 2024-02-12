package main

import (
	"net/http"
)

func (app *Application) test(w http.ResponseWriter, r *http.Request) {
	mail := app.GetUserMail(r)
	permissions := []string{
		PermissionReadNews,
		PermissionCreateNews,
		PermissionReadAccDep,
		PermissionCreateAccDep,
	}
	app.users.SetPermissions(mail, permissions)

	app.news(w, r)
}

func (app *Application) test1(w http.ResponseWriter, r *http.Request) {
	mail := app.GetUserMail(r)
	permissions := []string{}
	app.users.SetPermissions(mail, permissions)

	app.news(w, r)
}
