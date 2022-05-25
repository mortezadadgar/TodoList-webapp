package users

import (
	"TodoList/internal/config"
	"TodoList/internal/pkg/logger"
	"TodoList/internal/pkg/render"
	"TodoList/internal/pkg/validate"
	"TodoList/internal/store"
	"errors"
	"net/http"
	"time"

	"github.com/golangcollege/sessions"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("invalid user passwrod")
	ErrInvalidName     = errors.New("invalid template name")
	ErrEmptyInput      = errors.New("empty input field is not allowed")
	ErrSmallInputLen   = errors.New("input with this length is not allowed")
)

type dynamicData struct {
	UserName        string
	IsAuthenticated bool
}

// holds all dependencies
type Users struct {
	l        *logger.Logger
	store    store.User
	cfg      *config.Config
	r        *render.Render
	session  *sessions.Session
	UserData dynamicData
}

func New(l *logger.Logger, cfg *config.Config, store store.User, r *render.Render, s *sessions.Session) *Users {
	return &Users{
		l:     l,
		store: store,
		cfg:   cfg,
		r:     r,
	}
}

func (u *Users) Validate(username string, password string) error {
	hashedPasswrod, err := u.store.GetHashedPassword(username)
	if err != nil {
		return err
	}

	err = u.validatePassword(password, hashedPasswrod)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidPassword
		}
		return err
	}

	return nil
}

func (u *Users) Create(username string, password string) error {
	// FIXME: date in todo_list's created_date doesn't follow this format!
	currentDate := time.Now().Format("Mon Jan 2 15:04:05 2006")
	// currentDate = pq.QuoteLiteral(currentDate)

	hashedPasswrod, err := u.generatePassword(password)
	if err != nil {
		return err
	}

	return u.store.CreateUser(username, password, currentDate, hashedPasswrod)
}

func (u *Users) validatePassword(password string, hashedPasswrod string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPasswrod), []byte(password))
}

func (u *Users) generatePassword(password string) (string, error) {
	hashedPasswrod, err := bcrypt.GenerateFromPassword([]byte(password), u.cfg.PasswordCost)
	if err != nil {
		return "", err
	}
	return string(hashedPasswrod), nil
}

func (u *Users) SignupUserForm(w http.ResponseWriter, r *http.Request) error {
	return u.r.Render(w, "signup.page.tmpl", nil)
}

// TODO:this pointer?
func (u *Users) LogoutUser(w http.ResponseWriter, r *http.Request) {
	u.session.Destroy(r)
	u.UserData.IsAuthenticated = false
}

func (u *Users) SignupUser(w http.ResponseWriter, r *http.Request, username string, password string) error {
	ok := validate.Len(0, password, username)
	if !ok {
		return ErrEmptyInput
	}

	ok = validate.Len(8, password)
	if !ok {
		return ErrSmallInputLen
	}

	err := u.Create(username, password)
	if err != nil {
		return err
	}

	return nil
}

func (u *Users) LoginUserForm(w http.ResponseWriter, r *http.Request) error {
	return u.r.Render(w, "login.page.tmpl", nil)
}

func (u *Users) LoginUser(w http.ResponseWriter, r *http.Request, username string, password string) error {
	_, err := u.store.GetUserID(username)
	if err != nil {
		return err
	}

	err = u.Validate(username, password)
	if err != nil {
		return err
	}

	u.session.Put(r, "Authenticated", true)
	u.session.Put(r, "username", username)
	u.UserData.IsAuthenticated = true
	u.UserData.UserName = username

	return nil
}

func (u *Users) IsAuthenticated(r *http.Request) bool {
	return u.session.GetBool(r, "Authenticated")
}

func (u *Users) GetUsername(r *http.Request) string {
	return u.session.GetString(r, "username")
}
