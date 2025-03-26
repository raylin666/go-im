package model

type C2COfflineMessage struct {
	BaseModel

	FromAccount string `gorm:"column:from_account;type:string;size:30;uniqueIndex:uk_from_to_account;comment:发送者ID" json:"from_account"`
	ToAccount   string `gorm:"column:to_account;type:string;size:30;uniqueIndex:uk_from_to_account;comment:接收者ID" json:"to_account"`
	MessageId   int    `gorm:"column:message_id;default:0;comment:消息ID 当无离线消息时为0" json:"message_id"`
	UnreadNum   int    `gorm:"column:unread_num;default:0;comment:未读消息数量" json:"unread_num"`
}

func (*C2COfflineMessage) TableName() string {
	return "c2c_offline_message"
}
