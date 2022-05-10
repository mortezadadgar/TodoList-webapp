package main

import (
	"time"
)

type entry struct {
	Body        string
	isDone      bool
	createdDate string
}

type userTodo struct {
	user        string
	TodoEntries []entry
}

func (t *userTodo) createTodo(body string, isDone bool) {
	e := entry{
		Body:        body,
		isDone:      false,
		createdDate: time.Now().Format("Mon Jan 2 15:04:05 2006"),
	}
	t.TodoEntries = append(t.TodoEntries, e)
}
