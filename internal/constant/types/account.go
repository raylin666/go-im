package types

import "time"

type AccountCreateRequest struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Status   int8   `json:"status"`
	IsAdmin  bool   `json:"is_admin"`
}

type AccountCreateResponse struct {
	UserId    string    `json:"user_id"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
	Status    int8      `json:"status"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

type AccountCreateData struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Status   int8   `json:"status"`
	IsAdmin  int8   `json:"is_admin"`
}
