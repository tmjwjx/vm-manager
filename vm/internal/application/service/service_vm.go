package service

import (
	"context"
	"vm/internal/domain/virtualMachine/service"
	"vm/internal/infrastructure/globals"
	pvm "vm/internal/interfaces/grpc/proto/vm"
)

type VMServer struct {
	pvm.UnimplementedVMManagerServer
	//pveClient entity.IPVEClient
}

func NewVMServer() *VMServer {
	return &VMServer{
		//pveClient: pveClient,
	}
}

func (s *VMServer) CreateVM(ctx context.Context, req *pvm.CreateVMReq) (*pvm.CreateVMResp, error) {

	// 使用PVEClient创建虚拟机
	// 这里使用了一些默认配置，实际应用中可以从请求参数或配置中获取
	pveClient := service.NewPVEClient(globals.Conn)
	err := pveClient.CreateVM(req)
	if err != nil {
		return nil, err
	}

	return &pvm.CreateVMResp{}, nil
}
