package runs

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"vm/internal/application/service"
	pvm "vm/internal/interfaces/grpc/proto/vm"
)

func GRPCRun() {
	// 开启端口监听
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Printf("监听失败: %v", err)
	}

	// 注册grpc服务
	grpcServer := grpc.NewServer()
	pvm.RegisterVMManagerServer(grpcServer, service.NewVMServer())

	// 启动服务
	if err = grpcServer.Serve(listen); err != nil {
		log.Printf("启动服务失败: %v", err)
	}
}
