package repository

import (
	"context"
	"github.com/go-impatient/gaia/internal/schema"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// UserRepository interface
type UserRepository interface {
	Exist(ctx context.Context, model *schema.User) (bool, error)
	List(ctx context.Context, limit, page int, sort string, model *schema.User) (total int, users []*schema.User, err error)
	Get(ctx context.Context, id int64) (*schema.User, error)
	Create(ctx context.Context, model *schema.User) (*schema.User, error)
	Update(ctx context.Context, model *schema.User) (*schema.User, error)
	DeleteFull(ctx context.Context, model *schema.User) (*schema.User, error)
	Delete(ctx context.Context, id int64) (*schema.User, error)
	Count(ctx context.Context) (int, error)
}

// userRepository struct
type userRepository struct {
	DB *gorm.DB
}

// NewUserRepository returns an instance of `UserRepository`.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (repo *userRepository) Exist(ctx context.Context, model *schema.User) (bool, error) {
	log.Info().Msg("Received UserRepository.Exist")

	return true, nil
}

func (repo *userRepository) List(ctx context.Context, limit, page int, sort string, model *schema.User) (total int, users []*schema.User, err error) {
	log.Info().Msg("Received UserRepository.List")

	return 0, nil, nil
}

func (repo *userRepository) Get(ctx context.Context, id int64) (*schema.User, error) {
	log.Info().Msg("Received UserRepository.Get")

	return nil, nil
}

func (repo *userRepository) Create(ctx context.Context, model *schema.User) (*schema.User, error) {
	log.Info().Msg("Received UserRepository.Create")

	return nil, nil
}

func (repo *userRepository) Update(ctx context.Context, model *schema.User) (*schema.User, error) {
	log.Info().Msg("Received UserRepository.Update")

	return nil, nil
}

func (repo *userRepository) DeleteFull(ctx context.Context, model *schema.User) (*schema.User, error) {
	log.Info().Msg("Received UserRepository.DeleteFull")

	return nil, nil
}

func (repo *userRepository) Delete(ctx context.Context, id int64) (*schema.User, error) {
	log.Info().Msg("Received UserRepository.Delete")

	return nil, nil
}

func (repo *userRepository) Count(ctx context.Context) (int, error) {
	log.Info().Msg("Received UserRepository.Count")

	return 0, nil
}
