package types

import "time"

type AccountCreateRequest struct {
	AccountId string `json:"account_id"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	IsAdmin   bool   `json:"is_admin"`
}

type AccountCreateResponse struct {
	AccountId string    `json:"account_id"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

type AccountCreateData struct {
	AccountId string `json:"account_id"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	IsAdmin   int8   `json:"is_admin"`
}

type AccountUpdateRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	IsAdmin  bool   `json:"is_admin"`
}

type AccountUpdateResponse struct {
	AccountId string    `json:"account_id"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
}

type AccountUpdateData struct {
	AccountId string `json:"account_id"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	IsAdmin   int8   `json:"is_admin"`
}

type AccountDeleteResponse struct {
	AccountId string `json:"account_id"`
}

type AccountGenerateTokenResponse struct {
	AccountId   string `json:"account_id"`
	Token       string `json:"token"`
	TokenExpire int64  `json:"token_expire"`
}
