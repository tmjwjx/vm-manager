package repository

import (
	"vm/internal/domain/virtualMachine/entity"
)

// IVirtualMachineRepository 虚拟机仓储接口
type IVirtualMachineRepository interface {
	// Save 保存虚拟机信息
	Save(vm *entity.VirtualMachine) error

	// FindByUUID 根据UUID查找虚拟机
	FindByUUID(uuid string) (*entity.VirtualMachine, error)

	// Delete 删除虚拟机记录
	Delete(uuid string) error

	// Update 更新虚拟机信息
	Update(vm *entity.VirtualMachine) error

	// FindAll 查询所有虚拟机
	FindAll() ([]*entity.VirtualMachine, error)

	// FindByUserEmail 查询用户的所有虚拟机
	FindByUserEmail(email string) ([]*entity.VirtualMachine, error)
}
