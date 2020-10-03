package service

import (
	"github.com/go-impatient/gaia/internal/repository"
)

type Services struct {
	App  AppService
	User UserService
}

func NewServices(repo *repository.Repositories) *Services {
	return &Services{
		App:  NewAppService(repo.App),
		User: NewUserService(repo.User),
	}
}
