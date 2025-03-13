package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"net"
	"net/http"
	"vm/internal/application/service"
	"vm/internal/domain/virtualMachine/entity"
	"vm/internal/infrastructure/pve"
	pvm "vm/internal/interfaces/grpc/proto/vm"
)

type WebsocketServer struct {
	conn *websocket.Conn
}

var wss WebsocketServer

var UP = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := UP.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	wss.conn = conn
	for {
		m, p, err := conn.ReadMessage() // m 为消息类型，p 为消息内容，err 为错误
		if err != nil {
			break
		}
		_ = wss.conn.WriteMessage(websocket.TextMessage, []byte("客户端发送消息:"+string(p)))

		fmt.Println(m, string(p))
	}
	defer conn.Close()
	log.Println("websocket服务关闭")
}

func main() {

	// 初始化配置
	configInit()
	db := dbInit()
	tableInit(db)

	// 启动websocket服务
	go WebSocketInit()

	// 启动grpc服务
	GRPCInit()
}

func WebSocketInit() {
	http.HandleFunc("/", handler)
	_ = http.ListenAndServe(":8088", nil)
}

func GRPCInit() {
	// 开启端口监听
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Printf("监听失败: %v", err)
	}
	// 初始化PVE客户端
	pveClient := pve.NewPVEClient(
		viper.GetString("pve.api_url"),
		viper.GetString("pve.api_token"),
		viper.GetString("pve.node_name"),
	)

	// 注册grpc服务
	grpcServer := grpc.NewServer()
	pvm.RegisterVMManagerServer(grpcServer, service.NewVMServer(pveClient))

	// 启动服务
	if err = grpcServer.Serve(listen); err != nil {
		log.Printf("启动服务失败: %v", err)
	}
}

func configInit() {
	env := "local"
	viper.SetConfigName(env)                                // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")                             // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./internal/infrastructure/config") // 查找配置文件所在的路径

	viper.AddConfigPath(".")    // 还可以在工作目录中查找配置
	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {             // 处理读取配置文件的错误
		log.Printf("config配置错误: %s \n", err)
	}
}

func dbInit() *gorm.DB {
	type MySQLConfig struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
	}
	var config MySQLConfig
	if err := viper.UnmarshalKey("database", &config); err != nil {
		log.Printf("无法解码为结构database: %s \n", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 取消外键约束
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "vm_", // 设置表前缀
		},
	})
	if err != nil {
		log.Printf("连接database失败: %v", err)
	}

	return db
}

func tableInit(db *gorm.DB) {
	err := db.AutoMigrate(
		&entity.VirtualMachine{},
	)
	if err != nil {
		log.Printf("初始化table失败: %v", err)
	}
}
