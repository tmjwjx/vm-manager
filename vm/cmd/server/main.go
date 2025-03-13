package main

import (
	"log"
	_ "vm/internal/infrastructure/inits"
	"vm/internal/interfaces/runs"
)

func main() {
	log.Printf("服务启动成功")

	runs.Run()
}
