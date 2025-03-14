package repository

import (
	"vm/internal/domain/virtualMachine/entity"
)

// IVirtualMachineRepository 虚拟机仓储接口
type IVirtualMachineRepository interface {
	// CreateVM 创建虚拟机
	CreateVM(vm *entity.VirtualMachine) error
	// DestroyVM 销毁虚拟机
	DestroyVM(vm *entity.VirtualMachine) error
}
