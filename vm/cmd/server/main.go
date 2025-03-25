package main

import (
	vmProto "github.com/world-fish/proto/vm"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	pveAppServices "vm/internal/application/pve"
	vmAppServices "vm/internal/application/virtualMachine"
	"vm/internal/infrastructure/inits"
	"vm/internal/infrastructure/persistence/mysql"
	websocket2 "vm/internal/infrastructure/websocket"
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

	// 基础设施层：

	// 创建websocket连接指针
	wsConn := websocket2.NewWebSocketClient()
	vmRepo := mysql.NewVmRepo(db)

	// 应用层：
	pve := pveAppServices.NewPVEServer(wsConn, vmRepo)
	vm := vmAppServices.NewVMServer(wsConn, vmRepo)

	// 接口层：

	// grpc服务
	g := grpcInterface.NewGRPCServer(pve, vm)
	// websocket服务
	w := wsInterface.NewWSServer(pve, vm)

	// 开启http服务(websocket)：
	go func() {
		http.Handle("/", w)
		_ = http.ListenAndServe(":8088", nil)
	}()

	// 开启gRPC服务：
	// 开启端口监听
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Printf("监听失败: %v", err)
	}

	// 注册grpc服务：

	// 拦截器：身份验证
	jwt := middleware.JWTInterceptor(rdb)
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(jwt))
	vmProto.RegisterVMManagerServer(grpcServer, g)
	// 启动服务
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Printf("启动失败: %v", err)
	}

}
