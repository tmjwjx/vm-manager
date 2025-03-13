package service

import (
	"context"
	pvm "vm/internal/interfaces/grpc/proto/vm"
	"vm/internal/infrastructure/pve"
)

type VMServer struct {
	pvm.UnimplementedVMManagerServer
	pveClient pve.IPVEClient
}

func NewVMServer(pveClient pve.IPVEClient) *VMServer {
	return &VMServer{
		pveClient: pveClient,
	}
}

func (s *VMServer) CreateVM(ctx context.Context, req *pvm.CreateVMReq) (*pvm.CreateVMResp, error) {
	// 使用PVEClient创建虚拟机
	// 这里使用了一些默认配置，实际应用中可以从请求参数或配置中获取
	vmId, err := s.pveClient.CreateVM(2, 4.0, "linux", "ubuntu-20.04")
	if err != nil {
		return nil, err
	}

	return &pvm.CreateVMResp{
		VmId: vmId,
	}, nil
}
