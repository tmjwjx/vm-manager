package globals

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	Conn *websocket.Conn
	Log  *zap.SugaredLogger
)
