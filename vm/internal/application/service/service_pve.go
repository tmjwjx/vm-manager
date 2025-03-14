package service

import (
	"encoding/json"
	"vm/internal/infrastructure/globals"
)

func ServerPVE(mes []byte) {
	var data globals.Data
	_ = json.Unmarshal(mes, &data)
	switch data.Type {
	case globals.CreateType:
		// 创建虚拟机
	case globals.DestroyType:
		// 删除虚拟机

	}
}
