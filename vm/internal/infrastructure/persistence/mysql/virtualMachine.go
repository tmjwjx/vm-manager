package mysql

import "vm/internal/domain/virtualMachine/entity"

type VmRepo struct{}

func (vmRepo *VmRepo) CreateVM(vm *entity.VirtualMachine) error {
	return nil
}

// DestroyVM 销毁虚拟机
func (vmRepo *VmRepo) DestroyVM(vm *entity.VirtualMachine) error {

	return nil
}
