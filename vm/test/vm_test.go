package test

import (
	"github.com/spf13/viper"
	"log"
	"testing"
	pveApp "vm/internal/application/pve"
	etcd2 "vm/internal/infrastructure/etcd"
	"vm/internal/infrastructure/inits"
)

func TestGetEtcd(t *testing.T) {

	// 读取配置文件
	env := "local"
	viper.SetConfigName(env)                                 // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")                              // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("../internal/infrastructure/config") // 查找配置文件所在的路径
	viper.AddConfigPath(".")                                 // 还可以在工作目录中查找配置

	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {             // 处理读取配置文件的错误
		log.Panicf("config配置错误: %s \n", err)
	}

	// etcd连接
	etcd := inits.EtcdInit()

	// 基础设施层

	// 创建websocket连接指针
	etcdConn := etcd2.NewEtcdService(etcd)

	// 应用层：
	pve := pveApp.NewPVEServer(nil, nil, etcdConn)

	v, err := pve.GetEtcd("testkey")
	if err != nil {
		log.Println("获取etcd失败", err)
	} else {
		log.Println("获取etcd成功", v)
	}
}
