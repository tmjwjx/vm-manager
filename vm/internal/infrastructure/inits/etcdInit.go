package inits

import (
	"github.com/spf13/viper"
	"go.etcd.io/etcd/client/v3"
	"log"
	"strconv"
	"time"
)

func EtcdInit() (cli *clientv3.Client) {
	type EtcdConfig struct {
		Host        string
		Port        int
		DialTimeout time.Duration `mapstructure:"dial_timeout"`
	}
	var config EtcdConfig
	if err := viper.UnmarshalKey("etcd", &config); err != nil {
		log.Printf("无法解码为结构 etcd: %s \n", err)
	}

	// 连接 etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{config.Host + ":" + strconv.Itoa(config.Port)},
		DialTimeout: config.DialTimeout,
	})
	if err != nil {
		log.Fatal("连接失败：", err)
	}

	// 测试连接

	log.Println("成功连接 etcd！")

	return cli
}
