package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"pve/controller"
	"pve/pkg/globals"
)

type WebsocketClient struct {
	conn *websocket.Conn
}

func (W *WebsocketClient) ProcessMessage(mes []byte) {
	var data globals.Data
	err := json.Unmarshal(mes, &data)
	if err != nil {
		return
	}
	switch data.Type {
	case globals.CreateType:
		// 创建虚拟机
		controller.CreateVM(W.conn, data.Data)
	case globals.DestroyType:
	
	}
}

func main() {
	dl := websocket.Dialer{}
	conn, _, err := dl.Dial("ws://127.0.0.1:8088", nil)
	if err != nil {
		log.Println(err)
		return
	}
	wsc := WebsocketClient{conn: conn}
	for {
		_, p, err := wsc.conn.ReadMessage()
		if err != nil {
			break
		}
		wsc.ProcessMessage(p)
	}
}
