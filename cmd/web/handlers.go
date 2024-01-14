package main

import (
	"DAKExDUCK/assignment_1/pkg/models"
	"errors"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

func (app *application) news(w http.ResponseWriter, r *http.Request) {
	var news []models.News
	var err error
	Title := "News"

	id, err_id := strconv.Atoi(r.URL.Query().Get("id"))
	category, err_category := strconv.Atoi(r.URL.Query().Get("category"))

	if err_id == nil {
		news, err = app.snippets.Get(id)
	} else if err_category == nil {
		news, err = app.snippets.GetByCategory(category)
	} else if err_category != nil && err_id != nil {
		news, err = app.snippets.Latest()
		Title = "Latest news"
	}

	if err != nil {
		if !errors.Is(err, models.ErrNoRecord) {
			app.serverError(w, err)
			return
		}
	}

	files := []string{
		"./ui/templates/home.page.tmpl",
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
		News  []models.News
		Title string
	}{News: news, Title: Title})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
