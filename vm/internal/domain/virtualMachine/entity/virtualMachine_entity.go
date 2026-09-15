package entity

import (
	"github.com/jinzhu/gorm"
	"time"
)

type VirtualMachine struct {
	gorm.Model
	VMID           string     `json:"vm_id" gorm:"unique;not null;size:64"`   // 虚拟机ID，唯一且不能为空，最大长度为64
	Email          string     `json:"email" gorm:"unique;type:varchar(128);"` // 用户邮箱
	IPAddr         string     `json:"ip_addr" gorm:"size:64"`                 // IP地址，最大长度为64
	ExpirationTime *time.Time `json:"expiration_time"`                        // 过期时间，允许为NULL
}
