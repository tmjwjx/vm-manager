package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"vm/internal/infrastructure/globals"
)

type Data struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

const (
	CreateType = 1
	DeleteType = 2
)

func WS(w http.ResponseWriter, r *http.Request) {
	var up = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := up.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	globals.Conn = conn
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("收到消息: %s", message)
		// 处理信息
		
	}
}
