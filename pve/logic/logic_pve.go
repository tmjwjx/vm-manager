package logic

import (
	"crypto/tls"
	"fmt"
	"github.com/goccy/go-json"
	"io/ioutil"
	"net/http"
	"pve/pkg/globals"
	"strings"
	"time"
)

func CreateVM() (vmid string, err error) {
	// 获取可用的最小id
	vmid = GetMinID()

	url := fmt.Sprintf("%s/api2/json/nodes/lezhi/qemu/%s/clone", globals.PVEURL, globals.TEMPLATE)
	method := "POST"

	// 传递的参数
	t := "newid=" + vmid + "&name=lezhi-" + vmid
	payload := strings.NewReader(t)

	// 跳过证书验证(小组服务器好像没有证书)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 跳过证书验证
		},
	}

	// 创建客户端
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 添加请求头
	req.Header.Add("Authorization", "PVEAPIToken=root@pam!vmManager=613102db-dc90-41f4-960c-4754d05096b8")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("pve接口连接失败: %v\n", err)
		return
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}

// StartVM 启动虚拟机，直到启动成功为止
// @Description: 启动虚拟机，直到启动成功为止
// @param        vmid string
// @return       err
// @Author tianjiajie 2025-03-26 15:32:36
func StartVM(vmid string) (err error) {
	url := fmt.Sprintf("%s/api2/json/nodes/lezhi/qemu/%s/status/start", globals.PVEURL, vmid)
	method := "POST"
	token := fmt.Sprintf("PVEAPIToken %s", globals.TOKEN)

	// 跳过证书验证(小组服务器没有证书)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 跳过证书验证
		},
	}

	client := &http.Client{Transport: tr}

	// 最大重试次数
	maxRetries := 5
	// 重试间隔
	retryInterval := 2 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			fmt.Println(err)
			return err
		}

		req.Header.Add("Authorization", token)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		res, err := client.Do(req)
		if err != nil {
			fmt.Printf("pve接口连接失败: %v\n", err)
			return err
		}
		defer res.Body.Close()

		_, err = ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// 检查启动请求是否成功
		if res.StatusCode != http.StatusOK {
			fmt.Printf("启动虚拟机请求失败，状态码: %d，尝试次数: %d/%d\n", res.StatusCode, attempt+1, maxRetries)
			time.Sleep(retryInterval)
			continue
		}

		// 轮询虚拟机状态，直到虚拟机启动成功或超时
		statusURL := fmt.Sprintf("%s/api2/json/nodes/lezhi/qemu/%s/status/current", globals.PVEURL, vmid)

		statusReq, err := http.NewRequest("GET", statusURL, nil)
		if err != nil {
			fmt.Println(err)
			return err
		}

		statusReq.Header.Add("Authorization", token)
		statusRes, err := client.Do(statusReq)
		if err != nil {
			fmt.Printf("获取虚拟机状态失败: %v\n", err)
			return err
		}
		defer statusRes.Body.Close()

		body, err := ioutil.ReadAll(statusRes.Body)
		if err != nil {
			fmt.Println(err)
			return err
		}

		var statusResponse struct {
			Data struct {
				Status string `json:"status"`
			} `json:"data"`
		}

		err = json.Unmarshal(body, &statusResponse)
		if err != nil {
			fmt.Println(err)
			return err
		}

		if statusResponse.Data.Status == "running" {
			fmt.Printf("虚拟机 %s 启动成功\n", vmid)
			return nil
		}

		time.Sleep(retryInterval)
	}

	return fmt.Errorf("虚拟机 %s 启动失败，达到最大重试次数", vmid)
}

func SSHSetStaticIP(vmid string, ip string) (err error) {

	return err
}

// SetStaticIP
// @Description: 设置虚拟机静态IP
// @param        vmid string
// @param        ip string
// @return       err
// @Author tianjiajie 2025-03-26 16:23:47
func SetStaticIP(vmid string, ip string) (err error) {
	url := fmt.Sprintf("%s/api2/json/nodes/lezhi/qemu/%s/agent/network-config-set", globals.PVEURL, vmid)
	method := "POST"
	token := fmt.Sprintf("PVEAPIToken %s", globals.TOKEN)

	// 构建网络配置参数
	networkConfig := fmt.Sprintf(`{
		"type": "manual",
		"address": "%s",
		"netmask": "255.255.255.0",
		"gateway": "192.168.10.1",
		"interface": "eth0"
	}`, ip)

	payload := strings.NewReader(networkConfig)

	// 跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{Transport: tr}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("设置静态IP失败: %v", err)
	}
	defer res.Body.Close()

	// 检查响应状态
	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("设置静态IP失败，状态码: %d，响应: %s", res.StatusCode, string(body))
	}

	// 等待网络配置生效
	time.Sleep(5 * time.Second)

	return nil
}

// GetIPAddr
// @Description: 获取虚拟机的IP地址
// @param        vmId string
// @return       ip
// @return       err
// @Author tianjiajie 2025-03-26 16:36:00
func GetIPAddr(vmId string) (ip string, err error) {

	type NetworkInfo struct {
		Data struct {
			Result []struct {
				Name        string `json:"name"`
				IpAddresses []struct {
					IpAddressType string `json:"ip-address-type"`
					IpAddress     string `json:"ip-address"`
					Prefix        int    `json:"prefix"`
				} `json:"ip-addresses"`
			} `json:"result"`
		} `json:"data"`
	}

	url := fmt.Sprintf("%s/api2/json/nodes/lezhi/qemu/%s/agent/network-get-interfaces", globals.PVEURL, vmId)
	method := "GET"
	token := fmt.Sprintf("PVEAPIToken %s", globals.TOKEN)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", token)

	// 跳过证书验证(小组服务器好像没有证书)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // 跳过证书验证
		},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var networkInfo NetworkInfo
	err = json.Unmarshal(body, &networkInfo)
	if err != nil {
		return "", err
	}

	// 遍历接口信息，获取第一个 IPv4 地址
	for _, interfaceInfo := range networkInfo.Data.Result {
		// 仅获取以 eth 或 ens 开头的接口
		if strings.HasPrefix(interfaceInfo.Name, "eth") ||
			strings.HasPrefix(interfaceInfo.Name, "ens") {
			for _, ipInfo := range interfaceInfo.IpAddresses {
				// 仅获取 IPv4 地址
				if ipInfo.IpAddressType == "ipv4" {
					return ipInfo.IpAddress, nil
				}
			}
		}
	}

	return "", fmt.Errorf("no IPv4 address found")
}

// GetMinID
// @Description: 获取pve中可用的最小id(100以上)
// @return       string
// @Author tianjiajie 2025-03-20 18:10:20
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

func DestroyVM(data []byte) {
	url := globals.PVEURL + "/api2/json/nodes/lezhi/qemu/" + string(data)
	method := "DELETE"

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
