package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"vm/internal/domain/pve/services"
)

var _ services.IPVEService = (*PVEService)(nil)

type PVEService struct {
	ws *websocket.Conn
}

func (P PVEService) ReceiveMessage() (p []byte, err error) {
	_, p, err = P.ws.ReadMessage()
	if err != nil {
		log.Printf("读取消息失败: %v", err)
		return p, err
	}
	return p, nil
}

func NewWebSocketClient() *PVEService {
	return &PVEService{}
}

func (P PVEService) SetConn(conn *websocket.Conn) {
	P.ws = conn
	return
}

func (P PVEService) CreateVM(email string) (err error) {
	// 构造请求消息
	req := CreateVMReq{
		Email: email,
	}
	// 序列化请求消息
	b, err := json.Marshal(req)
	if err != nil {
		log.Printf("json 序列化失败: %v", err)
		return err
	}

	// 封装请求消息
	data := Data{
		Type: CreateType,
		Data: b,
	}
	// 再次序列化
	b, err = json.Marshal(data)
	if err != nil {
		log.Printf("json 序列化失败: %v", err)
		return err
	}

	// 发送创建虚拟机请求
	err = (P.ws).WriteMessage(websocket.TextMessage, b)
	if err != nil {
		log.Printf("发送创建虚拟机请求失败: %v", err)
		return err
	}
	return
}

func (P PVEService) DestroyVM(vmId string) error {
	//TODO implement me
	panic("implement me")
}

func (P PVEService) StartVM(vmId string) error {
	//TODO implement me
	panic("implement me")
}

func (P PVEService) StopVM(vmId string) error {
	//TODO implement me
	panic("implement me")
}

func (P PVEService) RenewVM(vmId string) error {
	//TODO implement me
	panic("implement me")
}

func (P PVEService) GetVMInfo(vmId string) error {
	//TODO implement me
	panic("implement me")
}
