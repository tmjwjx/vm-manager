package utils

import (
	"fmt"
	"pve/pkg/globals"
	
	"github.com/gorilla/websocket"
)

func Write(conn *websocket.Conn, t string, b []byte) {
	data := globals.Data{
		Type: t,
		Data: b,
	}
	err := conn.WriteJSON(data)
	if err != nil {
		fmt.Println(err)
		return
	}
}
