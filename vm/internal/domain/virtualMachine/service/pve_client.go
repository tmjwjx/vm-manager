package service

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"vm/internal/infrastructure/globals"
	pvm "vm/internal/interfaces/grpc/proto/vm"
)

// IPVEClient PVE客户端接口
// 作为防腐层隔离外部PVE服务的变化
type IPVEClient interface {
	// CreateVM 在PVE上创建虚拟机
	CreateVM(ctx context.Context, req *pvm.CreateVMReq) error
}

// PVEClient PVE客户端实现
type PVEClient struct {
	conn *websocket.Conn
}

// NewPVEClient 创建PVE客户端
func NewPVEClient(conn *websocket.Conn) *PVEClient {
	return &PVEClient{conn: conn}
}

// CreateVM 实现创建虚拟机
func (c *PVEClient) CreateVM(ctx context.Context, req *pvm.CreateVMReq) error {

	// 构造请求消息
	type PVECreateReq struct {
		Email string `json:"email"`
	}

	email := ctx.Value("email")

	mes := PVECreateReq{
		Email: email.(string),
	}
	b, err := json.Marshal(mes)

	// 发送创建虚拟机请求
	data := globals.Data{
		Type: globals.CreateType,
		Data: b,
	}
	b, err = json.Marshal(data)
	if err != nil {
		return err
	}
	_ = c.conn.WriteMessage(websocket.TextMessage, b)

	return nil // 暂时返回空字符串
}

// DestroyVM 销毁虚拟机
func (c *PVEClient) DestroyVM(req *pvm.DestroyVMReq) error {

	return nil
}
