package server

import (
	"TodoList/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var (
	ErrInvalidName   = errors.New("invalid template name")
	ErrEmptyInput    = errors.New("empty input field is not allowed")
	ErrSmallInputLen = errors.New("input with this length is not allowed")
)

// move this func to somewhere else
func (s server) index(w http.ResponseWriter, r *http.Request) {
	// TODO: HACKY?!
	// * Absolute *
	s.u.UserData.IsAuthenticated = s.u.IsAuthenticated(r)
	s.u.UserData.UserName = s.u.GetUsername(r)
	err := s.r.Render(w, "index.page.tmpl", s.u.UserData)
	if err != nil {
		s.l.Fatal(err)
	}
}

func (s server) listTodos(w http.ResponseWriter, r *http.Request) {
	todosText, err := s.t.GetTodos(s.u.GetUsername(r))
	if err != nil {
		s.l.Fatal(err)
	}

	_, err = w.Write(todosText)
	if err != nil {
		s.l.Fatal(err)
	}
}

type envelope struct {
	Data interface{} `json:"data"`
}

func (s server) createTodo(w http.ResponseWriter, r *http.Request) {
	text, err := io.ReadAll(r.Body)
	if err != nil {
		s.l.Fatal(err)
	}

	if len(text) == 0 {
		s.l.Info("Empty todo item")
		return
	}

	err = s.t.CreateTodo(string(text), s.u.GetUsername(r))
	if err != nil {
		s.l.Fatal(err)
	}
}

func (s server) updateTodo(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		s.l.Fatal(err)
	}

	ut := &models.UserTodo{}
	err = json.Unmarshal(b, &envelope{ut})
	if err != nil {
		s.l.Fatal(err)
	}

	err = s.t.UpdateTodo(ut.ID, ut.Text, ut.IsDone)
	if err != nil {
		s.l.Fatal(err)
	}
}

func (s server) deleteTodo(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		s.l.Fatal(err)
	}

	ut := &models.UserTodo{}
	err = json.Unmarshal(b, &envelope{ut})
	if err != nil {
		s.l.Fatal(err)
	}

	err = s.t.DeleteTodo(ut.ID, s.u.GetUsername(r))
	if err != nil {
		s.l.Fatal(err)
	}
}

func (s server) lastIDTodo(w http.ResponseWriter, r *http.Request) {
	id, err := s.t.GetLastID()
	if err != nil {
		s.l.Fatal(err)
	}

	_, err = w.Write([]byte(strconv.Itoa(id)))
	if err != nil {
		s.l.Fatal(err)
	}
}

func (s server) signupUserForm(w http.ResponseWriter, r *http.Request) {
	err := s.u.SignupUserForm(w, r)
	if err != nil {
		s.l.Fatal(err)
	}
}

func (s server) logoutUser(w http.ResponseWriter, r *http.Request) {
	s.u.LogoutUser(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther) //NOTE: JS ignores this
}

func (s server) signupUser(w http.ResponseWriter, r *http.Request) {
	v, err := getFormValues(r, "password", "username")
	if err != nil {
		s.l.Error(err)
		http.Redirect(w, r, "/user/signup", http.StatusSeeOther)
	}

	err = s.u.SignupUser(w, r, v["username"], v["password"])
	if err != nil {
		s.l.Error(err)
		http.Redirect(w, r, "/user/signup", http.StatusSeeOther)
	}

	if err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (s server) loginUserForm(w http.ResponseWriter, r *http.Request) {
	err := s.u.LoginUserForm(w, r)
	if err != nil {
		s.l.Fatal(err)
	}
}

func (s server) loginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Trying to login")
	v, err := getFormValues(r, "password", "username")
	if err != nil {
		s.l.Error(err)
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}

	err = s.u.LoginUser(w, r, v["username"], v["password"])
	if err != nil {
		s.l.Error(err)
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}

	if err == nil {
		s.l.Info("Authenticated!\n")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (s server) authUser(w http.ResponseWriter, r *http.Request) {
	if !s.u.IsAuthenticated(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type userNameData struct {
	UserName string `json:"username"`
}

func (s server) getUsername(w http.ResponseWriter, r *http.Request) {
	name := s.u.GetUsername(r)

	un := userNameData{UserName: name}

	data, err := json.Marshal(&envelope{un})
	if err != nil {
		s.l.Fatal(err)
	}

	_, err = w.Write(data)
	if err != nil {
		s.l.Fatal(err)
	}
}

// func (s server) getUsername(w http.ResponseWriter, r *http.Request) {
// 	_, err := w.Write([]byte(s.u.GetUsername(r)))
// 	if err != nil {
// 		s.l.Error(err)
// 	}
// }
