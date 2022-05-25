package models

type UserTodo struct {
	Text        string `json:"text"`
	CreatedDate string `json:"-"`
	UserID      int    `json:"-"`
	ID          int    `json:"id"`
	IsDone      bool   `json:"is_done"`
}
