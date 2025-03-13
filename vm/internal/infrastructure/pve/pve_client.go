package pve

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// IPVEClient PVE客户端接口
// 作为防腐层隔离外部PVE服务的变化
type IPVEClient interface {
	// CreateVM 在PVE上创建虚拟机
	CreateVM() (string, error)

	// DestroyVM 在PVE上销毁虚拟机
	DestroyVM(vmId string) error

	// PowerOnVM 开启虚拟机
	PowerOnVM(vmId string) error

	// PowerOffVM 关闭虚拟机
	PowerOffVM(vmId string) error

	// GetVMStatus 获取虚拟机状态
	GetVMStatus(vmId string) (string, error)

	// ListVMs 获取虚拟机列表
	ListVMs() ([]map[string]interface{}, error)
}

// PVEClient PVE客户端实现
type PVEClient struct {
	ApiUrl   string
	ApiToken string
	NodeName string
}

// NewPVEClient 创建PVE客户端
func NewPVEClient(apiUrl, apiToken, nodeName string) *PVEClient {
	return &PVEClient{
		ApiUrl:   apiUrl,
		ApiToken: apiToken,
		NodeName: nodeName,
	}
}

// CreateVM 实现创建虚拟机
func (c *PVEClient) CreateVM() (string, error) {
	// 确保 URL 包含协议部分
	url := "http://localhost:8088/vm/create"

	// 发送 POST 请求
	conn, err := http.Post(url, "application/json", nil)
	if err != nil {
		return "", err
	}
	defer conn.Body.Close()

	// 正确读取响应体
	body, err := ioutil.ReadAll(conn.Body)
	if err != nil {
		return "", err
	}

	// 打印响应内容
	fmt.Println(string(body))

	return string(body), nil
}

// DestroyVM 实现销毁虚拟机
func (c *PVEClient) DestroyVM(vmId string) error {
	// TODO: 实现与PVE API的通信逻辑
	return nil
}

// PowerOnVM 实现开启虚拟机
func (c *PVEClient) PowerOnVM(vmId string) error {
	// TODO: 实现与PVE API的通信逻辑
	return nil
}

// PowerOffVM 实现关闭虚拟机
func (c *PVEClient) PowerOffVM(vmId string) error {
	// TODO: 实现与PVE API的通信逻辑
	return nil
}

// GetVMStatus 实现获取虚拟机状态
func (c *PVEClient) GetVMStatus(vmId string) (string, error) {
	// TODO: 实现与PVE API的通信逻辑
	return "running", nil
}

// ListVMs 实现获取虚拟机列表
func (c *PVEClient) ListVMs() ([]map[string]interface{}, error) {
	// TODO: 实现与PVE API的通信逻辑
	return []map[string]interface{}{}, nil
}
