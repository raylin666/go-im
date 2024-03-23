package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int            `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`               // 主键ID
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime;comment:更新时间" json:"updated_at"` // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;comment:删除时间" json:"deleted_at"`                // 删除时间
}
