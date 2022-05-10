package main

import (
	"fmt"
	"net/http"
)

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	app.render(w, "index.page.tmpl", app.addDefaultData(r, &templateData{
		UserTodo: app.userTodo,
	}))
}

func (app *application) indexForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err.Error())
	}

	body := r.FormValue("todo-input")
	if len(body) == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	app.userTodo.createTodo(body, false)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "sign up new users form")
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Destroy(r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "sign up new users")
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, "login.page.tmpl", app.addDefaultData(r, nil))
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err.Error())
	}

	password := r.Form.Get("password")
	username := r.Form.Get("username")
	app.infoLog.Println("pass:", password, "user:", username)

	// TODO: Flash a message to user
	// check
	// if len(password) < 8 {
	// 	app.infoLog.Println("pass is too small")
	// 	return
	// }
	// if len(password) == 0 || len(username) == 0 {
	// 	app.infoLog.Println("empty user or pass")
	// 	return
	// }

	// session
	app.session.Put(r, "Authenticated", true)
	app.session.Put(r, "username", username)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
