package services

import "github.com/gorilla/websocket"

type IPVEService interface {
	CreateVM(email string) error
	DestroyVM(vmId string) error
	StartVM(vmId string) error
	StopVM(vmId string) error
	RenewVM(vmId string) error
	GetVMInfo(vmId string) error
	SetConn(conn *websocket.Conn)
	ReceiveMessage() ([]byte, error)
}
