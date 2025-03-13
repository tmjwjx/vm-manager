package service

import (
	"github.com/gorilla/websocket"
)

// IPVEClient PVE客户端接口
// 作为防腐层隔离外部PVE服务的变化
type IPVEClient interface {
	// CreateVM 在PVE上创建虚拟机
	CreateVM() (string, error)
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
func (c *PVEClient) CreateVM() (string, error) {
	
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
	return "", nil // 暂时返回空字符串
}
