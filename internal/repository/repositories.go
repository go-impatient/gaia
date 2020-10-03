package repository

import (
	"gorm.io/gorm"
)

type Repositories struct {
	User  UserRepository
	App AppRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		App: NewAppRepository(db),
		User:  NewUserRepository(db),
	}
}
