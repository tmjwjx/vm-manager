package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"vm/internal/application/pve"
	"vm/internal/application/virtualMachine"
)

type WSServer struct {
	pveServer pve.IPVEServer
	vmServer  virtualMachine.IVMServer
}

func NewWSServer(pveServer pve.IPVEServer, vmServer virtualMachine.IVMServer) *WSServer {
	return &WSServer{pveServer: pveServer, vmServer: vmServer}
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
	// 调用 SetConn 保证grpc服务中的pve底层的conn链接与websocket服务中的conn链接一致
	W.pveServer.SetConn(conn)

	// 读取消息
	for {
		// 读取消息
		message, err := W.pveServer.ReceiveMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("收到消息: %s", message)
		// 处理信息（对数据库进行操作）
		W.vmServer.ProcessMessage(message)
	}
}
