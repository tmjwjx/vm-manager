package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"pve/models"
)

func CreateVMCallback(req *models.CreateResp) {
	// 要调用的接口地址
	url := "http://localhost:8081/vm/create_callback"
	
	// 将结构体转为 json
	jsonData, err := json.Marshal(req)
	if err != nil {
		fmt.Println("json.Marshal failed, err:", err)
		return
	}
	
	// 创建一个请求
	conn, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	
	// 发起 POST 请求
	client := &http.Client{}
	resp, err := client.Do(conn)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()
	// 没有响应，也不需要读取了
}
