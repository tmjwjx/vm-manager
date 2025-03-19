package main

import (
	"github.com/gorilla/websocket"
	vmProto "github.com/world-fish/proto/vm"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	pveAppServices "vm/internal/application/pve"
	vmAppServices "vm/internal/application/virtualMachine"
	pveDomainServices "vm/internal/domain/pve/services"
	vmDominServices "vm/internal/domain/virtualMachine/services"
	"vm/internal/infrastructure/inits"
	"vm/internal/infrastructure/persistence/mysql"
	grpcInterface "vm/internal/interfaces/grpc"
	"vm/internal/interfaces/middleware"
	wsInterface "vm/internal/interfaces/websocket"
)

func main() {
	
	// 配置文件初始化
	inits.ConfigInit()
	// 日志配置
	inits.LogInit()
	// 数据库连接
	db := inits.DbInit()
	// 数据库表初始化
	inits.TableInit(db)
	// redis连接
	rdb := inits.RedisInit()
	
	// 创建websocket连接指针
	var wsConn *websocket.Conn
	
	// websocket服务 和 grpc服务
	w := wsInterface.NewWSServer(&wsConn, vmAppServices.NewVMServer(vmDominServices.NewVMService(mysql.NewVmRepo(db))))
	g := grpcInterface.NewGRPCServer(pveAppServices.NewPVEServer(pveDomainServices.NewPVEService(&wsConn)))
	
	// 开启http服务
	go func() {
		http.Handle("/", w)
		_ = http.ListenAndServe(":8088", nil)
	}()
	
	// 开启gRPC服务
	// 开启端口监听
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Printf("监听失败: %v", err)
	}
	// 注册grpc服务
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(middleware.JWTInterceptor(rdb)))
	vmProto.RegisterVMManagerServer(grpcServer, g)
	// 启动服务
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Printf("启动失败: %v", err)
	}
	
}
