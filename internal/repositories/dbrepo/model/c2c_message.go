package model

import (
	"time"
)

const (
	C2CMessageStatusOff = 0
	C2CMessageStatusOn  = 1

	C2CMessageRevokeNo  = 0
	C2CMessageRevokeYes = 1
)

type C2CMessage struct {
	BaseModel

	FromAccountId  string     `gorm:"column:from_account_id;type:string;size:30;index:un_from_to_account;comment:发送者账号ID" json:"from_account_id"`
	ToAccountId    string     `gorm:"column:to_account_id;type:string;size:30;index:un_from_to_account;comment:接收者账号ID" json:"to_account_id"`
	MsgType        string     `gorm:"column:msg_type;type:string;size:30;comment:消息类型" json:"msg_type"`
	Data           string     `gorm:"column:data;type:string;size:600;comment:消息内容" json:"data"`
	Status         int8       `gorm:"column:status;default:1;comment:消息状态 0:隐藏 1:显示" json:"status"`
	IsRevoke       int8       `gorm:"column:is_revoke;default:0;comment:是否已撤回 0:否 1:是" json:"is_revoke"`
	RevokeTime     *time.Time `gorm:"column:revoke_time;comment:消息撤回时间" json:"revoke_time"`
	SendTime       time.Time  `gorm:"column:send_time;comment:消息发送时间" json:"send_time"`
	FromDeleteTime *time.Time `gorm:"column:from_delete_time;comment:发送者删除消息时间" json:"from_delete_time"`
	ToDeleteTime   *time.Time `gorm:"column:to_delete_time;comment:接收者删除消息时间" json:"to_delete_time"`
}
