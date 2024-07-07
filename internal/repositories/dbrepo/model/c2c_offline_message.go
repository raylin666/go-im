package model

type C2COfflineMessage struct {
	BaseModel

	FromAccountId string `gorm:"column:from_account_id;type:string;size:30;uniqueIndex:uk_from_to_account;comment:发送者账号ID" json:"from_account_id"`
	ToAccountId   string `gorm:"column:to_account_id;type:string;size:30;uniqueIndex:uk_from_to_account;comment:接收者用户ID" json:"to_account_id"`
	MessageId     int    `gorm:"column:message_id;comment:消息ID, 无离线消息时值为0" json:"message_id"`
}
