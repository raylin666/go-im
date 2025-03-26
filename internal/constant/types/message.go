package types

import "mt/internal/repositories/dbrepo/model"

type MessageSendC2CMessageRequest struct {
	Seq         string `json:"seq"`
	FromAccount string `json:"from_account"`
	ToAccount   string `json:"to_account"`
	Message     string `json:"message"`
}

type MessageSendC2CMessageResponse struct {
}

type MessageSendC2CMessageDataResult struct {
	C2CMessage        *model.C2CMessage
	C2COfflineMessage *model.C2COfflineMessage
	FromAccount       *model.Account
	ToAccount         *model.Account
	Message           string
	ToAccountOnline   bool
	Error             error
}
