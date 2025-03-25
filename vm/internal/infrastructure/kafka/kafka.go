package infrastructure

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
	"vm/internal/application/pve"
	"vm/internal/domain/kafka/entity"
	"vm/internal/domain/virtualMachine/repo"
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

// ConsumeTasks 消费任务
func (k *KafkaClient) ConsumeTasks(repo repo.IVirtualMachineRepository, pveServer pve.IPVEServer) {
	go func() {
		for {
			msg, err := k.reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}
			k.handleMessage(msg, repo, pveServer)
		}
	}()
}

// handleMessage 统一处理消息
func (k *KafkaClient) handleMessage(msg kafka.Message, repo repo.IVirtualMachineRepository, pveServer pve.IPVEServer) {
	var task entity.Task
	if err := json.Unmarshal(msg.Value, &task); err != nil {
		log.Printf("Failed to unmarshal task: %v", err)
		return
	}

	vm := repo.FindByID(task.VMID)
	if vm == nil {
		log.Printf("VM %s not found, skipping", task.VMID)
		return
	}

	now := time.Now()
	expireTime := time.Unix(task.ExpireAt, 0)
	diff := expireTime.Sub(now)

	if diff > 24*time.Hour {
		log.Printf("VM %s expire time too far (>24h), likely renewed, discarding", task.VMID)
		return
	}

	if diff > 0 {
		log.Printf("Waiting for %s task of VM %s, sleeping for %v", task.TaskType, task.VMID, diff)
		time.Sleep(diff)
	}

	// switch task.TaskType {
	// case "notify":
	// 	if err := pveServer.NotifyExpiration(task.VMID); err != nil {
	// 		log.Printf("Failed to notify expiration for VM %s: %v", task.VMID, err)
	// 	} else {
	// 		log.Printf("Notified expiration for VM: %s", task.VMID)
	// 	}
	// case "destroy":
	// 	if vm.ExpirationTime.After(now) {
	// 		log.Printf("VM %s has been renewed, skipping destroy", task.VMID)
	// 		return
	// 	}
	// 	if err := pveServer.DestroyVM(task.VMID); err != nil {
	// 		log.Printf("Failed to destroy VM %s: %v from pve", task.VMID, err)
	// 	} else {
	// 		log.Printf("Destroyed VM: %s", task.VMID)
	//
	// 		if err := repo.DestroyVM(task.VMID); err != nil {
	// 			log.Printf("Failed to destroy VM %s: %v from repo", task.VMID, err)
	// 		}
	// 	}
	// default:
	// 	log.Printf("Unknown task type: %s", task.TaskType)
	// }
}

// Close 关闭实例，释放资源
func (k *KafkaClient) Close() {
	k.writer.Close()
	k.reader.Close()
}
