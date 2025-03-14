package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pvm "vm/internal/interfaces/grpc/proto/vm"
)

func main() {
	// 连接grpc服务 默认禁用安全连接 没有加密和认证
	conn, err := grpc.Dial("localhost:8888", grpc.WithInsecure())
	if err != nil {
		log.Printf("连接失败: %v", err)
	}
	defer conn.Close()

	// 建立连接
	client := pvm.NewVMManagerClient(conn)

	// 调用服务
	resp, err := client.CreateVM(context.Background(), &pvm.CreateVMReq{
		UserEmail: "2937693919@qq.com",
		Token:     "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MzIsImVtYWlsIjoiMjkzNzY5MzkxOUBxcS5jb20iLCJpc3MiOiJzaXd1LXdlYi1zZXJ2aWNlIiwiZXhwIjoxNzQyMjU3OTU2LCJpYXQiOjE3NDE5NTU1NTZ9.tfX6T_Wa-dhYmK4ZvP5PGr8H8aUH4fMsTQD04VSJYVs",
	})
	if err != nil {
		log.Printf("调用失败: %v", err)
	} else {
		log.Printf("调用成功: %v", resp)
	}
}
