package service

import (
	"context"
	"github.com/go-impatient/gaia/internal/model"
	"github.com/go-impatient/gaia/internal/model/tpl"

	"github.com/rs/zerolog/log"

	"github.com/go-impatient/gaia/internal/repository"
)

type UserService interface {
	Login(ctx context.Context, email, password string) (*tpl.UserRequest, error)
	Register(ctx context.Context, email, password, phoneNumber string) (*tpl.UserRequest, error)
	ChangePassword(ctx context.Context, email, password string) error
	GetUserProfile(ctx context.Context, email string) (*tpl.UserRequest, error)

	GetRepo() repository.UserRepository
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{
		repo: r,
	}
}

func (srv *userService) Login(ctx context.Context, email, password string) (*tpl.UserRequest, error) {
	log.Info().Msg("Received UserService.Login")
	_, err := srv.repo.Get(ctx, 0)
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	return nil, nil
}

func (srv *userService) Register(ctx context.Context, email, password, phoneNumber string) (*tpl.UserRequest, error) {
	log.Info().Msg("Received UserService.Register")
	srv.repo.Create(ctx, &model.User{
		Email:    "moocss@163.com",
		Password: "123456",
	})

	return nil, nil
}

func (srv *userService) ChangePassword(ctx context.Context, email, password string) error {
	log.Info().Msg("Received UserService.Register")

	return nil
}

func (srv *userService) GetUserProfile(ctx context.Context, email string) (*tpl.UserRequest, error) {
	log.Info().Msg("Received UserService.Register")

	return nil, nil
}

func (srv *userService) GetRepo() repository.UserRepository {
	return srv.repo
}
