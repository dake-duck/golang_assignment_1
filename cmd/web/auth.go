package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
)

var (
	cookieHandler = securecookie.New(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32))
)

type User struct {
	Email        string
	PasswordHash string
}

func (app *Application) loginHandler(w http.ResponseWriter, r *http.Request) {
	if app.isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	files := []string{
		"./ui/templates/login.page.tmpl",
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
	}{})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (app *Application) loginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if app.isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.FormValue("username")
	password := r.FormValue("password")

	users, err := app.users.Get(email)
	if err != nil || len(users) == 0 || !checkPasswordHash(password, users[email].PasswordHash) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	setSession(email, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) registerHandler(w http.ResponseWriter, r *http.Request) {
	if app.isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	files := []string{
		"./ui/templates/register.page.tmpl",
		"./ui/templates/base.layout.tmpl",
		"./ui/templates/footer.layout.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, struct{}{})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (app *Application) registerSubmitHandler(w http.ResponseWriter, r *http.Request) {
	if app.isLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.FormValue("username")
	password := r.FormValue("password")

	_, exists := app.users.Get(email)
	if exists == nil {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	hashedPassword := hashPassword(password)

	app.users.Insert(email, hashedPassword)

	setSession(email, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) logoutHandler(w http.ResponseWriter, r *http.Request) {
	clearSession(w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func setSession(email string, w http.ResponseWriter) {
	value := map[string]string{
		"email": email,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}
}

func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func (app *Application) isLoggedIn(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false
	}

	value := make(map[string]string)
	if err = cookieHandler.Decode("session", cookie.Value, &value); err == nil {
		users, err := app.users.Get(value["email"])
		if err != nil {
			return false
		}
		_, exists := users[value["email"]]
		return exists
	}

	return false
}

func (app *Application) GetUserMail(r *http.Request) string {
	cookie, err := r.Cookie("session")
	if err != nil {
		return ""
	}

	value := make(map[string]string)
	if err = cookieHandler.Decode("session", cookie.Value, &value); err == nil {
		return value["email"]
	}

	return ""
}

func (app *Application) HavePermission(permissions []string, requiredPermission string) bool {
	for _, p := range permissions {
		if p == requiredPermission {
			return true
		}
	}
	return false
}

func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
