package server

import (
	"TodoList/internal/config"
	"TodoList/internal/pkg/logger"
	"TodoList/internal/pkg/render"
	"TodoList/internal/todos"
	"TodoList/internal/users"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/golangcollege/sessions"
)

// This holds all server dependencies!
type server struct {
	u       *users.Users
	cfg     *config.Config
	session *sessions.Session
	t       *todos.Todos
	r       *render.Render
	l       *logger.Logger
}

func New(l *logger.Logger, cfg *config.Config, s *sessions.Session, t *todos.Todos, u *users.Users, r *render.Render) (*server, error) {
	return &server{
		l:       l,
		cfg:     cfg,
		session: s,
		t:       t,
		u:       u,
		r:       r,
	}, nil
}

func (s server) Start() {
	server := &http.Server{
		Addr:    s.cfg.Address,
		Handler: s.routes(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		err := server.Shutdown(ctx)
		if err != nil {
			s.l.Fatal("unable to shutdown the server!")
		}
	}()

	s.l.Info("Started server on address %s", s.cfg.Address)

	s.l.Fatal(server.ListenAndServe())
}

func (s server) routes() *chi.Mux {
	r := chi.NewRouter()

	// middlewares
	// if !s.cfg.Production {
	// 	r.Use(s.logger) // should come first
	// }
	r.Use(s.recoverer)
	r.Use(s.session.Enable)

	// static files
	fs := http.FileServer(http.Dir("./static"))
	r.With(s.restrictStatic).
		Handle("/static/*", http.StripPrefix("/static/", fs))

	// index route
	r.Get("/", s.index)

	// todo route
	r.With(s.authMiddleware).Route("/todo", func(r chi.Router) {
		r.Post("/", s.createTodo)
		r.Put("/", s.updateTodo)
		r.Delete("/", s.deleteTodo)
		r.Get("/list", s.listTodos)
		r.Get("/id", s.lastIDTodo) // FIXME: bad route name
	})

	// user route
	r.Route("/user", func(r chi.Router) {
		r.Get("/login", s.loginUserForm)
		r.Post("/login", s.loginUser)
		r.Get("/signup", s.signupUserForm)
		r.Post("/signup", s.signupUser)
		r.With(s.authMiddleware).Post("/logout", s.logoutUser)
		r.Get("/auth", s.authUser)
		r.Get("/name", s.getUsername)
	})

	return r
}
