package account

import (
	"time"
)

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

type DeleteResponse struct {
	AccountId string `json:"account_id"`
}

type GetInfoResponse struct {
	AccountId      string     `json:"account_id"`
	Nickname       string     `json:"nickname"`
	Avatar         string     `json:"avatar"`
	IsAdmin        bool       `json:"is_admin"`
	IsOnline       bool       `json:"is_online"`
	LastLoginIp    string     `json:"last_login_ip"`
	FirstLoginTime *time.Time `json:"first_login_time"`
	LastLoginTime  *time.Time `json:"last_login_time"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

type LoginRequest struct {
	ClientIp string `json:"client_ip"`
}

type LoginResponse struct {
	AccountId      string     `json:"account_id"`
	Nickname       string     `json:"nickname"`
	Avatar         string     `json:"avatar"`
	IsAdmin        bool       `json:"is_admin"`
	IsOnline       bool       `json:"is_online"`
	LastLoginIp    string     `json:"last_login_ip"`
	FirstLoginTime *time.Time `json:"first_login_time"`
	LastLoginTime  *time.Time `json:"last_login_time"`
	OnlineId       int        `json:"online_id"`
}

type GenerateTokenResponse struct {
	AccountId   string `json:"account_id"`
	Token       string `json:"token"`
	TokenExpire int64  `json:"token_expire"`
}
