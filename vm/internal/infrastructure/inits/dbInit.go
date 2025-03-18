package inits

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

func DbInit() (db *gorm.DB) {
	// 读取配置文件
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

	// 参数设置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)
	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 取消外键约束
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "vm_", // 设置表前缀
		},
	})
	if err != nil {
		log.Panicf("连接database失败: %v", err)
	}

	return db
}
