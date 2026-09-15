package services

import (
	"vm/internal/infrastructure/scheduler"
)

// IVMSchedulerService 接口
type IVMSchedulerService interface {
	// Start(vmRepo vmRepo.IVirtualMachineRepository, kafkaClient kafkaRepo.IKafkaService)
	RegisterTask(taskFunc func())
	Start()
	Stop()
}

// 接口实现检查
var _ IVMSchedulerService = (*scheduler.VMScheduler)(nil)
