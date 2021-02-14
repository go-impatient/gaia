package tpl

import (
	"time"
)

type UserListRequestData struct {
	PageNumber int `form:"page_number" json:"page_number"`
	PageSize   int `form:"page_size" json:"page_size"`
}

// UserBody ...
type UserBody struct {
	Username string `form:"username" json:"username" binding:"required"` // 用户名
	Password string `form:"password" json:"password" binding:"required"` // 密码
}

type UserInfo struct {
	Username              string    `json:"username"`                 // 用户名
	Realname              string    `json:"realname"`                 // 真实姓名
	Nickname              string    `json:"nickname"`                 // 昵称
	Email                 string    `json:"email"`                    // 邮箱
	Phone                 string    `json:"phone"`                    // 手机号码
	Sex                   int       `json:"sex"`                      // 性别: 0 男性, 1 女性, 2未知
	Status                int       `json:"status"`                   // 状态: 1启用, 0禁用
	IsDeleted             int       `json:"is_deleted"`               // 是否已删除 : 1删除, 0未删除
	CreateUser            string    `json:"create_user"`              // 创建人
	UpdateUser            string    `json:"update_user"`              // 修改人
	PasswordErrorLastTime time.Time `json:"password_error_last_time"` // 最后一次输错密码时间
	PasswordErrorNum      int       `json:"password_error_num"`       // 密码错误次数
	PasswordExpireTime    time.Time `json:"password_expire_time"`     // 密码过期时间
}

type UsersResponse struct {
	SuccessResponseType
	Data []UserInfo `json:"data"`
}

type UserResponse struct {
	SuccessResponseType
	Data UserInfo `json:"data"`
}
