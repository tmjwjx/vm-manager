package application

import (
	vmEntity "vm/internal/domain/virtualMachine/entity"
	"vm/internal/domain/virtualMachine/vo"
)

// IVMApplicationService
// @Description: 虚拟机应用服务接口
// @Author tianjiajie 2025-03-12 14:30:12
type IVMApplicationService interface {
	// CreateVMForUser 为用户创建虚拟机
	CreateVMForUser(userEmail string, spec *vo.VMSpec) (*vmEntity.VirtualMachine, error)

	// DestroyUserVM 销毁用户的虚拟机
	DestroyUserVM(userEmail string, vmUUID string) error

	// PowerOnUserVM 开启用户的虚拟机
	PowerOnUserVM(userEmail string, vmUUID string) error

	// PowerOffUserVM 关闭用户的虚拟机
	PowerOffUserVM(userEmail string, vmUUID string) error

	// GetUserVMStatus 获取用户虚拟机状态
	GetUserVMStatus(userEmail string, vmUUID string) (*vmEntity.VirtualMachine, error)

	// ListAllVMs 管理员获取所有虚拟机列表
	ListAllVMs() ([]*vmEntity.VirtualMachine, error)
}
