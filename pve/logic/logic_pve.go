package logic

import (
	"crypto/tls"
	"fmt"
	"github.com/goccy/go-json"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"pve/models"
	"pve/pkg/globals"
	"strings"
	"time"
)

// GetVMInfo
// @Description: 获取虚拟机信息
// @param        vmid string
// @return       *models.GetVMInfoResp
// @return       error
// @Author tianjiajie 2025-03-27 21:59:37
func GetVMInfo(vmid string) (*models.GetVMInfoResp, error) {
	// ----------------------------
	// 1. 获取虚拟机配置信息
	// ----------------------------
	url := fmt.Sprintf("%s/api2/json/nodes/lezhi/qemu/%s/config",
		globals.PVEURL, vmid)

	client := createHTTPClient()
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("PVEAPIToken %s", globals.TOKEN))
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("API返回状态码: %d", res.StatusCode)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("读取响应失败: %v", err)
		return nil, err
	}

	// 解析配置数据
	var config struct {
		Data map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(b, &config); err != nil {
		return nil, fmt.Errorf("配置解析失败: %v", err)
	}

	// ----------------------------
	// 2. 获取虚拟机当前状态
	// ----------------------------
	statusURL := fmt.Sprintf("%s/api2/json/nodes/lezhi/qemu/%s/status/current",
		globals.PVEURL, vmid)

	req, _ = http.NewRequest("GET", statusURL, nil)
	req.Header.Add("Authorization", fmt.Sprintf("PVEAPIToken %s", globals.TOKEN))
	res, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("API返回状态码: %d", res.StatusCode)
	}
	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("读取响应失败: %v", err)
		return nil, err
	}

	// 解析状态数据
	var status struct {
		Data struct {
			Status string  `json:"status"`
			Uptime int     `json:"uptime"` // 运行时长(秒)
			CPUs   float64 `json:"cpus"`   // CPU核心数（可能带小数）
			Mem    int     `json:"mem"`    // 已用内存(bytes)
			MaxMem int     `json:"maxmem"` // 最大内存(bytes)
		} `json:"data"`
	}
	if err := json.Unmarshal(b, &status); err != nil {
		return nil, fmt.Errorf("状态解析失败: %v", err)
	}

	// ----------------------------
	// 3. 数据整合与转换
	// ----------------------------
	resp := &models.GetVMInfoResp{
		VMID:   vmid,
		Status: status.Data.Status,
		Cores:  int(status.Data.CPUs),
		Memory: status.Data.MaxMem / 1024 / 1024, // 转换为MB
	}

	// 填充配置字段
	if val, ok := config.Data["ostype"]; ok {
		resp.OSType = val.(string)
	}

	// 解析磁盘信息
	resp.Disk = make(map[string]string)
	for key, value := range config.Data {
		if strings.HasPrefix(key, "scsi") || strings.HasPrefix(key, "sata") {
			resp.Disk[key] = value.(string)
		}
	}

	// 解析网络信息（示例解析net0）
	if netConfig, ok := config.Data["net0"]; ok {
		parts := strings.Split(netConfig.(string), ",")
		netInfo := models.NetworkInfo{Interface: "net0"}
		for _, part := range parts {
			if strings.HasPrefix(part, "virtio=") {
				netInfo.MAC = strings.Split(part, "=")[1]
			} else if strings.HasPrefix(part, "bridge=") {
				netInfo.Bridge = strings.Split(part, "=")[1]
			}
		}
		resp.Networks = append(resp.Networks, netInfo)
	}

	return resp, nil
}

// createHTTPClient 创建HTTP客户端，跳过证书验证
func createHTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

