package inits

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
	"time"
)

func RedisInit() *redis.Client {
	type RedisConfig struct {
		Host         string
		Port         int
		Password     string
		DB           int
		PoolSize     int
		MinIdleConns int
		DialTimeout  time.Duration
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
		MaxRetries   int
	}
	var config RedisConfig
	if err := viper.UnmarshalKey("redis", &config); err != nil {
		log.Printf("无法解码为结构 redis: %s \n", err)
	}
	
	// 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
		MaxRetries:   config.MaxRetries,
	})
	
	// 测试 Redis 连接
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("连接 redis 失败: %v", err)
	}
	
	return rdb
}
