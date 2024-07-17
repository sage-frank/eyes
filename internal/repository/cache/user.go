package cache

import (
	"context"
	"encoding/json"
	"eyes/internal/domain"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type UserCache interface {
	Set(ctx context.Context, u domain.User) error
}

type redisUserCache struct {
	client     redis.Cmdable
	pattern    string
	expiration time.Duration
}

func NewRedisUserCache(client redis.Cmdable) UserCache {
	return &redisUserCache{
		client:     client,
		pattern:    "user_%d",
		expiration: time.Minute * 15,
	}
}

func (r redisUserCache) Set(ctx context.Context, u domain.User) error {
	values, err := json.Marshal(u)
	if err != nil {
		return err
	}
	_, err = r.client.Set(fmt.Sprintf(r.pattern, u.ID), values, r.expiration).Result()
	return err
}
