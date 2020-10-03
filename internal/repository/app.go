package repository

import (
	"context"
	"github.com/go-impatient/gaia/internal/model"

	"gorm.io/gorm"
)

// AdminRepository interface
type AdminRepository interface {
	// Create User 创建用户
	Save(ctx context.Context, model *model.Admin) (bool, error)

	// Update User 更新用户
	Update(ctx context.Context, id uint64, data map[string]interface{}) (bool, error)

	// Delete User 删除用户
	Delete(ctx context.Context, id uint64) (bool, error)

	// Find User 根据用户ID, 获取用户数据
	FindOne(ctx context.Context, id uint64) (*model.Admin, error)

	// Find User 根据多个用户ID, 获取用户数据
	FindAll(ctx context.Context, ids []uint64) (*model.Admin, error)
}

// adminRepository struct
type adminRepository struct {
	DB *gorm.DB
}

// NewAdminRepository returns an instance of `AdminRepository`.
func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{
		DB: db,
	}
}

func (repo *adminRepository) Save(ctx context.Context, model *model.Admin) (bool, error) {
	panic("implement me")
}

func (repo *adminRepository) Update(ctx context.Context, id uint64, model map[string]interface{}) (bool, error) {
	panic("implement me")
}

func (repo *adminRepository) Delete(ctx context.Context, id uint64) (bool, error) {
	panic("implement me")
}

func (repo *adminRepository) FindOne(ctx context.Context, id uint64) (*model.Admin, error) {
	panic("implement me")
}

func (repo *adminRepository) FindAll(ctx context.Context, ids []uint64) (*model.Admin, error) {
	panic("implement me")
}
