package tpl

import (
	"errors"

	"github.com/asaskevich/govalidator"
)

// AdminRequest 请求数据结果
type AdminRequest struct {
	ID       int64  `json:"-" valid:"-"`
	Username string `json:"username" valid:"-"`
	Password string `json:"password" valid:"-"`
	Email    string `json:"email" valid:"email"`
}

// AdminResponse 返回数据结构
type AdminResponse struct {
}

// Validate the fields.
func (u *AdminRequest) Validate() error {
	_, err := govalidator.ValidateStruct(u)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
