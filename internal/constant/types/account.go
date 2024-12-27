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

type AccountDeleteResponse struct {
	AccountId string `json:"account_id"`
}

type AccountGetInfoResponse struct {
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

type AccountLoginRequest struct {
	ClientIp   string `json:"client_ip"`
	ClientAddr string `json:"client_addr"`
	ServerAddr string `json:"server_addr"`
	DeviceId   string `json:"device_id"`
	Os         string `json:"os"`
	System     string `json:"system"`
}

type AccountLoginResponse struct {
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

type AccountLogoutRequest struct {
	OnlineId int     `json:"online_id"`
	ClientIp *string `json:"client_ip"`
	State    int8    `json:"state"`
}

type AccountGenerateTokenResponse struct {
	AccountId   string `json:"account_id"`
	Token       string `json:"token"`
	TokenExpire int64  `json:"token_expire"`
}
