package service

import (
	"context"
	"github.com/go-impatient/gaia/internal/model"
	"github.com/go-impatient/gaia/internal/repository"
)

type AdminService interface {
	Login(ctx context.Context, email, password string) (*model.AdminRequest, error)
	Register(ctx context.Context, email, password, phoneNumber string) (*model.AdminRequest, error)
	ChangePassword(ctx context.Context, email, password string) error
	GetUserProfile(ctx context.Context, email string) (*model.AdminRequest, error)

	GetRepo() repository.AdminRepository
}

type adminService struct {
	repo repository.AdminRepository
}

func NewAdminService(r repository.AdminRepository) AdminService {
	return &adminService{
		repo: r,
	}
}

func (srv *adminService) Login(ctx context.Context, email, password string) (*model.AdminRequest, error) {
	panic("implement me")
}

func (srv *adminService) Register(ctx context.Context, email, password, phoneNumber string) (*model.AdminRequest, error) {
	panic("implement me")
}

func (srv *adminService) ChangePassword(ctx context.Context, email, password string) error {
	panic("implement me")
}

func (srv *adminService) GetUserProfile(ctx context.Context, email string) (*model.AdminRequest, error) {
	panic("implement me")
}

func (srv *adminService) GetRepo() repository.AdminRepository {
	return srv.repo
}
