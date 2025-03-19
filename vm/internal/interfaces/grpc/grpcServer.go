package grpc

import (
	"context"
	vmProto "github.com/world-fish/proto/vm"
	"log"
	services "vm/internal/application/pve"
)

type GRPCServer struct {
	// grpc的服务端实现
	vmProto.UnimplementedVMManagerServer
	PVEServer services.IPVEServer
}

var _ vmProto.VMManagerServer = (*GRPCServer)(nil)

func NewGRPCServer(PVEServer services.IPVEServer) *GRPCServer {
	return &GRPCServer{PVEServer: PVEServer}

}

func (G GRPCServer) CreateVM(ctx context.Context, req *vmProto.CreateVMReq) (*vmProto.CreateVMResp, error) {
	// 验证参数
	log.Printf("siwu发送创建虚拟机请求")

	// 调用业务逻辑
	//vmServer := services.NewVMServer(G.DB, G.WS)
	resp, err := G.PVEServer.CreateVM(ctx, req)
	if err != nil {
		// 返回错误
		return nil, err
	} else {
		// 返回结果
		return resp, nil
	}
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
