package models

// TODO: add email
type User struct {
	Username       string
	Password       string
	hashedPasswrod string
	createdDate    string
	ID             int
}
