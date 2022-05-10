package main

import (
	"net/http"
)

func (app *application) render(w http.ResponseWriter, name string, data interface{}) {
	tmpl, ok := app.cacheTemplate[name]
	if !ok {
		app.serverError(w, "invalid template name")
	}

	err := tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		app.serverError(w, err.Error())
	}
}

func (app *application) serverError(w http.ResponseWriter, error string) {
	http.Error(w, error, http.StatusInternalServerError)
	w.WriteHeader(http.StatusInternalServerError)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.session.GetBool(r, "Authenticated")
}
