package schema

import "gorm.io/gorm"

// 创建一个Admin数据模型
type Admin struct {
	gorm.Model
	Username string `gorm:"type:varchar(100);column:username;not null" json:"username" valid:"-"`
	Password string `gorm:"type:varchar(50);column:password;not null" json:"password" valid:"-"`
	Email    string `gorm:"type:varchar(100);column:email;unique;not null;" json:"email" valid:"email"`
}

// TableName, 获取User表名称
func (u *Admin) TableName() string {
	return "tb_user"
}
