package globals

import (
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	Conn *websocket.Conn
)
