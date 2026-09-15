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

func (vmRepo *VMRepo) DestroyVM(vmId string) error {
	// 开启事务
	tx := vmRepo.db.Begin()

	// 删除虚拟机
	if err := tx.Where("vm_id = ?", vmId).Delete(&entity.VirtualMachine{}).Error; err != nil {
		log.Printf("删除虚拟机失败: %v", err)
		tx.Rollback()
		return err
	}

	// 提交
	tx.Commit()
	return nil

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

// GetExpiringVMs 获取即将过期的数据，支持指定天数范围，并按过期时间排序
func (vmRepo *VMRepo) GetExpiringVMs(days int) ([]*entity.VirtualMachine, error) {
	var vms []*entity.VirtualMachine
	now := time.Now()
	// 计算查询的结束时间，根据传入的天数
	end := now.Add(time.Duration(days*24) * time.Hour)

	// 查询条件：过期时间在当前时间之后且在指定天数之内
	err := vmRepo.db.
		Where("expire_at > ? AND expire_at <= ?", now, end).
		Order("expire_at ASC"). // 按过期时间升序排序
		Find(&vms).Error

	if err != nil {
		log.Printf("查询虚拟机错误（%d天内到期）：%v\n", days, err)
		return nil, err
	}

	return vms, nil
}

// FindByID 根据ID查找虚拟机
func (vmRepo *VMRepo) FindByID(id string) *entity.VirtualMachine {
	var vm entity.VirtualMachine
	err := vmRepo.db.Where("vm_id = ?", id).First(&vm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		log.Printf("查找虚拟机失败: %v", err)
		return nil
	}
	return &vm
}
