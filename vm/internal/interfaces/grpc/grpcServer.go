package grpc

import (
	"context"
	vmProto "github.com/world-fish/proto/vm"
	"log"
	"vm/internal/application/pve"
	"vm/internal/application/virtualMachine"
)

var _ vmProto.VMManagerServer = (*GRPCServer)(nil)

type GRPCServer struct {
	// grpc的服务端实现
	vmProto.UnimplementedVMManagerServer
	pveServer pve.IPVEServer
	vmServer  virtualMachine.IVMServer
}

func NewGRPCServer(pveServer pve.IPVEServer, vmServer virtualMachine.IVMServer) *GRPCServer {
	return &GRPCServer{pveServer: pveServer, vmServer: vmServer}
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
	vm, err := G.pveServer.StartVM(ctx, req)
	if err != nil {
		return nil, err
	}
	return vm, nil
}

func (G GRPCServer) StopVM(ctx context.Context, req *vmProto.StopVMReq) (*vmProto.StopVMResp, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) RenewVM(ctx context.Context, req *vmProto.RenewVMReq) (*vmProto.RenewVMResp, error) {
	// 参数验证
	log.Printf("siwu发送续费虚拟机请求")
	// 调用业务逻辑
	vm, err := G.vmServer.RenewVM(ctx, req)
	if err != nil {
		return nil, err
	}
	// 返回结果
	return vm, nil
}

func (G GRPCServer) GetVMInfo(ctx context.Context, req *vmProto.GetVMInfoReq) (*vmProto.GetVMInfoResp, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) mustEmbedUnimplementedVMManagerServer() {
	//TODO implement me
	panic("implement me")
}
