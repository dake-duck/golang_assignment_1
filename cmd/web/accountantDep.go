package main

import (
	"DAKExDUCK/assignment_1/pkg/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

func (app *Application) createAccDep(w http.ResponseWriter, r *http.Request) {
	if !app.isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		app.createAccDepGET(w, r)
	} else if r.Method == http.MethodPost {
		app.createAccDepPOST(w, r)
	} else {
		http.Error(w, "Method not found", 405)
	}
}

func (app *Application) createAccDepGET(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/templates/account_dep_create.page.tmpl",
		"./ui/templates/base.layout.tmpl",
		"./ui/templates/footer.layout.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, struct {
		Title string
	}{Title: "Accounts of Department: Create a new one"})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (app *Application) createAccDepPOST(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	name := r.Form.Get("name")
	sname := r.Form.Get("sname")
	age := r.Form.Get("age")
	ageInt, err := strconv.Atoi(age)
	if err != nil {
		panic(err)
	}

	id, err := app.accountDep.Insert(name, sname, ageInt)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"status": "success",
		"ID":     id,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	// Set content type and write the JSON response
	w.Header().Set("Content-Type", "Application/json")
	w.Write(jsonResponse)
}

func (app *Application) showAccDepList(w http.ResponseWriter, r *http.Request) {
	if !app.isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var accounts []models.AccountDep
	var err error
	id, err_id := strconv.Atoi(r.URL.Query().Get("id"))

	if err_id == nil {
		accounts, err = app.accountDep.Get(id)
	} else if err_id != nil {
		accounts, err = app.accountDep.GetAll()
	}

	if err != nil {
		if !errors.Is(err, models.ErrNoRecord) {
			app.serverError(w, err)
			return
		}
	}

	files := []string{
		"./ui/templates/account_dep.page.tmpl",
		"./ui/templates/base.layout.tmpl",
		"./ui/templates/footer.layout.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, struct {
		Accounts []models.AccountDep
		Title    string
	}{Accounts: accounts, Title: "Accounts of Department"})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}
