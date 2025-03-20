package grpc

import (
	"context"
	vmProto "github.com/world-fish/proto/vm"
	"log"
	"vm/internal/application/pve"
	"vm/internal/application/virtualMachine"
)

type GRPCServer struct {
	// grpc的服务端实现
	vmProto.UnimplementedVMManagerServer
	pveServer pve.IPVEServer
	vmServer  virtualMachine.IVMServer
}

var _ vmProto.VMManagerServer = (*GRPCServer)(nil)

func NewGRPCServer(PVEServer pve.IPVEServer) *GRPCServer {
	return &GRPCServer{pveServer: PVEServer}

}

func (G GRPCServer) CreateVM(ctx context.Context, req *vmProto.CreateVMReq) (*vmProto.CreateVMResp, error) {
	// 验证参数
	log.Printf("siwu发送创建虚拟机请求")

	// 调用业务逻辑
	resp, err := G.pveServer.CreateVM(ctx, req)
	if err != nil {
		// 返回错误
		return nil, err
	}
	// 返回结果
	return resp, nil

}

func (G GRPCServer) DestroyVM(ctx context.Context, req *vmProto.DestroyVMReq) (*vmProto.DestroyVMResp, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) StartVM(ctx context.Context, req *vmProto.StartVMReq) (*vmProto.StartVMResp, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) StopVM(ctx context.Context, req *vmProto.StopVMReq) (*vmProto.StopVMResp, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) RenewVM(ctx context.Context, req *vmProto.RenewVMReq) (*vmProto.RenewVMResp, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) GetVMInfo(ctx context.Context, req *vmProto.GetVMInfoReq) (*vmProto.GetVMInfoResp, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) mustEmbedUnimplementedVMManagerServer() {
	//TODO implement me
	panic("implement me")
}
