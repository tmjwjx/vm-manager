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
	DestroyVM(vm *entity.VirtualMachine) error
	// RenewVM 虚拟机续期
	RenewVM(email string, day int) error
	// GetVMInfo 获取虚拟机信息
	GetVMInfo(vm *entity.VirtualMachine) error
	// VerifyEmail 验证邮箱是否存在
	VerifyEmail(email string) bool
}

var _ IVirtualMachineRepository = (*mysql.VMRepo)(nil)
