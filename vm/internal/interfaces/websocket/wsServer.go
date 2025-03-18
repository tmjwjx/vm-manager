package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"vm/internal/application/virtualMachine"
)

type WSServer struct {
	WS       **websocket.Conn
	VMServer virtualMachine.IVMServer
}

func NewWSServer(WS **websocket.Conn, VMServer virtualMachine.IVMServer) *WSServer {
	return &WSServer{WS: WS, VMServer: VMServer}
}

func (W WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// websocket升级器
	var up = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	// 升级连接
	conn, err := up.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	*W.WS = conn

	// 读取消息
	for {
		_, message, err := (*W.WS).ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("收到消息: %s", message)
		// 处理信息
		W.VMServer.ProcessMessage(message)
	}
}
