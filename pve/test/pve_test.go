package test

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"pve/logic"
	"testing"
)

func TestGetVMInfo(t *testing.T) {
	vmInfo, err := logic.GetVMInfo("111")
	if err != nil {
		fmt.Println(err)
	}
	spew.Dump("VM Info:", vmInfo)
}

func TestSSHSetStaticIP(t *testing.T) {
	err := logic.SSHSetStaticIP("111", "192.168.10.59/24")
	if err != nil {
		fmt.Println(err)
	}
}

func TestGetIPAddr(t *testing.T) {
	ip, err := logic.GetIPAddr("112")
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
