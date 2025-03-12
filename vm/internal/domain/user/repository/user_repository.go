package repository

import (
	"vm/internal/domain/user/entity"
)

// IUserRepository 用户仓储接口
type IUserRepository interface {
	// Save 保存用户信息
	Save(user *entity.User) error

	// FindByEmail 根据邮箱查找用户
	FindByEmail(email string) (*entity.User, error)

	// FindAll 查询所有用户
	FindAll() ([]*entity.User, error)
}
