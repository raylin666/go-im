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

	FromAccount   string     `gorm:"column:from_account;type:string;size:30;index:un_from_to_account;comment:发送者ID" json:"from_account"`
	ToAccount     string     `gorm:"column:to_account;type:string;size:30;index:un_from_to_account;comment:接收者ID" json:"to_account"`
	MsgType       int8       `gorm:"column:msg_type;default:3;comment:消息类型 目前只支持自定义消息。1:文本 2:图片 3:自定义" json:"msg_type"`
	Data          string     `gorm:"column:data;type:string;comment:消息内容" json:"data"`
	Status        int8       `gorm:"column:status;default:1;comment:消息状态 0:隐藏 1:显示" json:"status"`
	IsRevoke      int8       `gorm:"column:is_revoke;default:0;comment:消息是否已撤回 0:否 1:是" json:"is_revoke"`
	RevokedAt     *time.Time `gorm:"column:revoked_at;comment:消息撤回时间" json:"revoked_at"`
	SendAt        time.Time  `gorm:"column:send_at;comment:消息发送时间" json:"send_at"`
	FromDeletedAt *time.Time `gorm:"column:from_deleted_at;comment:发送者删除消息时间" json:"from_deleted_at"`
	ToDeletedAt   *time.Time `gorm:"column:to_deleted_at;comment:接收者删除消息时间" json:"to_deleted_at"`
}
