package controller

import (
	"fmt"
	"pve/internal/logic"
)

func CreateVM(data []byte) {

	fmt.Println("create vm")
	logic.CreateVM(data)

}
