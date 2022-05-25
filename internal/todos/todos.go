package todos

import (
	"TodoList/internal/models"
	"TodoList/internal/pkg/logger"
	"TodoList/internal/store"
	"database/sql"
	"encoding/json"
	"time"
)

type Todos struct {
	todoStore store.Todo
	userStore store.User
	db        *sql.DB
	l         *logger.Logger
	UserTodo  *models.UserTodo
}

func New(l *logger.Logger, db *sql.DB, todoStore store.Todo, userStore store.User) *Todos {
	return &Todos{
		l:         l,
		db:        db,
		todoStore: todoStore,
		userStore: userStore,
	}
}

func (t *Todos) CreateTodo(text string, username string) error {
	userID, err := t.userStore.GetUserID(username)
	if err != nil {
		return err
	}

	CreatedDate := time.Now().Format("Mon Jan 2 15:04:05 2006") //TODO: make it config

	err = t.todoStore.CreateTodo(userID, text, CreatedDate)
	if err != nil {
		return err
	}

	return nil
}

type envelope struct {
	Data interface{} `json:"data"`
}

// DONE: send as a json request so the isDone can be send as well
func (t *Todos) GetTodos(username string) ([]byte, error) {
	userID, err := t.userStore.GetUserID(username)
	if err != nil {
		return nil, err
	}

	todos, err := t.todoStore.GetTodos(userID)
	if err != nil {
		return nil, err
	}

	encodedTodos, err := json.Marshal(envelope{todos})
	if err != nil {
		return nil, err
	}

	return encodedTodos, nil
}

func (t *Todos) DeleteTodo(id int, username string) error {
	userID, err := t.userStore.GetUserID(username)
	if err != nil {
		return err
	}

	err = t.todoStore.DeleteTodo(id, userID)
	if err != nil {
		return err
	}

	return nil
}

func (t *Todos) UpdateTodo(id int, text string, isDone bool) error {
	err := t.todoStore.UpdateTodo(id, text, isDone)
	if err != nil {
		return err
	}

	return nil
}

func (t *Todos) GetLastID() (int, error) {
	id, err := t.todoStore.GetLastID()
	if err != nil {
		return 0, err
	}

	return id, nil
}
