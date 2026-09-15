package models

type CreateResp struct {
	VMID   string `json:"vm_id"`
	Email  string `json:"email"`
	IPAddr string `json:"ip_addr"`
}

type GetVMInfoResp struct {
	VMID       string            `json:"vm_id"`       // 虚拟机唯一标识
	IPAddr     string            `json:"ip_addr"`     // 主IP地址
	Status     string            `json:"status"`      // 运行状态(running/stopped)
	Cores      int               `json:"cores"`       // CPU核心数
	Memory     int               `json:"memory"`      // 内存容量(MB)
	Disk       map[string]string `json:"disk"`        // 磁盘信息(key=磁盘接口名, value=磁盘配置)
	Networks   []NetworkInfo     `json:"networks"`    // 网络接口详情
	OSType     string            `json:"os_type"`     // 操作系统类型(l26=Linux, w2k=Windows等)
	QemuVer    string            `json:"qemu_ver"`    // QEMU版本
	BootOrder  []string          `json:"boot_order"`  // 启动顺序(如[net0,scsi0])
	Agent      int               `json:"agent"`       // QEMU代理状态(0=禁用,1=启用)
	Template   string            `json:"template"`    // 基础模板ID（如果基于模板创建）
	Notes      string            `json:"notes"`       // 备注信息
	CreateTime string            `json:"create_time"` // 创建时间(时间戳或格式化时间)
}

type NetworkInfo struct {
	Interface string   `json:"interface"` // 接口名(如net0)
	MAC       string   `json:"mac"`       // MAC地址
	Bridge    string   `json:"bridge"`    // 桥接网络名称(如vmbr0)
	IPs       []string `json:"ips"`       // 分配的IP地址列表(IPv4+IPv6)
}
