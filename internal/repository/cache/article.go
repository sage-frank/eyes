package cache

import (
	"context"
	"encoding/json"
	"eyes/internal/domain"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type ArtCache interface {
	Set(ctx context.Context, a domain.Article) error
}

type redisArtCache struct {
	client     redis.Cmdable
	pattern    string
	expiration time.Duration
}

func (r redisArtCache) Set(ctx context.Context, a domain.Article) error {
	values, err := json.Marshal(a)
	if err != nil {
		return err
	}
	_, err = r.client.Set(fmt.Sprintf(r.pattern, a.ID), values, r.expiration).Result()
	return err
}

func NewRedisArticleCache(client redis.Cmdable) ArtCache {
	return &redisArtCache{
		client:     client,
		pattern:    "article_%d",
		expiration: time.Minute * 15,
	}
}
