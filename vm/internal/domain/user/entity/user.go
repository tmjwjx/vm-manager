package entity

import "github.com/jinzhu/gorm"

// User
// @Description: 用户实体
// @Author tianjiajie 2025-03-12 11:30:00
type User struct {
	gorm.Model
	Email string `gorm:"type:varchar(128);uniqueIndex"`
}