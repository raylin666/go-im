package websocket

type Account struct {
	ID string `json:"id"`
}

func NewAccount(accountId string) *Account {
	return &Account{
		ID: accountId,
	}
}