// StartVM 启动虚拟机，并等待其状态为running且获取IP地址
// @Description: 启动虚拟机，等待其运行并获取IP地址
// @param        vmid string 虚拟机ID
// @return       error 错误信息
func StartVM(vmid string) error {
	client := createHTTPClient()
	token := fmt.Sprintf("PVEAPIToken %s", globals.TOKEN)

	// 尝试发送启动请求，最多重试5次
	maxRetries := 5
	retryInterval := 2 * time.Second
	for attempt := 0; attempt < maxRetries; attempt++ {
		url := fmt.Sprintf("%s/api2/json/nodes/lezhi/qemu/%s/status/start", globals.PVEURL, vmid)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			log.Printf("创建请求失败: %v", err)
			return err
		}
		req.Header.Add("Authorization", token)
		res, err := client.Do(req)
		if err != nil {
			log.Printf("pve接口连接失败: %v", err)
			time.Sleep(retryInterval)
			continue
		}
		res.Body.Close()
		if res.StatusCode == http.StatusOK {
			break
		}
		log.Printf("启动虚拟机请求失败，状态码: %d，尝试次数: %d/%d", res.StatusCode, attempt+1, maxRetries)
		if attempt == maxRetries-1 {
			return fmt.Errorf("虚拟机 %s 启动失败，达到最大重试次数", vmid)
		}
		time.Sleep(retryInterval)
	}

	// 等待虚拟机状态为running，最多60秒
	statusMaxWait := 30 * time.Second
	statusInterval := 5 * time.Second
	startTime := time.Now()
	for {
		if time.Since(startTime) > statusMaxWait {
			return fmt.Errorf("等待虚拟机 %s 启动超时", vmid)
		}
		status, err := GetVMStatus(vmid)
		if err != nil {
			log.Printf("获取虚拟机状态失败: %v", err)
			time.Sleep(statusInterval)
			continue
		}
		if status == "running" {
			break
		}
		log.Printf("等待虚拟机启动... 当前状态: %s", status)
		time.Sleep(statusInterval)
	}

	// 等待获取IP地址，最多120秒
	ipMaxWait := 30 * time.Second
	ipInterval := 5 * time.Second
	ipStartTime := time.Now()
	for {
		if time.Since(ipStartTime) > ipMaxWait {
			return fmt.Errorf("等待虚拟机 %s 获取IP地址超时", vmid)
		}
		ip, err := GetIPAddr(vmid)
		if err == nil && ip != "" {
			log.Printf("虚拟机 %s 启动成功，IP地址: %s", vmid, ip)
			return nil
		}
		log.Printf("等待虚拟机获取IP地址...")
		time.Sleep(ipInterval)
	}
}

// RebootVM 重启虚拟机，并等待其状态为running且获取IP地址
// @Description: 重启虚拟机，等待其运行并获取IP地址
// @param        vmid string 虚拟机ID
// @return       error 错误信息
func RebootVM(vmid string) error {
	client := createHTTPClient()
	token := fmt.Sprintf("PVEAPIToken %s", globals.TOKEN)

	// 尝试发送重启请求，最多重试5次
	maxRetries := 5
	retryInterval := 2 * time.Second
	for attempt := 0; attempt < maxRetries; attempt++ {
		url := fmt.Sprintf("%s/api2/json/nodes/lezhi/qemu/%s/status/reboot", globals.PVEURL, vmid)
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			log.Printf("创建请求失败: %v", err)
			return err
		}
		req.Header.Add("Authorization", token)
		res, err := client.Do(req)
		if err != nil {
			log.Printf("pve接口连接失败: %v", err)
			time.Sleep(retryInterval)
			continue
		}
		res.Body.Close()
		if res.StatusCode == http.StatusOK {
			break
		}
		log.Printf("重启虚拟机请求失败，状态码: %d，尝试次数: %d/%d", res.StatusCode, attempt+1, maxRetries)
		if attempt == maxRetries-1 {
			return fmt.Errorf("虚拟机 %s 重启失败，达到最大重试次数", vmid)
		}
		time.Sleep(retryInterval)
	}

	// 等待虚拟机状态为running，最多60秒
	statusMaxWait := 30 * time.Second
	statusInterval := 5 * time.Second
	startTime := time.Now()
	for {
		if time.Since(startTime) > statusMaxWait {
			return fmt.Errorf("等待虚拟机 %s 重启超时", vmid)
		}
		status, err := GetVMStatus(vmid)
		if err != nil {
			log.Printf("获取虚拟机状态失败: %v", err)
			time.Sleep(statusInterval)
			continue
		}
		if status == "running" {
			return nil
		}
		log.Printf("等待虚拟机重启... 当前状态: %s", status)
		time.Sleep(statusInterval)
	}
}

