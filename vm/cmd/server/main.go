package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"vm/internal/application/service"
	"vm/internal/infrastructure/init"
	pvm "vm/internal/interfaces/grpc/proto/vm"
	ws "vm/internal/interfaces/websocket"
)

func main() {
	
	// 初始化配置
	init.Init()
	
	// 启动grpc服务
	go GRPCInit()
	
	// 启动websocket服务
	WebSocketInit()
}

func WebSocketInit() {
	http.HandleFunc("/", ws.WS)
	_ = http.ListenAndServe(":8088", nil)
}

func GRPCInit() {
	// 开启端口监听
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Printf("监听失败: %v", err)
	}
	
	//// 初始化PVE客户端
	//pveClient := entity.NewPVEClient(globals.Conn)
	
	// 注册grpc服务
	grpcServer := grpc.NewServer()
	pvm.RegisterVMManagerServer(grpcServer, service.NewVMServer())
	
	// 启动服务
	if err = grpcServer.Serve(listen); err != nil {
		log.Printf("启动服务失败: %v", err)
	}
}
