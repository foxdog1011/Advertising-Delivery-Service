// 在 cache/redis.go 中

package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

type RedisClientInterface interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error) // 增加 Get 方法
}

func InitRedis() RedisClientInterface {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "cache:6379", // 或者 "localhost:6379" 如果您直接在宿主机上运行 Redis
		Password: "",           // 如果设置了密码
		DB:       0,            // 通常使用默认的 DB
	})

	return &RedisClient{Client: rdb}
}
