package websocket

import (
	"net/http"
)

func WebSocketRun() {
	http.HandleFunc("/", WS)
	_ = http.ListenAndServe(":8088", nil)
}
