package service

import (
	"context"
	pvm "github.com/world-fish/proto/vm"
	"vm/internal/domain/virtualMachine/service"
	"vm/internal/infrastructure/globals"
)

type VMServer struct {
	pvm.UnimplementedVMManagerServer
}

func NewVMServer() *VMServer {
	return &VMServer{}
}

func (s *VMServer) CreateVM(ctx context.Context, req *pvm.CreateVMReq) (*pvm.CreateVMResp, error) {
	
	// 使用PVEClient创建虚拟机
	// 这里使用了一些默认配置，实际应用中可以从请求参数或配置中获取
	pveClient := service.NewPVEClient(globals.Conn)
	err := pveClient.CreateVM(ctx, req)
	if err != nil {
		return nil, err
	}
	
	return &pvm.CreateVMResp{
		Result: true,
	}, nil
}

func (s *VMServer) DestroyVM(ctx context.Context, req *pvm.DestroyVMReq) (*pvm.DestroyVMResp, error) {
	
	// 使用PVEClient创建虚拟机
	// 这里使用了一些默认配置，实际应用中可以从请求参数或配置中获取
	pveClient := service.NewPVEClient(globals.Conn)
	err := pveClient.DestroyVM(ctx, req)
	if err != nil {
		return nil, err
	}
	
	return &pvm.DestroyVMResp{}, nil
}
