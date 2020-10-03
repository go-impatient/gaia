package model

import (
	"gorm.io/gorm"
)

type UserToken struct {
	gorm.Model
	UserID     int64 `gorm:"column:user_id;type:bigint(64);not null;default:0" json:"user_id"`         // 对应应用的主键
	Status     int   `gorm:"index:idx_user_status;column:status;type:int(2);default:1" json:"status"`  // 状态: 1启用, 0禁用
	ExpiredAt  int64 `gorm:"column:expired_at;type:bigint(64);not null;default:0" json:"expired_at"`   // 过期时间
	CreateTime int64 `gorm:"column:create_time;type:bigint(64);not null;default:0" json:"create_time"` // 创建时间
}

// TableName, 获取UserToken表名称
func (m *UserToken) TableName() string {
	return "sys_user_token"
}
