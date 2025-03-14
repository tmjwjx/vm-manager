package logic

import (
	"crypto/tls"
	"fmt"
	"github.com/goccy/go-json"
	"io/ioutil"
	"net/http"
	"pve/pkg/globals"
	"strings"
)

func CreateVM(data []byte) {
	vmid := GetMinID()
	
	url := globals.PVEURL + "/api2/json/nodes/lezhi/qemu/112/clone"
	method := "POST"
	
	// 传递的参数
	t := "newid=" + vmid + "&name=lezhi-" + vmid
	payload := strings.NewReader(t)
	
	// 跳过证书验证（仅限测试环境）
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 跳过证书验证
		},
	}
	
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	req.Header.Add("Authorization", "PVEAPIToken=root@pam!vmManager=613102db-dc90-41f4-960c-4754d05096b8")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("pve接口连接失败: %v\n", err)
		return
	}
	defer res.Body.Close()
	
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func GetMinID() string {
	type NextIDResponse struct {
		Data string `json:"data"`
	}
	url := globals.PVEURL + "/api2/json/cluster/nextid"
	method := "GET"
	
	// 跳过证书验证（仅限测试环境）
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 跳过证书验证
		},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	
	req.Header.Add("Authorization", "PVEAPIToken=root@pam!vmManager=613102db-dc90-41f4-960c-4754d05096b8")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("pve接口连接失败: %v\n", err)
		return ""
	}
	defer res.Body.Close()
	
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	
	var nextIDResponse NextIDResponse
	err = json.Unmarshal(body, &nextIDResponse)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return nextIDResponse.Data
}
