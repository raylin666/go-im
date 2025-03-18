package types

type MessageSendC2CMessageRequest struct {
	Seq         string `json:"seq"`
	FromAccount string `json:"from_account"`
	ToAccount   string `json:"to_account"`
	Message     string `json:"message"`
}

type MessageSendC2CMessageResponse struct {
}
