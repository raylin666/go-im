package types

type MessageSendC2CMessageRequest struct {
	Seq       string `json:"seq"`
	ToAccount string `json:"to_account"`
	Message   string `json:"message"`
}

type MessageSendC2CMessageResponse struct {
}
