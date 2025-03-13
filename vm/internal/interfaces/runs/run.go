package runs

func Run() {

	// 启动grpc服务
	go GRPCRun()

	// 启动websocket服务
	WebSocketRun()
}
