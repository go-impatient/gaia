package service

import (
	"context"

	"github.com/go-impatient/gaia/internal/model/tpl"
	"github.com/go-impatient/gaia/internal/repository"
)

type AppService interface {
	Login(ctx context.Context, email, password string) (*tpl.AppRequest, error)
	Register(ctx context.Context, email, password, phoneNumber string) (*tpl.AppRequest, error)
	ChangePassword(ctx context.Context, email, password string) error
	GetUserProfile(ctx context.Context, email string) (*tpl.AppRequest, error)

	GetRepo() repository.AppRepository
}

type appService struct {
	repo repository.AppRepository
}

func NewAppService(r repository.AppRepository) AppService {
	return &appService{
		repo: r,
	}
}

func (srv *appService) Login(ctx context.Context, email, password string) (*tpl.AppRequest, error) {
	panic("implement me")
}

func (srv *appService) Register(ctx context.Context, email, password, phoneNumber string) (*tpl.AppRequest, error) {
	panic("implement me")
}

func (srv *appService) ChangePassword(ctx context.Context, email, password string) error {
	panic("implement me")
}

func (srv *appService) GetUserProfile(ctx context.Context, email string) (*tpl.AppRequest, error) {
	panic("implement me")
}

func (srv *appService) GetRepo() repository.AppRepository {
	return srv.repo
}
