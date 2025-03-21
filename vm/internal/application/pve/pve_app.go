package pve

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	vmProto "github.com/world-fish/proto/vm"
	"log"
	pveDomainServices "vm/internal/domain/pve/services"
	"vm/internal/domain/virtualMachine/repo"
	"vm/internal/infrastructure/utils"
)

/*
对pve进行操作(本地服务器虚拟机里的创建、删除...)
*/

type IPVEServer interface {
	CreateVM(ctx context.Context, req *vmProto.CreateVMReq) (*vmProto.CreateVMResp, error)
	DestroyVM(ctx context.Context, req *vmProto.DestroyVMReq) (*vmProto.DestroyVMResp, error)
	StartVM(ctx context.Context, req *vmProto.StartVMReq) (*vmProto.StartVMResp, error)
	StopVM(ctx context.Context, req *vmProto.StopVMReq) (*vmProto.StopVMResp, error)
	RenewVM(ctx context.Context, req *vmProto.RenewVMReq) (*vmProto.RenewVMResp, error)
	GetVMInfo(ctx context.Context, req *vmProto.GetVMInfoReq) (*vmProto.GetVMInfoResp, error)
	SetConn(conn *websocket.Conn)
	ReceiveMessage() ([]byte, error)
}

var _ IPVEServer = &PVEServer{}

type PVEServer struct {
	pveService pveDomainServices.IPVEService
	vmRepo     repo.IVirtualMachineRepository
}

func NewPVEServer(pveService pveDomainServices.IPVEService, vmService repo.IVirtualMachineRepository) *PVEServer {
	return &PVEServer{pveService: pveService, vmRepo: vmService}
}

func (P PVEServer) ReceiveMessage() ([]byte, error) {
	message, err := P.pveService.ReceiveMessage()
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (P PVEServer) SetConn(conn *websocket.Conn) {
	P.pveService.SetConn(conn)
	return
}

func (P PVEServer) CreateVM(ctx context.Context, req *vmProto.CreateVMReq) (resp *vmProto.CreateVMResp, err error) {
	// 解析参数
	email := ctx.Value("email").(string)
	// 验证邮箱格式
	ok := utils.VerifyEmail(email)
	if !ok {
		log.Printf("邮箱格式错误")
		return nil, errors.New("邮箱格式错误")
	}
	// 验证邮箱是否存在
	P.vmRepo.VerifyEmail(email)

	// 构建请求消息
	createVMReq := &CreateVMReq{
		Email: email,
	}
	// 序列化请求消息
	b, err := json.Marshal(createVMReq)
	if err != nil {
		log.Printf("json 序列化失败: %v", err)
		return nil, err
	}
	// 封装请求消息
	data := Data{
		Type: CreateType,
		Data: b,
	}
	// 再次序列化
	b, err = json.Marshal(data)
	if err != nil {
		log.Printf("json 序列化失败: %v", err)
		return nil, err
	}

	// 发送websocket消息
	err = P.pveService.SendMessage(b)
	if err != nil {
		return nil, err
	}

	// 返回结果
	return
}

func (P PVEServer) DestroyVM(ctx context.Context, req *vmProto.DestroyVMReq) (*vmProto.DestroyVMResp, error) {
	//TODO implement me
	panic("implement me")
}

func (P PVEServer) StartVM(ctx context.Context, req *vmProto.StartVMReq) (*vmProto.StartVMResp, error) {
	//TODO implement me
	panic("implement me")
}

func (P PVEServer) StopVM(ctx context.Context, req *vmProto.StopVMReq) (*vmProto.StopVMResp, error) {
	//TODO implement me
	panic("implement me")
}

func (P PVEServer) RenewVM(ctx context.Context, req *vmProto.RenewVMReq) (*vmProto.RenewVMResp, error) {
	//TODO implement me
	panic("implement me")
}

func (P PVEServer) GetVMInfo(ctx context.Context, req *vmProto.GetVMInfoReq) (*vmProto.GetVMInfoResp, error) {
	//TODO implement me
	panic("implement me")
}
