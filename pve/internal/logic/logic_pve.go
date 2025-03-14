package logic

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func CreateVM(data []byte) {
	url := "http://192.168.10.2:8006/api2/json/nodes/lezhi/qemu/112/clone"
	method := "POST"

	payload := strings.NewReader("newid=123&name=tianjiajie")

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
