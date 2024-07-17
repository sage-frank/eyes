package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client redis.Cmdable
}

func NewRedisCache(client redis.Cmdable) Cache {
	return &RedisCache{client: client}
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	values, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, values, expiration).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(val, dest)
}

func (r *RedisCache) SetCnt(ctx context.Context, key string, cnt int64, expiration time.Duration) error {
	return r.client.Set(ctx, key, cnt, expiration).Err()
}

func (r *RedisCache) GetCnt(ctx context.Context, key string) (int64, error) {
	return r.client.Get(ctx, key).Int64()
}
