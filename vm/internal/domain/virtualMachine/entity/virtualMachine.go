package entity

import (
	"github.com/jinzhu/gorm"
	"time"
)

type VmStatus string

const (
	StatusRunning VmStatus = "running"
	StatusStopped VmStatus = "stopped"
	StatusPaused  VmStatus = "paused"
	StatusError   VmStatus = "error"
)

// VirtualMachine
// @Description: 虚拟机实体
// @Author tianjiajie 2025-03-12 17:24:32
type VirtualMachine struct {
	gorm.Model
	Name           string     `json:"name" gorm:"unique;not null;size:64"`                                             // 虚拟机名称，唯一且不能为空，最大长度为64
	HostID         uint       `json:"host_id" gorm:"not null"`                                                         // 宿主机ID
	UserID         int        `json:"user_id" gorm:"not null"`                                                         // 关联用户逻辑外键，不能为空
	OSType         string     `json:"os_type" gorm:"not null;size:32"`                                                 // 操作系统类型，不能为空，最大长度为32
	CpuCores       int        `json:"cpu_cores" gorm:"not null"`                                                       // CPU核心数，不能为空
	MemoryMB       int        `json:"memory_mb" gorm:"not null"`                                                       // 内存大小（MB），不能为空
	DiskGB         int        `json:"disk_gb" gorm:"not null"`                                                         // 磁盘大小（GB），不能为空
	Status         VmStatus   `json:"status" gorm:"type:enum('running','stopped','paused','error');default:'stopped'"` // 虚拟机状态，枚举类型，默认值为stopped
	IPAddress      string     `json:"ip_address" gorm:"size:15"`                                                       // IP地址，最大长度为15
	ExpirationTime *time.Time `json:"expiration_time"`                                                                 // 过期时间，允许为NULL
}
