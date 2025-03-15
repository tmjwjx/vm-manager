package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"pve/pkg/globals"
)

func Write(conn *websocket.Conn, t globals.DataType, b []byte) {
	data1 := globals.Data{
		Type: t,
		Data: b,
	}
	data2, _ := json.Marshal(data1)
	err := conn.WriteJSON(data2)
	if err != nil {
		fmt.Println(err)
	}
}
