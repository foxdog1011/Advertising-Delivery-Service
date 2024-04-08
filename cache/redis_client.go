// cache/redis_client.go
package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	Client *redis.Client
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		// 如果键不存在，redis.Nil 被返回
		return "", nil
	} else if err != nil {
		return "", err // 返回遇到的错误
	}
	return val, nil // 返回键對應的值
}
