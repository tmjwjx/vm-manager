package service

import (
	"context"
	"net/http"
	pvm "vm/internal/interfaces/grpc/proto/vm"
)

type VMServer struct {
	pvm.UnimplementedVMManagerServer
}

func (s *VMServer) CreateVM(ctx context.Context, req *pvm.CreateVMReq) (*pvm.CreateVMResp, error) {

	_, err := http.Get("localhost:8088")
	if err != nil {
		return nil, err
	}

	return &pvm.CreateVMResp{
		VmId: "123",
	}, nil
}
