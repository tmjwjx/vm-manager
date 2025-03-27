package infrastructure

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
	"vm/internal/domain/kafka/entity"
)

type KafkaClient struct {
	writer *kafka.Writer // 统一的写入器
	reader *kafka.Reader // 统一的读取器
}

func NewKafkaClient(brokerList []string) (*KafkaClient, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokerList...),
		Balancer: &kafka.LeastBytes{},
	}

	// 使用消费者组消费多个主题
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokerList,
		GroupID:     "vm-consumer-group",
		MinBytes:    10e3,
		MaxBytes:    10e6,
		GroupTopics: []string{"vm_notify", "vm_destroy"},
	})

	return &KafkaClient{
		writer: writer,
		reader: reader,
	}, nil
}

// ProduceTask 生产任务消息
func (k *KafkaClient) ProduceTask(task *entity.Task, topic string) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return k.writer.WriteMessages(context.Background(), kafka.Message{
		Topic: topic, // 在消息中指定目标主题
		Value: data,
	})
}

// // ConsumeTasks 消费任务
// func (k *KafkaClient) ConsumeTasks(services services.IVirtualMachineRepository, pveServer pve.IPVEServer) {
// 	go func() {
// 		for {
// 			msg, err := k.reader.ReadMessage(context.Background())
// 			if err != nil {
// 				log.Printf("Error reading message: %v", err)
// 				continue
// 			}
// 			k.handleMessage(msg, services, pveServer)
// 		}
// 	}()
// }
//
// // handleMessage 统一处理消息
// func (k *KafkaClient) handleMessage(msg kafka.Message, services services.IVirtualMachineRepository, pveServer pve.IPVEServer) {
// 	var task entity.Task
// 	if err := json.Unmarshal(msg.Value, &task); err != nil {
// 		log.Printf("Failed to unmarshal task: %v", err)
// 		return
// 	}
//
// 	vm := services.FindByID(task.VMID)
// 	if vm == nil {
// 		log.Printf("VM %s not found, skipping", task.VMID)
// 		return
// 	}
//
// 	now := time.Now()
// 	expireTime := time.Unix(task.ExpireAt, 0)
// 	diff := expireTime.Sub(now)
//
// 	if diff > 24*time.Hour {
// 		log.Printf("VM %s expire time too far (>24h), likely renewed, discarding", task.VMID)
// 		return
// 	}
//
// 	if diff > 0 {
// 		log.Printf("Waiting for %s task of VM %s, sleeping for %v", task.TaskType, task.VMID, diff)
// 		time.Sleep(diff)
// 	}
//
// 	// switch task.TaskType {
// 	// case "notify":
// 	// 	if err := pveServer.NotifyExpiration(task.VMID); err != nil {
// 	// 		log.Printf("Failed to notify expiration for VM %s: %v", task.VMID, err)
// 	// 	} else {
// 	// 		log.Printf("Notified expiration for VM: %s", task.VMID)
// 	// 	}
// 	// case "destroy":
// 	// 	if vm.ExpirationTime.After(now) {
// 	// 		log.Printf("VM %s has been renewed, skipping destroy", task.VMID)
// 	// 		return
// 	// 	}
// 	// 	if err := pveServer.DestroyVM(task.VMID); err != nil {
// 	// 		log.Printf("Failed to destroy VM %s: %v from pve", task.VMID, err)
// 	// 	} else {
// 	// 		log.Printf("Destroyed VM: %s", task.VMID)
// 	//
// 	// 		if err := services.DestroyVM(task.VMID); err != nil {
// 	// 			log.Printf("Failed to destroy VM %s: %v from services", task.VMID, err)
// 	// 		}
// 	// 	}
// 	// default:
// 	// 	log.Printf("Unknown task type: %s", task.TaskType)
// 	// }
// }

// Close 关闭实例，释放资源
func (k *KafkaClient) Close() {
	k.writer.Close()
	k.reader.Close()
}

// ConsumeTasks 消费任务，返回消息类型和 VMID
func (k *KafkaClient) ConsumeTasks() (<-chan struct {
	TaskType string
	VMID     string
}, error) {
	// 创建一个通道用于返回任务信息
	taskChan := make(chan struct {
		TaskType string
		VMID     string
	})

	go func() {
		defer close(taskChan) // 确保 goroutine 退出时关闭通道

		for {
			msg, err := k.reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}
			taskType, vmID := k.handleMessage(msg)
			if taskType != "" && vmID != "" { // 仅在有效时发送
				taskChan <- struct {
					TaskType string
					VMID     string
				}{TaskType: taskType, VMID: vmID}
			}
		}
	}()

	return taskChan, nil
}

// handleMessage 解析消息并处理休眠逻辑，返回任务类型和 VMID
func (k *KafkaClient) handleMessage(msg kafka.Message) (string, string) {
	var task entity.Task
	if err := json.Unmarshal(msg.Value, &task); err != nil {
		log.Printf("Failed to unmarshal task: %v", err)
		return "", ""
	}

	// // 检查 VM 是否存在
	// vm := repo.FindByID(task.VMID)
	// if vm == nil {
	// 	log.Printf("VM %s not found, skipping", task.VMID)
	// 	return "", "" // 如果 VM 不存在，返回空值
	// }

	// 检查任务时间并休眠
	now := time.Now()
	expireTime := time.Unix(task.ExpireAt, 0)
	diff := expireTime.Sub(now)

	if diff > 24*time.Hour {
		log.Printf("VM %s expire time too far (>24h), likely renewed, discarding", task.VMID)
		return "", "" // 如果时间太远，丢弃任务
	}

	if diff > 0 {
		log.Printf("Waiting for %s task of VM %s, sleeping for %v", task.TaskType, task.VMID, diff)
		time.Sleep(diff) // 休眠直到任务时间到达
	}
	return task.TaskType, task.VMID
}
