package types

import "time"

type AccountCreateRequest struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	IsAdmin  bool   `json:"is_admin"`
}

type AccountCreateResponse struct {
	UserId    string    `json:"user_id"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

type AccountCreateData struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	IsAdmin  int8   `json:"is_admin"`
}
