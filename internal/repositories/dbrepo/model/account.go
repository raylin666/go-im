package model

import "time"

type Account struct {
	BaseModel

	UserId        string     `gorm:"column:user_id;type:string;size:30;unique:uk_user;comment:用户ID" json:"user_id"` // 用户ID
	Username      string     `gorm:"column:username;type:string;size:30;comment:用户昵称" json:"username"`              // 用户昵称
	Avatar        string     `gorm:"column:avatar;type:string;size:255;comment:用户头像" json:"avatar"`                 // 用户头像
	Status        int8       `gorm:"column:status;default:0;comment:用户状态 0停用 1启用" json:"status"`                    // 用户状态 0停用 1启用
	IsAdmin       int8       `gorm:"column:is_admin;default:0;index:idx_admin;comment:是否管理员 0否 1是" json:"is_admin"` // 是否管理员 0否 1是
	LastLoginTime *time.Time `gorm:"column:last_login_time;comment:最后登录时间" json:"last_login_time"`                  // 最后登录时间
}
