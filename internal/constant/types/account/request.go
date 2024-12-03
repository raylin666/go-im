package account

import "time"

type CreateRequest struct {
	AccountId string `json:"account_id"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	IsAdmin   bool   `json:"is_admin"`
}

type CreateResponse struct {
	AccountId string    `json:"account_id"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateData struct {
	AccountId string `json:"account_id"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	IsAdmin   int8   `json:"is_admin"`
}

type UpdateRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	IsAdmin  bool   `json:"is_admin"`
}

type UpdateResponse struct {
	AccountId string    `json:"account_id"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateData struct {
	AccountId string `json:"account_id"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	IsAdmin   int8   `json:"is_admin"`
}

type DeleteResponse struct {
	AccountId string `json:"account_id"`
}

type GenerateTokenResponse struct {
	AccountId   string `json:"account_id"`
	Token       string `json:"token"`
	TokenExpire int64  `json:"token_expire"`
}
