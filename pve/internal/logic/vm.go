package logic

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func CreateVM() {
	url := "http://192.168.10.2:8006/api2/json/nodes/lezhi/qemu/112/clone"
	method := "POST"

	payload := strings.NewReader("newid=111&name=tianjiajie")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "PVEAPIToken=root@pam!vmManager=613102db-dc90-41f4-960c-4754d05096b8")
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "192.168.10.2:8006")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", "http://192.168.10.2:8006/api2/json/nodes/lezhi/qemu/112/clone")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
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
