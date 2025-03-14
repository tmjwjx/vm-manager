package globals

import (
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	DB   *gorm.DB
	Conn *websocket.Conn
	Log  *zap.SugaredLogger
	// RDB redis链接
	RDB *redis.Client
)
