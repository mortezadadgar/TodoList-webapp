package store

import (
	"TodoList/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

//TODO: manage multiple call to same row!!

var (
	ErrUserNotFound    = errors.New("user not found in database")
	ErrInvalidPassword = errors.New("invalid user passwrod")
)

type postgreStore struct {
	db *sql.DB
}

func NewPostgreStore(db *sql.DB) (*postgreStore, error) {
	return &postgreStore{
		db: db,
	}, nil
}

func (p postgreStore) CreateUser(username string, password string, currentDate string, hashedPasswrod string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.db.ExecContext(ctx,
		`INSERT INTO users(username, hashed_password, created_date) VALUES($1, $2, $3)`,
		username,
		hashedPasswrod,
		currentDate,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p postgreStore) GetHashedPassword(username string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var hashedPasswrod string
	err := p.db.QueryRowContext(ctx,
		`SELECT hashed_password FROM users WHERE username = $1`, username).Scan(&hashedPasswrod)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrUserNotFound
		}
		return "", err
	}

	return hashedPasswrod, nil
}

func (p postgreStore) GetUserID(username string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var userID int
	err := p.db.QueryRowContext(ctx,
		`SELECT id FROM users WHERE username = $1`, username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, ErrUserNotFound
		}
		return 0, err
	}

	return userID, nil
}

func (p postgreStore) DeleteUser(u *models.User) error {
	return nil
}

func (p postgreStore) CreateTodo(userID int, text string, createdDate string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.db.ExecContext(ctx, `INSERT INTO todos(user_id, body, is_done, created_date) VALUES($1, $2, $3, $4)`,
		userID,
		text,
		false,
		createdDate)
	if err != nil {
		return err
	}

	fmt.Printf("Added userID: %d, text: %s\n", userID, text)

	return nil
}

// NOTE: No need for calling Close here as Next() returns false for last row
func (p postgreStore) GetTodos(userId int) (*[]models.UserTodo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// TODO: Duplicate todos?
	todos := []models.UserTodo{}
	rows, err := p.db.QueryContext(ctx, `SELECT id, body, is_done FROM todos WHERE user_id = $1`, userId)
	if err != nil {
		return nil, err
	}

	var (
		todoText string
		isDone   bool
		id       int
	)
	for rows.Next() {
		err := rows.Scan(&id, &todoText, &isDone)
		if err != nil {
			return nil, err
		}
		todos = append(todos, models.UserTodo{ID: id, Text: todoText, IsDone: isDone})
	}

	return &todos, nil
}

func (p postgreStore) UpdateTodo(id int, text string, isDone bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.db.ExecContext(ctx, `UPDATE todos SET body = $1, is_done = $2 WHERE id = $3`, text, isDone, id)
	if err != nil {
		return nil
	}

	return nil
}

func (p postgreStore) GetLastID() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	err := p.db.QueryRowContext(ctx, `SELECT id FROM todos ORDER BY id DESC LIMIT 1`).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p postgreStore) DeleteTodo(id int, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.db.ExecContext(ctx, `DELETE from todos WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return err
	}

	return nil
}
