package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/go-impatient/gaia/internal/model"
)

// AppRepository interface
type AppRepository interface {
	// Create User 创建用户
	Save(ctx context.Context, model *model.App) (bool, error)

	// Update User 更新用户
	Update(ctx context.Context, id uint64, data map[string]interface{}) (bool, error)

	// Delete User 删除用户
	Delete(ctx context.Context, id uint64) (bool, error)

	// Find User 根据用户ID, 获取用户数据
	FindOne(ctx context.Context, id uint64) (*model.App, error)

	// Find User 根据多个用户ID, 获取用户数据
	FindAll(ctx context.Context, ids []uint64) (*model.App, error)
}

// appRepository struct
type appRepository struct {
	DB *gorm.DB
}

// NewAppRepository returns an instance of `AppRepository`.
func NewAppRepository(db *gorm.DB) AppRepository {
	return &appRepository{
		DB: db,
	}
}

func (repo *appRepository) Save(ctx context.Context, model *model.App) (bool, error) {
	panic("implement me")
}

func (repo *appRepository) Update(ctx context.Context, id uint64, model map[string]interface{}) (bool, error) {
	panic("implement me")
}

func (repo *appRepository) Delete(ctx context.Context, id uint64) (bool, error) {
	panic("implement me")
}

func (repo *appRepository) FindOne(ctx context.Context, id uint64) (*model.App, error) {
	panic("implement me")
}

func (repo *appRepository) FindAll(ctx context.Context, ids []uint64) (*model.App, error) {
	panic("implement me")
}
