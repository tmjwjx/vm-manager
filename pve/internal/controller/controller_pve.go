package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"pve/internal/logic"
	"pve/pkg/globals"
	"pve/pkg/utils"
)

func CreateVM(conn *websocket.Conn, data []byte) {
	
	fmt.Println("create vm")
	vmid, err := logic.CreateVM(data)
	if err != nil {
		fmt.Println(err)
	}
	
	// 通知 vm 服务
	s := struct {
		VMID string `json:"vm_id"`
	}{
		VMID: vmid,
	}
	b, err := json.Marshal(s)
	if err != nil {
		return
	}
	utils.Write(conn, globals.CreateType, b)
}
