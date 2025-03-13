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
	VMID           uint       `json:"vm_id" gorm:"unique;not null;size:64"` // 虚拟机ID，唯一且不能为空，最大长度为64
	Email          string     `gorm:"type:varchar(128);"`                   // 用户邮箱
	ExpirationTime *time.Time `json:"expiration_time"`                      // 过期时间，允许为NULL
}
