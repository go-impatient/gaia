package repository

import (
	"context"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/go-impatient/gaia/internal/model"
)

// UserRepository interface
type UserRepository interface {
	Exist(ctx context.Context, model *model.User) (bool, error)
	List(ctx context.Context, limit, page int, sort string, model *model.User) (total int, users []*model.User, err error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Create(ctx context.Context, model *model.User) (*model.User, error)
	Update(ctx context.Context, model *model.User) (*model.User, error)
	DeleteFull(ctx context.Context, model *model.User) (*model.User, error)
	Delete(ctx context.Context, id int64) (*model.User, error)
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

func (repo *userRepository) Exist(ctx context.Context, model *model.User) (bool, error) {
	log.Info().Msg("Received UserRepository.Exist")

	return true, nil
}

func (repo *userRepository) List(ctx context.Context, limit, page int, sort string, model *model.User) (total int, users []*model.User, err error) {
	log.Info().Msg("Received UserRepository.List")

	return 0, nil, nil
}

func (repo *userRepository) Get(ctx context.Context, id int64) (*model.User, error) {
	log.Info().Msg("Received UserRepository.Get")

	return nil, nil
}

func (repo *userRepository) Create(ctx context.Context, model *model.User) (*model.User, error) {
	log.Info().Msg("Received UserRepository.Create")

	ok := Create(model)

	if ok {
		log.Info().Msg("添加成功")
	}

	return nil, nil
}

func (repo *userRepository) Update(ctx context.Context, model *model.User) (*model.User, error) {
	log.Info().Msg("Received UserRepository.Update")

	return nil, nil
}

func (repo *userRepository) DeleteFull(ctx context.Context, model *model.User) (*model.User, error) {
	log.Info().Msg("Received UserRepository.DeleteFull")

	return nil, nil
}

func (repo *userRepository) Delete(ctx context.Context, id int64) (*model.User, error) {
	log.Info().Msg("Received UserRepository.Delete")

	return nil, nil
}

func (repo *userRepository) Count(ctx context.Context) (int, error) {
	log.Info().Msg("Received UserRepository.Count")

	return 0, nil
}
