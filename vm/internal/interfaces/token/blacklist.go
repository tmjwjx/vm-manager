package token

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// BlacklistKey 定义黑名单的 Redis Sorted Set 键名，与 siwu 项目一致
const BlacklistKey = "blacklist:tokens"

// IsTokenBlacklisted 检查 token 是否在黑名单中
func IsTokenBlacklisted(rdb *redis.Client, tokenString string) (bool, error) {
	ctx := context.Background()
	
	// 获取 token 的过期时间戳
	score, err := rdb.ZScore(ctx, BlacklistKey, tokenString).Result()
	if errors.Is(err, redis.Nil) {
		// token 不存在
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to check blacklist: %v", err)
	}
	
	// 检查是否过期
	currentTime := time.Now().Unix()
	if currentTime > int64(score) {
		// 已过期，异步移除
		go rdb.ZRem(ctx, BlacklistKey, tokenString)
		return false, nil
	}
	
	return true, nil
}
