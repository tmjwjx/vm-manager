package scheduler

import (
	"github.com/robfig/cron/v3"
	"log"
	"strconv"
	"time"
	kafkaVm "vm/internal/infrastructure/kafka"
	mysqlVm "vm/internal/infrastructure/persistence/mysql"
)

// VMScheduler 定义调度器结构
type VMScheduler struct {
	db          *mysqlVm.VMRepo
	kafkaClient *kafkaVm.KafkaClient
	cron        *cron.Cron
}

// NewVMScheduler 创建调度器实例
func NewVMScheduler(db *mysqlVm.VMRepo, kafkaClient *kafkaVm.KafkaClient) *VMScheduler {
	return &VMScheduler{
		db:          db,
		kafkaClient: kafkaClient,
		cron:        cron.New(),
	}
}

// Start 启动调度器
func (s *VMScheduler) Start() {
	_, err := s.cron.AddFunc("0 8 * * *", func() {
		log.Println("Running cron job at 8:00 AM")
		vms := s.db.GetExpiringVMs()
		now := time.Now()

		for _, vm := range vms {
			notifyTime := time.Unix(vm.ExpirationTime.Unix()-2*3600, 0)
			if notifyTime.After(now) && notifyTime.Before(now.Add(24*time.Hour)) {
				taskNotify := &kafkaVm.Task{
					VMID:     strconv.Itoa(int(vm.ID)),
					ExpireAt: notifyTime.Unix(),
					TaskType: "notify", // 设置任务类型
				}
				if err := s.kafkaClient.ProduceTask(taskNotify, "vm_notify"); err != nil {
					log.Printf("Failed to produce notify task for VM %s: %v", vm.ID, err)
				}
			}

			if vm.ExpirationTime.After(now) && vm.ExpirationTime.Before(now.Add(24*time.Hour)) {
				taskDestroy := &kafkaVm.Task{
					VMID:     strconv.Itoa(int(vm.ID)),
					ExpireAt: vm.ExpirationTime.Unix(),
					TaskType: "destroy", // 设置任务类型
				}
				if err := s.kafkaClient.ProduceTask(taskDestroy, "vm_destroy"); err != nil {
					log.Printf("Failed to produce destroy task for VM %v: %v", vm.ID, err)
				}
			}
		}
	})
	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}
	s.cron.Start()
}
