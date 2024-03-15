package model

import "time"

type App struct {
	BaseModel

	Name      string     `gorm:"column:name;uniqueIndex:uk_key_name" json:"name"` // 应用名称
	Key       int        `gorm:"column:key;uniqueIndex:uk_key_name" json:"key"`   // 应用KEY
	Secret    string     `gorm:"column:secret" json:"secret"`                     // 应用密钥
	Status    int8       `gorm:"column:status;default:0" json:"status"`           // 应用状态 0停用 1启用 2冻结
	ExpiredAt *time.Time `gorm:"column:expired_at" json:"expired_at"`             // 过期时间
}
