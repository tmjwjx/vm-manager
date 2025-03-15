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
	vmid, err := logic.CreateVM()
	if err != nil {
		fmt.Println(err)
		return
	}
	
	// 解压数据
	e := struct {
		Email string `json:"email"`
	}{}
	err = json.Unmarshal(data, &e)
	
	// 通知 vm 服务
	s := struct {
		VMID  string `json:"vm_id"`
		Email string `json:"email"`
	}{
		VMID:  vmid,
		Email: e.Email,
	}
	b, err := json.Marshal(s)
	if err != nil {
		return
	}
	utils.Write(conn, globals.CreateType, b)
}

func DestroyVM(data []byte) {
	
	fmt.Println("Destroy vm")
	logic.DestroyVM(data)
	
}
