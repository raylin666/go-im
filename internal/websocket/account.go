package websocket

type Account struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

func NewAccount(id, nickname, avatar string) *Account {
	return &Account{
		ID:       id,
		Nickname: nickname,
		Avatar:   avatar,
	}
}
