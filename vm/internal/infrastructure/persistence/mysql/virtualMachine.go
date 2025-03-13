package mysql

import "vm/internal/domain/virtualMachine/entity"

type VmRepo struct{}

func (vmRepo *VmRepo) Create(vm *entity.VirtualMachine) error {
	return nil
}
