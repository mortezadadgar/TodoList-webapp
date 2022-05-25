package store

import "TodoList/internal/models"

type User interface {
	CreateUser(string, string, string, string) error
	GetUserID(string) (int, error)
	GetHashedPassword(string) (string, error)
	DeleteUser(*models.User) error
}

type Todo interface {
	CreateTodo(int, string, string) error
	UpdateTodo(int, string, bool) error
	DeleteTodo(int, int) error
	GetTodos(int) (*[]models.UserTodo, error)
	GetLastID() (int, error)
}
