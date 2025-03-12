package service

import (
	"vm/internal/domain/virtualMachine/entity"
	"vm/internal/domain/virtualMachine/vo"
)

// IVirtualMachineService 虚拟机领域服务接口
type IVirtualMachineService interface {
	// CreateVM 创建虚拟机
	CreateVM(spec *vo.VMSpec) (*entity.VirtualMachine, error)

	// DestroyVM 销毁虚拟机
	DestroyVM(vmUUID string) error

	// PowerOnVM 开机
	PowerOnVM(vmUUID string) error

	// PowerOffVM 关机
	PowerOffVM(vmUUID string) error

	// GetVMStatus 获取虚拟机状态
	GetVMStatus(vmUUID string) (*entity.VirtualMachine, error)

	// ListVMs 获取虚拟机列表
	ListVMs() ([]*entity.VirtualMachine, error)
}
