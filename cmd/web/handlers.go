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
		return
	}
}

func (app *application) createNews(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		app.createNewsGet(w, r)
	} else if r.Method == http.MethodPost {
		app.createNewsPost(w, r)
	} else {
		http.Error(w, "Method not found", 405)
	}
}

func (app *application) createNewsGet(w http.ResponseWriter, r *http.Request) {
	var err error
	Title := "Create News"

	Category, _ := strconv.Atoi(r.URL.Query().Get("category"))

	files := []string{
		"./ui/templates/createNews.page.tmpl",
		"./ui/templates/base.layout.tmpl",
		"./ui/templates/footer.layout.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	Tags, err := app.snippets.GetTags()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, struct {
		Title    string
		Category int
		Tags     []models.Tag
	}{Title: Title, Category: Category, Tags: Tags})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (app *application) createNewsPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		log.Println(err.Error())
		return
	}

	title := r.Form.Get("title")
	body := r.Form.Get("body")

	var categories []int
	for _, catStr := range r.Form["category"] {
		catInt, err := strconv.Atoi(catStr)
		if err != nil {
			http.Error(w, "Error converting category to integer", http.StatusBadRequest)
			log.Println(err.Error())
			return
		}
		categories = append(categories, catInt)
	}

	newsID, err := app.snippets.InsertWithTags(0, title, body, categories)
	if err != nil {
		http.Error(w, "Error inserting news post", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	response := map[string]interface{}{
		"status": "success",
		"newsID": newsID,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	// Set content type and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
