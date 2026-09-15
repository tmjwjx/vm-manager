package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"pve/logic"
	"pve/models"
	"pve/pkg/globals"
	"pve/pkg/utils"
	"time"
)

func GetVMInfo(conn *websocket.Conn, data []byte) {
	// 解压数据
	var e models.GetVMInfoReq
	err := json.Unmarshal(data, &e)
	if err != nil {
		return
	}

	// 获取虚拟机信息
	vmInfo, err := logic.GetVMInfo(e.VMID)
	if err != nil {
		fmt.Println("获取虚拟机信息错误 err =", err)
		return
	}

	// 回调通知思悟

	fmt.Println("获取虚拟机信息成功 vmInfo =", vmInfo)

}

func RebootVM(conn *websocket.Conn, data []byte) {
	// 解压数据
	var e models.RebootReq
	err := json.Unmarshal(data, &e)
	if err != nil {
		return
	}

	// 重启虚拟机
	err = logic.RebootVM(e.VMID)
	if err != nil {
		fmt.Println("重启虚拟机错误 err =", err)
		return
	}
}

func StartVM(conn *websocket.Conn, data []byte) {
	// 解压数据
	var e models.StartReq
	err := json.Unmarshal(data, &e)
	if err != nil {
		return
	}

	// 启动虚拟机
	err = logic.StartVM(e.VMID)
	if err != nil {
		fmt.Println("启动虚拟机错误 err =", err)
		return
	}
}

func CreateVM(conn *websocket.Conn, data []byte) {
	// 创建虚拟机
	vmId, err := logic.CreateVM()
	if err != nil {
		fmt.Println("创建虚拟机错误 err =", err)
		return
	}

	// 解压数据
	var e models.CreateReq
	err = json.Unmarshal(data, &e)

	// 启动虚拟机
	err = logic.StartVM(vmId)
	if err != nil {
		fmt.Println("启动虚拟机错误 err =", err)
	}

	time.Sleep(2 * time.Second)

	// 获取虚拟机的IP地址
	ipAddr, err := logic.GetIPAddr(vmId)
	if err != nil {
		fmt.Println("获取虚拟机IP地址错误 err =", err)
	}

	// 设置虚拟机静态IP
	err = logic.SSHSetStaticIP(vmId, ipAddr)
	if err != nil {
		fmt.Println("设置虚拟机静态IP错误 err =", err)
	}

	// 重启虚拟机
	err = logic.RebootVM(vmId)
	if err != nil {
		fmt.Println("重启虚拟机错误 err =", err)
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
