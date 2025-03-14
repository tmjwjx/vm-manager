package main

import (
	_ "vm/internal/infrastructure/inits"
	vmGrpc "vm/internal/interfaces/grpc"
	"vm/internal/interfaces/websocket"
)

func main() {

	// 启动grpc服务
	go vmGrpc.GRPCRun()

	// 启动websocket服务
	websocket.WebSocketRun()

}
