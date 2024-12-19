package model

import (
	"time"
)

const (
	AccountOnlineOsWeb     = "web"
	AccountOnlineOsAndroid = "android"
	AccountOnlineOsiOS     = "ios"
)

type AccountOnline struct {
	ID         int        `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"` // 主键ID
	AccountId  string     `gorm:"column:account_id;type:string;size:30;index:un_accountid_logintime;comment:账号ID" json:"account_id"`
	LoginTime  time.Time  `gorm:"column:login_time;not null;index:un_accountid_logintime;comment:登录时间" json:"login_time"`
	LogoutTime *time.Time `gorm:"column:logout_time;comment:登出时间" json:"logout_time"`
	LoginIp    string     `gorm:"column:login_ip;not null;type:string;size:16;comment:登录IP" json:"login_ip"`
	LogoutIp   string     `gorm:"column:logout_ip;type:string;size:16;comment:登出IP" json:"logout_ip"`
	ClientAddr string     `gorm:"column:client_addr;type:string;size:24;index:un_csaddr;comment:客户端连接本地地址" json:"client_addr"`
	ServerAddr string     `gorm:"column:server_addr;type:string;size:24;index:un_csaddr;comment:服务端连接远程地址" json:"server_addr"`
	DeviceId   string     `gorm:"column:device_id;type:string;comment:设备ID" json:"device_id"`
	Os         string     `gorm:"column:os;type:string;default:web;comment:系统类型, 目前有 web|android|ios 值" json:"os"`
	System     string     `gorm:"column:system;type:string;comment:设备信息, 例如: HUAWEI#EML-AL00#HWEML#28#9" json:"system"`
}
