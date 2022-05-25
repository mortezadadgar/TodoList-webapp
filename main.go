package main

import (
	"TodoList/internal/config"
	"TodoList/internal/pkg/cachetemplate"
	"TodoList/internal/pkg/logger"
	"TodoList/internal/pkg/render"
	"TodoList/internal/server"
	"TodoList/internal/store"
	"TodoList/internal/todos"
	"TodoList/internal/users"
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"github.com/golangcollege/sessions"
)

func main() {
	l := logger.New("02 Jan 15:04", 3)

	cfg, err := config.New()
	if err != nil {
		l.Fatal(err)
	}

	cachedTemplate, err := cachetemplate.New("./static/templates")
	if err != nil {
		l.Fatal(err)
	}

	db, err := openDB(cfg.ConnectionString)
	if err != nil {
		l.Fatal(err)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			l.Fatal(err)
		}
	}()

	ps, err := store.NewPostgreStore(db)
	if err != nil {
		l.Fatal(err)
	}

	r := render.NewRender(cachedTemplate)

	session := sessions.New([]byte(cfg.SecretKey))

	t := todos.New(l, db, ps, ps)

	u := users.New(l, cfg, ps, r, session)

	// TODO: hell of dependencies
	s, err := server.New(l, cfg, session, t, u, r)
	if err != nil {
		l.Fatal(err)
	}

	s.Start()
}

// TODO: should we here?
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, err
}
