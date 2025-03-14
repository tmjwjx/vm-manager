package service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"vm/internal/infrastructure/globals"
	pvm "vm/internal/interfaces/grpc/proto/vm"
)

// IPVEClient PVE客户端接口
// 作为防腐层隔离外部PVE服务的变化
type IPVEClient interface {
	// CreateVM 在PVE上创建虚拟机
	CreateVM() error
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
func (c *PVEClient) CreateVM(req *pvm.CreateVMReq) error {

	//c.conn.WriteMessage(websocket.TextMessage, []byte("create"))

	//// 确保 URL 包含协议部分
	//url := "http://localhost:8088/vm/create"
	//
	//// 发送 POST 请求
	//conn, err := http.Post(url, "application/json", nil)
	//if err != nil {
	//	return "", err
	//}
	//defer conn.Body.Close()
	//
	//// 正确读取响应体
	//body, err := ioutil.ReadAll(conn.Body)
	//if err != nil {
	//	return "", err
	//}
	//
	//// 打印响应内容
	//fmt.Println(string(body))
	//return string(body), nil

	// 构造请求消息
	type PVECreateReq struct {
		Email string `json:"email"`
	}
	mes := PVECreateReq{
		Email: req.UserEmail,
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