// GetVMStatus 获取虚拟机状态
// @Description: 获取虚拟机状态
// @param        vmid string 虚拟机ID
// @return       string 状态
// @return       error 错误信息
func GetVMStatus(vmid string) (string, error) {
	url := fmt.Sprintf("%s/api2/json/nodes/lezhi/qemu/%s/status/current", globals.PVEURL, vmid)
	token := fmt.Sprintf("PVEAPIToken %s", globals.TOKEN)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", token)

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := &http.Client{Transport: tr}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var statusResponse struct {
		Data struct {
			Status string `json:"status"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &statusResponse); err != nil {
		return "", err
	}

	return statusResponse.Data.Status, nil
}

// SSHSetStaticIP 通过 SSH 连接 CentOS 7 虚拟机设置静态 IP
// 参数说明：
// - vmid: 虚拟机 ID（仅用于日志标识，实际连接需通过 sshHost）
// - ip: 静态 IP 地址（格式：192.168.10.59/24, 网关和 DNS 需预设或从参数扩展）
// 注意：需根据实际环境调整 SSH 连接参数（主机、端口、用户名、密码/密钥）
func SSHSetStaticIP(vmid string, ip string) (err error) {
	// ----------------------------
	// 1. SSH 连接配置（根据实际环境修改）
	// ----------------------------
	host, err := GetIPAddr(vmid)
	if err != nil {
		return fmt.Errorf("获取虚拟机 IP 失败: %v", err)
	}
	sshHost := host         // 虚拟机当前 IP（需可达）
	sshPort := "22"         // SSH 端口
	sshUser := "root"       // SSH 用户名（需有 sudo 权限）
	sshPassword := "123456" // SSH 密码或私钥

	// 解析 IP 和子网掩码
	ip += "/24" // 默认处理 /24
	ipParts := strings.Split(ip, "/")
	if len(ipParts) != 2 {
		return fmt.Errorf("IP 格式错误，应为 IP/CIDR（如 192.168.10.59/24）")
	}
	ipAddr := ipParts[0]
	cidr := ipParts[1]
	gateway := "192.168.10.1"       // 预设网关（可根据需要参数化）
	dnsServers := "8.8.8.8 8.8.4.4" // 预设 DNS

	// ----------------------------
	// 2. 创建 SSH 客户端
	// ----------------------------
	config := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 生产环境应验证主机密钥
		Timeout:         10 * time.Second,
	}

	client, err := ssh.Dial("tcp", net.JoinHostPort(sshHost, sshPort), config)
	if err != nil {
		return fmt.Errorf("SSH 连接失败: %v", err)
	}
	defer client.Close()

	// ----------------------------
	// 3. 修改网络配置文件
	// ----------------------------
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("创建 SSH 会话失败: %v", err)
	}
	defer session.Close()

	// 生成 ifcfg-eth0 内容
	configContent := fmt.Sprintf(`
DEVICE=eth0
BOOTPROTO=static
ONBOOT=yes
IPADDR=%s
NETMASK=%s
GATEWAY=%s
DNS1=%s
`, ipAddr, cidrToNetmask(cidr), gateway, strings.Split(dnsServers, " ")[0])

	// 通过 echo 写入配置（需要 root 权限）
	commands := []string{
		"cp /etc/sysconfig/network-scripts/ifcfg-eth0 /etc/sysconfig/network-scripts/ifcfg-eth0.bak", // 备份
		fmt.Sprintf("echo '%s' > /etc/sysconfig/network-scripts/ifcfg-eth0", strings.TrimSpace(configContent)),
		"systemctl restart network", // 重启网络服务
	}

	for _, cmd := range commands {
		session, err := client.NewSession()
		if err != nil {
			return fmt.Errorf("执行命令失败: %v", err)
		}
		defer session.Close()

		if err := session.Run(cmd); err != nil {
			return fmt.Errorf("命令执行错误: %s\n错误详情: %v", cmd, err)
		}
	}

	return nil
}

// ----------------------------
// 工具函数：CIDR 转子网掩码
// ----------------------------
func cidrToNetmask(cidr string) string {
	switch cidr {
	case "24":
		return "255.255.255.0"
	case "16":
		return "255.255.0.0"
	case "8":
		return "255.0.0.0"
	default:
		return "255.255.255.0" // 默认处理 /24
	}
}

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
