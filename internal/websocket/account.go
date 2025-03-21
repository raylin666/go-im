package websocket

type Account struct {
	ID          string `json:"id"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
	IsAdmin     bool   `json:"is_admin"`     // 是否管理员
	OnlineId    int    `json:"online_id"`    // 账号在线ID, account_online 表
	LogoutState int8   `json:"logout_state"` // 帐号登出状态 0:正常退出 1:超时退出 2:服务端退出
}

// NewAccount 创建账号
func NewAccount(id, nickname, avatar string, onlineId int, isAdmin bool) *Account {
	return &Account{
		ID:       id,
		Nickname: nickname,
		Avatar:   avatar,
		IsAdmin:  isAdmin,
		OnlineId: onlineId,
	}
}

// WithLogoutState 设置帐号登出状态
func (account *Account) WithLogoutState(value int8) *Account {
	account.LogoutState = value
	return account
}
