package pve

import (
	"context"
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
}

type PVEServer struct {
	PVEService pveDomainServices.IPVEService
}

var _ IPVEServer = &PVEServer{}
var _ IPVEServer = (*PVEServer)(nil)

func NewPVEServer(PVEService pveDomainServices.IPVEService) *PVEServer {
	return &PVEServer{PVEService: PVEService}
}

func (P PVEServer) CreateVM(ctx context.Context, req *vmProto.CreateVMReq) (resp *vmProto.CreateVMResp, err error) {
	// 解析参数
	email := ctx.Value("email").(string)
	// 验证邮箱格式
	// 验证邮箱是否存在

	// 发送websocket消息
	err = P.PVEService.CreateVM(email)
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
