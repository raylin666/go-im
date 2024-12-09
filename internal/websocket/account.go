package websocket

type Account struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	IsAdmin  bool   `json:"is_admin"`  // 是否管理员
	OnlineId int    `json:"online_id"` // 账号在线ID, account_online 表
}

// NewAccount 创建账号
func NewAccount(id, nickname, avatar string, isAdmin bool) *Account {
	return &Account{
		ID:       id,
		Nickname: nickname,
		Avatar:   avatar,
		IsAdmin:  isAdmin,
	}
}

// WithOnlineId 添加账号在线ID
func (account *Account) WithOnlineId(id int) *Account {
	account.OnlineId = id
	return account
}
