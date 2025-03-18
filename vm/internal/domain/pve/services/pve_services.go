package services

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"vm/internal/domain/pve/apis"
)

type IPVEService interface {
	CreateVM(email string) error
	DestroyVM(vmId string) error
	StartVM(vmId string) error
	StopVM(vmId string) error
	RenewVM(vmId string) error
	GetVMInfo(vmId string) error
}

type PVEService struct {
	WS **websocket.Conn
}

func NewPVEService(WS **websocket.Conn) *PVEService {
	return &PVEService{WS: WS}
}

func (P PVEService) CreateVM(email string) (err error) {
	// 构造请求消息
	req := apis.CreateVMReq{
		Email: email,
	}
	// 序列化请求消息
	b, err := json.Marshal(req)
	if err != nil {
		log.Printf("json 序列化失败: %v", err)
		return err
	}

	// 封装请求消息
	data := apis.Data{
		Type: apis.CreateType,
		Data: b,
	}
	// 再次序列化
	b, err = json.Marshal(data)
	if err != nil {
		log.Printf("json 序列化失败: %v", err)
		return err
	}

	// 发送创建虚拟机请求
	err = (*P.WS).WriteMessage(websocket.TextMessage, b)
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
