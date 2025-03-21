package services

import (
	"github.com/gorilla/websocket"
	websocket2 "vm/internal/infrastructure/websocket"
)

type IPVEService interface {
	SetConn(conn *websocket.Conn)
	SendMessage(data []byte) error
	ReceiveMessage() ([]byte, error)
}

var _ IPVEService = (*websocket2.PVEService)(nil)
