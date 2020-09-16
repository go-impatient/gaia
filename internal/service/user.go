package service

import (
	"context"
	"github.com/go-impatient/gaia/internal/model"
	"github.com/go-impatient/gaia/internal/repository"
	"github.com/go-impatient/gaia/internal/schema"
	"github.com/rs/zerolog/log"
)

type UserService interface {
	Login(ctx context.Context, email, password string) (*model.UserRequest, error)
	Register(ctx context.Context, email, password, phoneNumber string) (*model.UserRequest, error)
	ChangePassword(ctx context.Context, email, password string) error
	GetUserProfile(ctx context.Context, email string) (*model.UserRequest, error)

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

func (srv *userService) Login(ctx context.Context, email, password string) (*model.UserRequest, error) {
	log.Info().Msg("Received UserService.Login")
	_, err := srv.repo.Get(ctx, 0)
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	return nil, nil
}

func (srv *userService) Register(ctx context.Context, email, password, phoneNumber string) (*model.UserRequest, error) {
	log.Info().Msg("Received UserService.Register")
	srv.repo.Create(ctx, &schema.User{
		Email:    "",
		Password: "",
	})

	return nil, nil
}

func (srv *userService) ChangePassword(ctx context.Context, email, password string) error {
	log.Info().Msg("Received UserService.Register")

	return nil
}

func (srv *userService) GetUserProfile(ctx context.Context, email string) (*model.UserRequest, error) {
	log.Info().Msg("Received UserService.Register")

	return nil, nil
}

func (srv *userService) GetRepo() repository.UserRepository {
	return srv.repo
}
