package service

import (
	"encoding/json"
)

type Data struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

const (
	CreateType = "create"
	DeleteType = "delete"
)

func ServerPVE(mes []byte) {
	var data Data
	_ = json.Unmarshal(mes, &data)
	switch data.Type {
	case CreateType:
		// 创建虚拟机
	case DeleteType:
		// 删除虚拟机
	}
}
