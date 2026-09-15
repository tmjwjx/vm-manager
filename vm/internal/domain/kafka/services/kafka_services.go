package services

import (
	"vm/internal/domain/kafka/entity"
	kafkaClient "vm/internal/infrastructure/kafka"
)

// IKafkaService kafka接口
type IKafkaService interface {
	// ProduceTask 生产任务
	ProduceTask(task *entity.Task, topic string) error
	// ConsumeTasks 消费任务
	ConsumeTasks() (<-chan struct {
		TaskType string
		VMID     string
	}, error)
	// Close 释放资源
	Close()
}

// kafka接口实现检查
var _ IKafkaService = (*kafkaClient.KafkaClient)(nil)
