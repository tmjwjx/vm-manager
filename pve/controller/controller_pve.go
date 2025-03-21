package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"pve/logic"
	"pve/models"
	"pve/pkg/globals"
	"pve/pkg/utils"
)

func CreateVM(conn *websocket.Conn, data []byte) {
	// 创建虚拟机
	vmId, err := logic.CreateVM()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 解压数据
	var e models.CreateReq
	err = json.Unmarshal(data, &e)

	// 获取虚拟机的IP地址
	ipAddr, err := logic.GetIPAddr(vmId)
	if err != nil {
		fmt.Println(err)
	}

	// 通知 vm 服务
	s := &models.CreateResp{
		VMID:   vmId,
		Email:  e.Email,
		IPAddr: ipAddr,
	}
	b, err := json.Marshal(s)
	if err != nil {
		return
	}
	utils.Write(conn, globals.CreateType, b)

	// 回调siwu
	logic.CreateVMCallback(s)
}

func DestroyVM(data []byte) {

	fmt.Println("Destroy vm")
	logic.DestroyVM(data)

}
