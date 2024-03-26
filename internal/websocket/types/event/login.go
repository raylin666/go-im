package event

import "time"

// LoginRequest 事件登录请求
type LoginRequest struct {
	UserId  string `json:"user_id"` // 用户ID
	Usersig string `json:"usersig"` // 用户签名
}

// LoginResponse 事件登录响应
type LoginResponse struct {
	UserId         string    `json:"user_id"`          // 用户ID
	Username       string    `json:"username"`         // 用户名称
	Avatar         string    `json:"avatar"`           // 用户头像
	IsAdmin        bool      `json:"is_admin"`         // 是否管理员
	Status         string    `json:"status"`           // 在线状态 离线(Offline) 在线(Online)
	FirstLoginTime time.Time `json:"first_login_time"` // 用户首次登录时间
	LastLoginTime  time.Time `json:"last_login_time"`  // 用户最后登录时间
	LastLoginIp    string    `json:"last_login_ip"`    // 用户最后登录IP
	RepeatLogin    bool      `json:"repeat_login"`     // 在此之前用户是否已登录状态 (如果用户此前未登录将返回false)
}

// LogoutResponse 事件登出响应
type LogoutResponse struct {
	UserId     string    `json:"user_id"`     // 用户ID
	LogoutTime time.Time `json:"logout_time"` // 用户登出时间
}

// LoginStatusResponse 事件登录状态响应
type LoginStatusResponse struct {
	UserId string `json:"user_id"` // 用户ID, 未登录返回空字符串
	Status string `json:"status"`  // 登录状态 已登录(Login) 未登录(Logout)
}
