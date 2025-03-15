package service

import (
	"encoding/json"
	"log"
	"vm/internal/domain/virtualMachine/entity"
	"vm/internal/infrastructure/globals"
	"vm/internal/infrastructure/persistence/mysql"
)

type CreateRep struct {
	VMID  string `json:"vm_id"`
	Email string `json:"email"`
}

func ServerPVE(mes []byte) {
	var data globals.Data
	err := json.Unmarshal(mes, &data)
	if err != nil {
		log.Printf("解析数据失败: %v", err)
		return
	}
	switch data.Type {
	case globals.CreateType:
		// 创建虚拟机
		rep := CreateRep{}
		_ = json.Unmarshal(data.Data, &rep)
		
		vm := &entity.VirtualMachine{
			VMID:  rep.VMID,
			Email: rep.Email,
		}
		
		// 执行持久化操作
		log.Printf("存储虚拟机信息: %v", vm)
		err := mysql.NewVmRepo(globals.DB).CreateVM(vm)
		if err != nil {
			log.Printf("存储虚拟机信息失败: %v", err)
			return
		}
	case globals.DestroyType:
		// 删除虚拟机
		
	}
}
