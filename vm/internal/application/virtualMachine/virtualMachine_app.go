package virtualMachine

import (
	"context"
	"encoding/json"
	"errors"
	vmProto "github.com/world-fish/proto/vm"
	"log"
	"time"
	"vm/internal/domain/pve/services"
	"vm/internal/domain/virtualMachine/entity"
	"vm/internal/domain/virtualMachine/repo"
)

/*
根据pve发送的消息 对vm进行操作(数据库里的创建、删除...)
*/

type IVMServer interface {
	ProcessMessage(message []byte)
	RenewVM(ctx context.Context, req *vmProto.RenewVMReq) (*vmProto.RenewVMResp, error)
}

var _ IVMServer = (*VMServer)(nil)

type VMServer struct {
	pveService services.IPVEService
	vmRepo     repo.IVirtualMachineRepository
}

func NewVMServer(pveService services.IPVEService, vmRepo repo.IVirtualMachineRepository) *VMServer {
	return &VMServer{pveService: pveService, vmRepo: vmRepo}
}

func (V *VMServer) RenewVM(ctx context.Context, req *vmProto.RenewVMReq) (resp *vmProto.RenewVMResp, err error) {
	resp = &vmProto.RenewVMResp{}
	// 获取参数
	email := ctx.Value("email").(string)
	day := int(req.Day)
	// 判断虚拟机是否存在
	ok := V.vmRepo.VerifyEmail(email)
	if !ok {
		log.Printf("虚拟机不存在")
		return resp, errors.New("虚拟机不存在")
	}

	// 续期虚拟机
	err = V.vmRepo.RenewVM(email, day)
	if err != nil {
		log.Printf("续期虚拟机失败: %v", err)
		return resp, err
	}
	resp.Result = true

	// 返回结果
	return resp, nil
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

		var now time.Time = time.Now().AddDate(0, 1, 0)
		vm := &entity.VirtualMachine{
			VMID:           rep.VMID,
			Email:          rep.Email,
			IPAddr:         rep.IPAddr,
			ExpirationTime: &now,
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
