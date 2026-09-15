// package scheduler
//
// import (
// 	"github.com/robfig/cron/v3"
// 	"log"
// 	"strconv"
// 	"time"
// 	"vm/internal/domain/kafka/entity"
// 	kafkaRepo "vm/internal/domain/kafka/services"
// 	vmRepo "vm/internal/domain/virtualMachine/repo"
// )
//
// // VMScheduler 定义调度器结构
// type VMScheduler struct {
// 	cron *cron.Cron
// }
//
// // NewVMScheduler 创建调度器实例
// func NewVMScheduler() *VMScheduler {
// 	return &VMScheduler{
// 		// vmRepo:      vmRepo,
// 		// kafkaClient: kafkaClient,
// 		cron: cron.New(),
// 	}
// }
//
// // Start 启动调度器
// func (s *VMScheduler) Start(vmRepo vmRepo.IVirtualMachineRepository, kafkaClient kafkaRepo.IKafkaService) {
// 	// 每天早上 8 点执行
// 	_, err := s.cron.AddFunc("0 8 * * *", func() {
// 		log.Println("Running cron job at 8:00 AM")
// 		now := time.Now()
//
// 		// 1. 获取提前两天需要通知的虚拟机（两天后过期）
// 		vmsToNotify, err := vmRepo.GetExpiringVMs(2)
// 		if err != nil {
// 			log.Printf("Failed to get VMs expiring in 2 days: %v", err)
// 			return
// 		}
//
// 		// 2. 获取当天需要销毁的虚拟机（今天过期）
// 		vmsToDestroy, err := vmRepo.GetExpiringVMs(1)
// 		if err != nil {
// 			log.Printf("Failed to get VMs expiring today: %v", err)
// 			return
// 		}
//
// 		// 处理通知任务（提前2天通知）
// 		for _, vm := range vmsToNotify {
// 			// 计算通知时间（这里直接使用当前时间，因为我们已经在两天前通知）
// 			taskNotify := &entity.Task{
// 				VMID:     strconv.Itoa(int(vm.ID)),
// 				ExpireAt: vm.ExpirationTime.Unix(), // 记录实际过期时间
// 				TaskType: "notify",
// 			}
// 			if err := kafkaClient.ProduceTask(taskNotify, "vm_notify"); err != nil {
// 				log.Printf("Failed to produce notify task for VM %d: %v", vm.ID, err)
// 			} else {
// 				log.Printf("Sent notification for VM %d, expires at %v", vm.ID, vm.ExpirationTime)
// 			}
// 		}
//
// 		// 处理销毁任务（当天过期）
// 		for _, vm := range vmsToDestroy {
// 			// 如果虚拟机已过期或当天即将过期，立即安排销毁
// 			if vm.ExpirationTime.Before(now) || vm.ExpirationTime.Equal(now) {
// 				taskDestroy := &entity.Task{
// 					VMID:     strconv.Itoa(int(vm.ID)),
// 					ExpireAt: vm.ExpirationTime.Unix(),
// 					TaskType: "destroy",
// 				}
// 				if err := kafkaClient.ProduceTask(taskDestroy, "vm_destroy"); err != nil {
// 					log.Printf("Failed to produce destroy task for VM %d: %v", vm.ID, err)
// 				} else {
// 					log.Printf("Scheduled destruction for VM %d, expired at %v", vm.ID, vm.ExpirationTime)
// 				}
// 			}
// 		}
// 	})
// 	if err != nil {
// 		log.Fatalf("Failed to add cron job: %v", err)
// 	}
//
// 	// 启动 cron 调度器
// 	s.cron.Start()
// }

package scheduler

import (
	"github.com/robfig/cron/v3"
	"log"
)

// VMScheduler 定义调度器结构
type VMScheduler struct {
	cron *cron.Cron
}

// NewVMScheduler 创建调度器实例
func NewVMScheduler() *VMScheduler {
	return &VMScheduler{
		cron: cron.New(),
	}
}

// RegisterTask 注册定时任务，接受一个无参数的函数作为任务
func (s *VMScheduler) RegisterTask(taskFunc func()) {
	_, err := s.cron.AddFunc("0 8 * * *", func() {
		log.Println("Running cron job at 8:00 AM")
		taskFunc() // 执行注入的任务函数
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}
}

// Start 启动调度器
func (s *VMScheduler) Start() {
	s.cron.Start()
}

// Stop 停止调度器（可选，用于优雅关闭）
func (s *VMScheduler) Stop() {
	s.cron.Stop()
}
