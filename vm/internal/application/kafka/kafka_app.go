package kafka

import (
	"log"
	"strconv"
	"time"
	"vm/internal/domain/kafka/entity"
	kafkaDomainServices "vm/internal/domain/kafka/services"
	"vm/internal/domain/virtualMachine/repo"
	vmSchedulerDomainServices "vm/internal/domain/vmScheduler/services"
)

// kafka服务

type IKafkaServer interface {
	// ProduceTask 添加任务
	ProduceTask()
	// TaskFunc 过期数据逻辑
	TaskFunc()
	// StartConsuming 消费任务
	StartConsuming(handler func(string, string)) // 消费任务并调用回调函数

}

var _ IKafkaServer = &KafkaServer{}

// KafkaServer kafka服务
type KafkaServer struct {
	// pveService   pveDomainServices.IPVEService
	vmRepo       repo.IVirtualMachineRepository
	kafkaService kafkaDomainServices.IKafkaService
	vmScheduler  vmSchedulerDomainServices.IVMSchedulerService
}

// func NewKafkaServer(pveService pveDomainServices.IPVEService, vmRepo repo.IVirtualMachineRepository, kafkaService kafkaDomainServices.IKafkaService, vmScheduler vmSchedulerDomainServices.IVMSchedulerService) *KafkaServer {
// 	return &KafkaServer{
// 		pveService:   pveService,
// 		vmRepo:       vmRepo,
// 		kafkaService: kafkaService,
// 		vmScheduler:  vmScheduler,
// 	}
// }

func NewKafkaServer(vmRepo repo.IVirtualMachineRepository, kafkaService kafkaDomainServices.IKafkaService, vmScheduler vmSchedulerDomainServices.IVMSchedulerService) *KafkaServer {
	return &KafkaServer{
		vmRepo:       vmRepo,
		kafkaService: kafkaService,
		vmScheduler:  vmScheduler,
	}
}

// ProduceTask 添加任务
func (k *KafkaServer) ProduceTask() {
	// 定时器，定时为kafka添加任务：
	// 注册任务
	k.vmScheduler.RegisterTask(k.TaskFunc)
	// 开始定时任务
	k.vmScheduler.Start()
}

// TaskFunc 过期数据逻辑
func (k *KafkaServer) TaskFunc() {
	now := time.Now()

	// 获取需要通知和销毁的虚拟机
	vmsToNotify, err := k.vmRepo.GetExpiringVMs(2) // 两天后过期
	if err != nil {
		log.Printf("Failed to get VMs expiring in 2 days: %v", err)
		return
	}

	vmsToDestroy, err := k.vmRepo.GetExpiringVMs(1) // 今天过期
	if err != nil {
		log.Printf("Failed to get VMs expiring today: %v", err)
		return
	}

	// 处理通知任务
	for _, vm := range vmsToNotify {
		taskNotify := &entity.Task{
			VMID:     strconv.Itoa(int(vm.ID)),
			ExpireAt: vm.ExpirationTime.Unix(),
			TaskType: "notify",
		}
		if err := k.kafkaService.ProduceTask(taskNotify, "vm_notify"); err != nil {
			log.Printf("Failed to produce notify task for VM %d: %v", vm.ID, err)
		} else {
			log.Printf("Sent notification for VM %d, expires at %v", vm.ID, vm.ExpirationTime)
		}
	}
	// 处理销毁任务
	for _, vm := range vmsToDestroy {
		if vm.ExpirationTime.Before(now) || vm.ExpirationTime.Equal(now) {
			taskDestroy := &entity.Task{
				VMID:     strconv.Itoa(int(vm.ID)),
				ExpireAt: vm.ExpirationTime.Unix(),
				TaskType: "destroy",
			}
			if err := k.kafkaService.ProduceTask(taskDestroy, "vm_destroy"); err != nil {
				log.Printf("Failed to produce destroy task for VM %d: %v", vm.ID, err)
			} else {
				log.Printf("Scheduled destruction for VM %d, expired at %v", vm.ID, vm.ExpirationTime)
			}
		}
	}
}

// StartConsuming 开始消费 Kafka 任务，并将任务交给回调函数处理
func (k *KafkaServer) StartConsuming(handler func(taskType, vmID string)) {
	taskChan, err := k.kafkaService.ConsumeTasks()
	if err != nil {
		log.Fatalf("Failed to start consuming Kafka tasks: %v", err)
	}

	go func() {
		for taskInfo := range taskChan {
			handler(taskInfo.TaskType, taskInfo.VMID)
		}
	}()
}

// // StartConsuming 开始消费 Kafka 任务
// func (k *KafkaServer) StartConsuming() {
// 	taskChan, err := k.kafkaService.ConsumeTasks()
// 	if err != nil {
// 		log.Fatalf("Failed to start consuming tasks: %v", err)
// 	}
//
// 	go func() {
// 		for taskInfo := range taskChan {
// 			k.handleTask(taskInfo.TaskType, taskInfo.VMID)
// 		}
// 	}()
// }
//
// // handleTask 处理消费到的任务
// func (k *KafkaServer) handleTask(taskType, vmID string) {
// 	log.Printf("Received task - Type: %s, VMID: %s", taskType, vmID)
//
// 	switch taskType {
// 	case "notify":
//
// 		if err := c.pveClient.NotifyExpiration(vmID); err != nil {
// 			log.Printf("Failed to notify expiration for VM %s: %v", vmID, err)
// 		} else {
// 			log.Printf("Notified expiration for VM: %s", vmID)
// 		}
// 	case "destroy":
// 		vm, err := c.vmRepo.FindByID(vmID)
// 		if err != nil || vm == nil {
// 			log.Printf("VM %s not found, skipping", vmID)
// 			return
// 		}
// 		now := time.Now()
// 		if vm.ExpirationTime.After(now) {
// 			log.Printf("VM %s has been renewed, skipping destroy", vmID)
// 			return
// 		}
// 		if err := c.pveClient.DestroyVM(vmID); err != nil {
// 			log.Printf("Failed to destroy VM %s: %v", vmID, err)
// 		} else {
// 			if err := c.vmRepo.DestroyVM(vmID); err != nil {
// 				log.Printf("Failed to delete VM %s from repo: %v", vmID, err)
// 			} else {
// 				log.Printf("Destroyed VM: %s", vmID)
// 			}
// 		}
// 	default:
// 		log.Printf("Unknown task type: %s for VM %s", taskType, vmID)
// 	}
// }
