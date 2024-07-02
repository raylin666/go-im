package model

import (
	"time"
)

type AccountOnline struct {
	ID         int        `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"` // 主键ID
	AccountId  string     `gorm:"column:account_id;type:string;size:30;unique:uk_accountid;comment:账号ID" json:"account_id"`
	LoginTime  time.Time  `gorm:"column:login_time;not null;comment:登录时间" json:"login_time"`
	LogoutTime *time.Time `gorm:"column:logout_time;comment:登出时间" json:"logout_time"`
	LoginIp    string     `gorm:"column:login_ip;not null;type:string;size:16;comment:登录IP" json:"login_ip"`
	LogoutIp   string     `gorm:"column:logout_ip;type:string;size:16;comment:登出IP" json:"logout_ip"`
	ClientIp   string     `gorm:"column:client_ip;type:string;size:16;comment:客户端连接IP" json:"client_ip"`
	ClientPort int        `gorm:"column:client_port;type:int;default:80;comment:客户端连接端口" json:"client_port"`
	ClientId   int        `gorm:"column:client_id;type:int;default:0;comment:客户端ID" json:"client_id"`
	DeviceId   string     `gorm:"column:device_id;type:string;comment:设备ID" json:"device_id"`
	Os         string     `gorm:"column:os;type:string;default:web;comment:系统类型, 目前有 web|android|ios 值" json:"os"`
	System     string     `gorm:"column:system;type:string;comment:设备信息, 例如: HUAWEI#EML-AL00#HWEML#28#9" json:"system"`
}
