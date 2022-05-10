package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/golangcollege/sessions"
	"github.com/jackc/pgx"
)

const port = "4050"

var secret = []byte("gHmCldWUsVCfaUOMe0PRsHL2JhiBCKva")

type application struct {
	cacheTemplate map[string]*template.Template
	session       sessions.Session
	templateData  templateData
	userTodo      userTodo
	infoLog       *log.Logger
	errLog        *log.Logger
}

func main() {
	infoLog := log.New(os.Stdout, "INFO: ", log.LstdFlags)
	errLog := log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lshortfile)

	cacheTemplate, err := newCacheTemplate("./static")
	if err != nil {
		errLog.Fatal(err)
	}

	db, err := pgx.Connect(pgx.ConnConfig{
		Host:     "localhost",
		Password: "123456",
	})

	defer db.Close()

	app := application{
		cacheTemplate: cacheTemplate,
		session:       *sessions.New(secret),
		infoLog:       infoLog,
		errLog:        errLog,
	}

	server := http.Server{
		Addr:     fmt.Sprintf(":%s", port),
		Handler:  app.routes(),
		ErrorLog: errLog,
	}

	infoLog.Printf("Started server on port %s\n", port)
	errLog.Fatal(server.ListenAndServe())
}
