// package kafka
//
// import (
// 	"log"
// 	"time"
// 	"vm/internal/application/pve"
// 	"vm/internal/application/virtualMachine"
// )
//
// type KafkaServer struct {
// 	pveServer pve.IPVEServer
// 	vmServer  virtualMachine.IVMServer
// }
//
// func NewKafkaServer(pveServer pve.IPVEServer, vmServer virtualMachine.IVMServer) *KafkaServer {
// 	return &KafkaServer{pveServer: pveServer, vmServer: vmServer}
// }
//
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

// interface/kafka/kafka_server.go
package kafka

import (
	"encoding/json"
	"log"
	"vm/internal/application/kafka" // 引入应用层的 Kafka 服务
	"vm/internal/application/pve"
	"vm/internal/application/virtualMachine"
)

// KafkaServer 定义接口层的 Kafka 服务
type KafkaServer struct {
	kafkaServer kafka.IKafkaServer       // 应用层的 Kafka 服务
	pveServer   pve.IPVEServer           // PVE 应用层服务
	vmServer    virtualMachine.IVMServer // VM 应用层服务
}

// NewKafkaServer 创建 KafkaServer 实例
func NewKafkaServer(kafkaServer kafka.IKafkaServer, pveServer pve.IPVEServer, vmServer virtualMachine.IVMServer) *KafkaServer {
	return &KafkaServer{
		kafkaServer: kafkaServer,
		pveServer:   pveServer,
		vmServer:    vmServer,
	}
}

// StartConsuming 启动 Kafka 消费任务
func (k *KafkaServer) StartConsuming() {
	// 委托应用层的 Kafka 服务启动消费
	go k.kafkaServer.StartConsuming(k.handleTask)
	log.Println("Kafka consumer started in interface layer")
}

// handleTask 处理任务的回调函数，根据任务类型调用应用层方法
func (k *KafkaServer) handleTask(taskType, vmID string) {
	log.Printf("Received Kafka task - Type: %v, VMID: %v", taskType, vmID)

	switch taskType {
	case "notify":
		// // 调用 PVE 应用层的通知方法
		// if err := k.pveServer.NotifyExpiration(vmID); err != nil {
		// 	log.Printf("Failed to notify expiration for VM %k: %v", vmID, err)
		// } else {
		// 	log.Printf("Notified expiration for VM: %k", vmID)
		// }
	case "destroy":
		var data struct {
			Type string `json:"type"`
			Data []byte `json:"data"`
		}
		data.Type = taskType
		data.Data = []byte(vmID)
		result, _ := json.Marshal(data)
		k.vmServer.ProcessMessage(result)

		// // 调用 VM 应用层的查询方法
		// vm, err := k.vmServer.GetVMByID(vmID)
		// if err != nil || vm == nil {
		// 	log.Printf("VM %k not found, skipping destroy", vmID)
		// 	return
		// }
		//
		// now := time.Now()
		// if vm.ExpirationTime.After(now) {
		// 	log.Printf("VM %k has been renewed, skipping destroy", vmID)
		// 	return
		// }
		//
		// // 调用 PVE 应用层的销毁方法
		// if err := k.pveServer.DestroyVM(vmID); err != nil {
		// 	log.Printf("Failed to destroy VM %k: %v", vmID, err)
		// } else {
		// 	// 如果销毁成功，调用 VM 应用层的删除方法
		// 	if err := k.vmServer.DeleteVM(vmID); err != nil {
		// 		log.Printf("Failed to delete VM %k from repo: %v", vmID, err)
		// 	} else {
		// 		log.Printf("Destroyed VM: %k", vmID)
		// 	}
		// }
	default:
		log.Printf("Unknown task type: %k for VM %k", taskType, vmID)
	}
}
