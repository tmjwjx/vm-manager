package test

import (
	"fmt"
	"pve/logic"
	"testing"
)

func TestGetIPAddr(t *testing.T) {
	ip, err := logic.GetIPAddr("111")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("ip:", ip)
}

func TestStartVM(t *testing.T) {
	err := logic.StartVM("111")
	if err != nil {
		fmt.Println(err)
	}
}

func TestSetStaticIP(t *testing.T) {
	err := logic.SetStaticIP("111", "ip=192.168.10.9/24,gw=192.168.10.1")
	if err != nil {
		fmt.Println(err)
	}
}
