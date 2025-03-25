package repo

import (
	"vm/internal/application/pve"
	"vm/internal/domain/kafka/entity"
	"vm/internal/domain/virtualMachine/repo"
	kafkaClient "vm/internal/infrastructure/kafka"
)

// IKafkaRepository kafka接口
type IKafkaRepository interface {
	// ProduceTask 生产任务
	ProduceTask(task *entity.Task, topic string) error
	// ConsumeTasks 消费任务
	ConsumeTasks(repo repo.IVirtualMachineRepository, pveServer pve.IPVEServer)
	// Close 释放资源
	Close()
}

// kafka接口实现检查
var _ IKafkaRepository = (*kafkaClient.KafkaClient)(nil)
