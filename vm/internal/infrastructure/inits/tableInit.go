package inits

import (
	"gorm.io/gorm"
	"log"
	"vm/internal/domain/virtualMachine/entity"
)

func TableInit(db *gorm.DB) {
	err := db.AutoMigrate(
		&entity.VirtualMachine{},
	)
	if err != nil {
		log.Printf("初始化table失败: %v", err)
	}
}
