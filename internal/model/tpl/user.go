package tpl

import (
	"time"
)

// UserRequest 请求数据结果
type UserRequest struct {
	Username              string    `gorm:"column:username;type:varchar(50);not null" json:"username"`                     // 用户名
	Password              string    `gorm:"column:password;type:varchar(50);not null" json:"password"`                     // 密码
	Realname              string    `gorm:"column:realname;type:varchar(50)" json:"realname"`                              // 真实姓名
	Nickname              string    `gorm:"column:nickname;type:varchar(50)" json:"nickname"`                              // 昵称
	Email                 string    `gorm:"column:email;type:varchar(80)" json:"email"`                                    // 邮箱
	Phone                 string    `gorm:"column:phone;type:varchar(20);not null" json:"phone"`                           // 手机号码
	Sex                   int       `gorm:"column:sex;type:int(2)" json:"sex"`                                             // 性别: 0 男性, 1 女性, 2未知
	Status                int       `gorm:"column:status;type:int(2)" json:"status"`                                       // 状态: 1启用, 0禁用
	IsDeleted             int       `gorm:"column:is_deleted;type:int(2)" json:"is_deleted"`                               // 是否已删除 : 1删除, 0未删除
	CreateUser            string    `gorm:"column:create_user;type:varchar(64)" json:"create_user"`                        // 创建人
	UpdateUser            string    `gorm:"column:update_user;type:varchar(64)" json:"update_user"`                        // 修改人
	PasswordErrorLastTime time.Time `gorm:"column:password_error_last_time;type:datetime" json:"password_error_last_time"` // 最后一次输错密码时间
	PasswordErrorNum      int       `gorm:"column:password_error_num;type:int(11)" json:"password_error_num"`              // 密码错误次数
	PasswordExpireTime    time.Time `gorm:"column:password_expire_time;type:datetime" json:"password_expire_time"`         // 密码过期时间
}

// UserResponse 返回数据结构
type UserResponse struct {
}
