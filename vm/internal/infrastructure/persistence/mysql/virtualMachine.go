package mysql

import (
	"gorm.io/gorm"
	"vm/internal/domain/virtualMachine/entity"
)

type VMRepo struct {
	DB *gorm.DB
}

func NewVmRepo(DB *gorm.DB) *VMRepo {
	return &VMRepo{DB: DB}
}

func (vmRepo *VMRepo) CreateVM(vm *entity.VirtualMachine) error {
	//TODO implement me
	panic("implement me")
}

func (vmRepo *VMRepo) DestroyVM(vm *entity.VirtualMachine) error {
	//TODO implement me
	panic("implement me")
}

func (vmRepo *VMRepo) RenewVM(vm *entity.VirtualMachine) error {
	//TODO implement me
	panic("implement me")
}

func (vmRepo *VMRepo) GetVMInfo(vm *entity.VirtualMachine) error {
	//TODO implement me
	panic("implement me")
}

// VerifyEmail
// @Description: 验证邮箱是否存在
// @receiver     vmRepo
// @param        email string
// @return       err
// @Author tianjiajie 2025-03-16 16:41:42
func (vmRepo *VMRepo) VerifyEmail(email string) (b bool) {
	// 查询邮箱是否存在
	var user entity.VirtualMachine
	if err := vmRepo.DB.Where("email = ?", email).First(&user).Error; err != nil {
		vmRepo.DB.Rollback()
		return
	}
	// 判断是否存在
	if user.ID == 0 {
		return false
	} else {
		return true
	}
}
