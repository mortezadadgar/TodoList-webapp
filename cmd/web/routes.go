package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() *chi.Mux {
	r := chi.NewRouter()

	// middlewares
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(app.session.Enable)

	// static files
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// index route
	r.Get("/", app.index)
	r.Post("/", app.indexForm)

	// user route
	r.Route("/user", func(r chi.Router) {
		r.Get("/login", app.loginUserForm)
		r.Post("/login", app.loginUser)
		r.Get("/signup", app.signupUserForm)
		r.Post("/signup", app.signupUser)
		r.Post("/logout", app.logoutUser)
	})

	return r
}
