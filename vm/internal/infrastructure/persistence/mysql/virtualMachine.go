package mysql

import (
	"gorm.io/gorm"
	"time"
	"vm/internal/domain/virtualMachine/entity"
)

type VmRepo struct {
	db *gorm.DB
}

func NewVmRepo(db *gorm.DB) *VmRepo {
	return &VmRepo{db: db}
}

func (vmRepo *VmRepo) CreateVM(vm *entity.VirtualMachine) (err error) {
	// 开启事务
	tx := vmRepo.db.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return tx.Error
	}

	// 设置过期时间
	end := time.Now().Add(time.Hour * 24 * 30)
	vm.ExpirationTime = &end

	// 存储虚拟机信息
	if err = tx.Create(vm).Error; err != nil {
		tx.Rollback()
		return
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return
	}

	return nil
}

// VerifyEmail
// @Description: 验证邮箱是否存在
// @receiver     vmRepo
// @param        email string
// @return       err
// @Author tianjiajie 2025-03-16 16:41:42
func (vmRepo *VmRepo) VerifyEmail(email string) (b bool, err error) {

	// 查询邮箱是否存在
	var user entity.VirtualMachine
	if err = vmRepo.db.Where("email = ?", email).First(&user).Error; err != nil {
		vmRepo.db.Rollback()
		return
	}
	if user.ID == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

// DestroyVM 销毁虚拟机
func (vmRepo *VmRepo) DestroyVM(vm *entity.VirtualMachine) (err error) {

	return nil
}
