package model

import (
	"time"
)

type Account struct {
	BaseModel

	AccountId      string     `gorm:"column:account_id;type:string;size:30;unique:uk_accountid;comment:账号ID" json:"account_id"`
	Nickname       string     `gorm:"column:nickname;type:string;size:30;comment:账号昵称" json:"nickname"`
	Avatar         string     `gorm:"column:avatar;type:string;size:120;comment:账号头像" json:"avatar"`
	IsAdmin        int8       `gorm:"column:is_admin;default:0;index:idx_admin;comment:是否管理员,管理员可向任何账号发送消息 0否 1是" json:"is_admin"`
	IsOnline       int8       `gorm:"column:is_online;default:0;comment:当前状态: 0离线 1在线" json:"is_online"`
	FirstLoginTime *time.Time `gorm:"column:first_login_time;comment:首次登录时间" json:"first_login_time"`
	LastLoginTime  *time.Time `gorm:"column:last_login_time;comment:最后登录时间" json:"last_login_time"`
	LastLoginIp    string     `gorm:"column:last_login_ip;type:string;size:16;comment:最后登录IP" json:"last_login_ip"`
}
