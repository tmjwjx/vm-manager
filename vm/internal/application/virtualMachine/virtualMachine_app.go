package virtualMachine

import (
	"encoding/json"
	"log"
	"vm/internal/domain/virtualMachine/entity"
	"vm/internal/domain/virtualMachine/repo"
)

/*
根据pve发送的消息 对vm进行操作(数据库里的创建、删除...)
*/

type IVMServer interface {
	ProcessMessage(message []byte)
}

var _ IVMServer = (*VMServer)(nil)

type VMServer struct {
	//VMService services.IVMService
	vmRepo repo.IVirtualMachineRepository
}

func NewVMServer(vmRepo repo.IVirtualMachineRepository) *VMServer {
	return &VMServer{vmRepo: vmRepo}
}

func (V *VMServer) ProcessMessage(message []byte) {
	var data Data
	err := json.Unmarshal(message, &data)
	if err != nil {
		log.Printf("解析数据失败: %v", err)
		return
	}
	switch data.Type {
	case CreateType:
		// 创建虚拟机
		rep := CreateResp{}
		_ = json.Unmarshal(data.Data, &rep)

		vm := &entity.VirtualMachine{
			VMID:  rep.VMID,
			Email: rep.Email,
		}

		// 执行持久化操作
		log.Printf("存储虚拟机信息: %v", vm)
		err := V.vmRepo.CreateVM(vm)
		if err != nil {
			log.Printf("存储虚拟机信息失败: %v", err)
			return
		}
	case DestroyType:
		// 删除虚拟机
	}
}
