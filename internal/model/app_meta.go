package model

// AppMeta App扩展表
type AppMeta struct {
	MetaID    int64  `gorm:"primary_key;column:meta_id;type:bigint(64);not null;default:0" json:"-"` // 扩展表的主键
	AppID     int64  `gorm:"column:app_id;type:bigint(64);not null;default:0" json:"app_id"`         // 对应应用的主键
	MetaKey   string `gorm:"column:meta_key;type:varchar(255);default:''" json:"meta_key"`
	MetaValue string `gorm:"column:meta_value;type:longtext;default:''" json:"meta_value"`
}

// TableName, 获取App扩展表名称
func (m *AppMeta) TableName() string {
	return "sys_app_meta"
}
