package model

import (
	"errors"

	"github.com/asaskevich/govalidator"
)

// UserRequest 请求数据结果
type UserRequest struct {
	ID       int64  `json:"-" valid:"-"`
	Username string `json:"username" valid:"-"`
	Password string `json:"password" valid:"-"`
	Email    string `json:"email" valid:"email"`
}

// UserResponse 返回数据结构
type UserResponse struct {
}

// Validate the fields.
func (u *UserRequest) Validate() error {
	_, err := govalidator.ValidateStruct(u)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
