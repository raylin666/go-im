package model

import (
	"github.com/go-kratos/kratos/v2/errors"
	"mt/internal/constant/defined"
	"time"
)

const (
	AppStatusClose  = 0
	AppStatusOpen   = 1
	AppStatusFreeze = 2
)

type App struct {
	BaseModel

	Ident     string    `gorm:"column:ident;uniqueIndex:uk_ident_name" json:"ident"` // 唯一标识, 用来标识来源
	Name      string    `gorm:"column:name;uniqueIndex:uk_ident_name" json:"name"`   // 应用名称
	Key       uint64    `gorm:"column:key;unique" json:"key"`                        // 应用KEY
	Secret    string    `gorm:"column:secret" json:"secret"`                         // 应用密钥
	Status    int8      `gorm:"column:status;default:0" json:"status"`               // 应用状态 0停用 1启用 2冻结
	ExpiredAt time.Time `gorm:"column:expired_at" json:"expired_at"`                 // 过期时间
}

// AppAvailableByKeyAndSecret 应用是否可用
func AppAvailableByKeyAndSecret(model *App) *errors.Error {
	if model == nil {
		return defined.ErrorNotVisitAuth
	}

	// 应用状态已关闭或已冻结
	if model.Status != AppStatusOpen {
		return defined.ErrorAppAuthClose
	}

	// 应用已过期
	if model.ExpiredAt.Before(time.Now()) {
		return defined.ErrorAppAuthExpired
	}

	return nil
}
