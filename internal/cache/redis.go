package cache

import (
	"context"
	"fintrack/config"
	"time"

	"github.com/redis/go-redis/v9"
)

// 定義緩存接口
type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Close() error
}

// Redis 實現 Cache 接口
type Redis struct {
	client *redis.Client
}

// NewCache 創建並返回一個 Cache 實例，並使用 Redis 作為緩存
func NewCache(config config.RedisConfig) Cache {
	return NewRedis(config)
}

// NewRedisCache 使用 RedisConfig 來初始化 RedisCache
func NewRedis(config config.RedisConfig) *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})
	return &Redis{client: rdb}
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *Redis) Close() error {
	return r.client.Close()
}
