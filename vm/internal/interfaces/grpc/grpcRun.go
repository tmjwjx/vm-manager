package grpc

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"vm/internal/application/service"
	"vm/internal/infrastructure/globals"
	"vm/internal/interfaces/grpc/interceptor"
	vmProto "vm/internal/interfaces/grpc/proto/vm"
)

func GRPCRun() {
	// 开启端口监听
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Printf("监听失败: %v", err)
	}
	
	// 注册grpc服务
	// 注册 JWT 拦截器
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptor.JWTInterceptor(globals.RDB)))
	//grpcServer := grpc.NewServer()
	
	vmProto.RegisterVMManagerServer(grpcServer, service.NewVMServer())
	
	// 启动服务
	if err = grpcServer.Serve(listen); err != nil {
		log.Printf("启动服务失败: %v", err)
	}
}
