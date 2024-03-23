package model

import (
	"fmt"
	"time"
)

const (
	AccountStatusOffline = 0
	AccountStatusOnline  = 1
)

type Account struct {
	BaseModel

	UserId         string     `gorm:"column:user_id;type:string;size:30;unique:uk_user;comment:用户ID" json:"user_id"`
	Username       string     `gorm:"column:username;type:string;size:30;comment:用户昵称" json:"username"`
	Avatar         string     `gorm:"column:avatar;type:string;size:255;comment:用户头像" json:"avatar"`
	IsAdmin        int8       `gorm:"column:is_admin;default:0;index:idx_admin;comment:是否管理员 0否 1是" json:"is_admin"`
	Status         int8       `gorm:"column:status;default:0;comment:当前状态: 0离线 1在线" json:"status"`
	LastLoginIp    string     `gorm:"column:last_login_ip;type:string;size:16;comment:最后登录IP" json:"last_login_ip"`
	LastLogoutIp   string     `gorm:"column:last_logout_ip;type:string;size:16;comment:最后登出IP" json:"last_logout_ip"`
	FirstLoginTime *time.Time `gorm:"column:first_login_time;comment:首次登录时间" json:"first_login_time"`
	LastLoginTime  *time.Time `gorm:"column:last_login_time;comment:最后登录时间" json:"last_login_time"`
	LastLogoutTime *time.Time `gorm:"column:last_logout_time;comment:最后登出时间" json:"last_logout_time"`
}

// AccountTableName 获取账号分表名称
func AccountTableName(key uint64) (tableName string) {
	tableName = fmt.Sprintf("account_%d", key)
	return
}

// AccountConvertStatus 转换账号在线状态
func AccountConvertStatus(status int8) string {
	if status == AccountStatusOnline {
		return "Online"
	}

	return "Offline"
}
