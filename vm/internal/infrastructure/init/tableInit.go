package init

import (
	"log"
	"vm/internal/domain/virtualMachine/entity"
	"vm/internal/infrastructure/globals"
)

func TableInit() {
	err := globals.DB.AutoMigrate(
		&entity.VirtualMachine{},
	)
	if err != nil {
		log.Printf("初始化table失败: %v", err)
	}
}
