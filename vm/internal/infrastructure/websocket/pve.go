package websocket

import (
	"github.com/gorilla/websocket"
	"log"
)

type PVEService struct {
	ws *websocket.Conn
}

func (P *PVEService) SendMessage(data []byte) (err error) {

	// 发送消息
	err = (P.ws).WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Printf("发送请求失败: %v", err)
		return err
	}
	return
}

func NewWebSocketClient() *PVEService {
	return &PVEService{}
}

func (P *PVEService) ReceiveMessage() (p []byte, err error) {
	// 读取消息
	_, p, err = P.ws.ReadMessage()
	if err != nil {
		log.Printf("读取消息失败: %v", err)
		return p, err
	}
	return p, nil
}

func (P *PVEService) SetConn(conn *websocket.Conn) {
	P.ws = conn
	return
}
