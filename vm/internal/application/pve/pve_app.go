package pve

import (
	"context"
	"github.com/gorilla/websocket"
	vmProto "github.com/world-fish/proto/vm"
	pveDomainServices "vm/internal/domain/pve/services"
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

type PVEServer struct {
	pveService pveDomainServices.IPVEService
}

func (P PVEServer) ReceiveMessage() ([]byte, error) {
	message, err := P.pveService.ReceiveMessage()
	if err != nil {
		return nil, err
	}
	return message, nil
}

var _ IPVEServer = &PVEServer{}

func NewPVEServer(pveService pveDomainServices.IPVEService) *PVEServer {
	return &PVEServer{pveService: pveService}
}

func (P PVEServer) SetConn(conn *websocket.Conn) {
	P.pveService.SetConn(conn)
	return
}

func (P PVEServer) CreateVM(ctx context.Context, req *vmProto.CreateVMReq) (resp *vmProto.CreateVMResp, err error) {
	// 解析参数
	email := ctx.Value("email").(string)
	// 验证邮箱格式
	// 验证邮箱是否存在

	// 发送websocket消息
	err = P.pveService.CreateVM(email)
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
