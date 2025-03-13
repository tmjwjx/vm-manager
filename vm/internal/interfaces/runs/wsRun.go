package runs

import (
	"net/http"
	ws "vm/internal/interfaces/websocket"
)

func WebSocketRun() {
	http.HandleFunc("/", ws.WS)
	_ = http.ListenAndServe(":8088", nil)
}
