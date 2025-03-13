package repository

import (
	"vm/internal/domain/virtualMachine/entity"
)

// IVirtualMachineRepository 虚拟机仓储接口
type IVirtualMachineRepository interface {
	// Create 创建虚拟机
	Create(vm *entity.VirtualMachine) error
}
