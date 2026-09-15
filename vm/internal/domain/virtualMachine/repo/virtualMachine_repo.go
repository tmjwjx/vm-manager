package repo

import (
	"vm/internal/domain/virtualMachine/entity"
	"vm/internal/infrastructure/persistence/mysql"
)

// IVirtualMachineRepository 虚拟机仓储接口
type IVirtualMachineRepository interface {
	// CreateVM 创建虚拟机
	CreateVM(vm *entity.VirtualMachine) error
	// DestroyVM 销毁虚拟机
	DestroyVM(vmId string) error
	// RenewVM 虚拟机续期
	RenewVM(email string, day int) error
	// GetVMInfoByEmail 通过邮箱获取虚拟机信息
	GetVMInfoByEmail(email string) (*entity.VirtualMachine, error)
	// VerifyEmail 验证邮箱是否存在
	VerifyEmail(email string) bool
	// GetExpiringVMs 获取要在days天后要过期的数据，并按过期时间排序
	GetExpiringVMs(days int) ([]*entity.VirtualMachine, error)
	// FindByID 根据ID查找虚拟机
	FindByID(id string) *entity.VirtualMachine
}

var _ IVirtualMachineRepository = (*mysql.VMRepo)(nil)
