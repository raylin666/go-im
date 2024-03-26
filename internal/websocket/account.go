package websocket

import "fmt"

type AccountLogin struct {
	UserId string
	Client *Client
}

func NewAccountLogin(userId string, client *Client) *AccountLogin {
	return &AccountLogin{
		UserId: userId,
		Client: client,
	}
}

func (account *AccountLogin) ManagerKey() string {
	return fmt.Sprintf("%d_%s", account.Client.AppKey, account.UserId)
}
