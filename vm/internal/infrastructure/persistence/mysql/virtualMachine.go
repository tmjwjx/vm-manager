package mysql

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"time"
	"vm/internal/domain/virtualMachine/entity"
)

type VMRepo struct {
	db *gorm.DB
}

func NewVmRepo(db *gorm.DB) *VMRepo {
	return &VMRepo{db: db}
}

func (vmRepo *VMRepo) CreateVM(vm *entity.VirtualMachine) error {
	// 开启事务
	tx := vmRepo.db.Begin()
	// 创建虚拟机
	if err := tx.Create(vm).Error; err != nil {
		log.Printf("创建虚拟机失败: %v", err)
		tx.Rollback()
		return err
	}
	// 提交
	tx.Commit()
	return nil
}

func (vmRepo *VMRepo) DestroyVM(vm *entity.VirtualMachine) error {
	//TODO implement me
	panic("implement me")
}

func (vmRepo *VMRepo) RenewVM(email string, day int) error {
	// 开启事务
	tx := vmRepo.db.Begin()
	// 查询虚拟机
	var vm entity.VirtualMachine
	if err := tx.Where("email = ?", email).First(&vm).Error; err != nil {
		log.Printf("查询虚拟机失败: %v", err)
		tx.Rollback()
		return err
	}
	// 续期
	if vm.ExpirationTime != nil {
		newTime := vm.ExpirationTime.AddDate(0, 0, day)
		vm.ExpirationTime = &newTime
	} else {
		// 处理ExpirationTime为nil的情况
		now := time.Now()
		newTime := now.AddDate(0, 0, day)
		vm.ExpirationTime = &newTime
	}
	// 更新
	if err := tx.Save(&vm).Error; err != nil {
		log.Printf("续期虚拟机失败: %v", err)
		tx.Rollback()
		return err
	}
	// 提交
	tx.Commit()
	return nil
}

func (vmRepo *VMRepo) GetVMInfoByEmail(email string) (vm *entity.VirtualMachine, err error) {
	// 查询虚拟机
	if err := vmRepo.db.Where("email = ?", email).First(vm).Error; err != nil {
		log.Printf("查询虚拟机失败: %v", err)
		return nil, err
	}
	return vm, nil
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
	if err := vmRepo.db.Where("email = ?", email).First(&user).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("查询邮箱失败: %v", err)
			return false
		}
	}
	// 判断是否存在
	if user.ID == 0 {
		return false
	} else {
		return true
	}
}
