package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"eyes/internal/domain"

	"github.com/go-redis/redis"
)

type AlertCache interface {
	Set(ctx context.Context, m domain.Monitor) error
	Get(ctx context.Context, m domain.Monitor) (*domain.Monitor, error)
	SetCnt(ctx context.Context, m int64) error
	GetCnt(ctx context.Context, m domain.Monitor) (int64, error)
}

type redisMonitorCache struct {
	client     redis.Cmdable
	pattern    string
	patternINT string
	expiration time.Duration
}

//func NewRedisMonitorCache(client redis.Cmdable) cache.Cache {
//	return &redisMonitorCache{
//		client:     client,
//		pattern:    "monitor_%s",
//		patternINT: "monitor_%d",
//		expiration: time.Minute * 15,
//	}
//}

func (r *redisMonitorCache) Set(_ context.Context, m domain.Monitor) error {
	values, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = r.client.Set(fmt.Sprintf(r.pattern, m.ID), values, r.expiration).Result()
	return err
}

func (r *redisMonitorCache) SetCnt(_ context.Context, m int64) error {
	_, err := r.client.Set(fmt.Sprintf(r.patternINT, m), m, r.expiration).Result()
	return err
}

func (r *redisMonitorCache) Get(_ context.Context, m domain.Monitor) (*domain.Monitor, error) {
	ret, err := r.client.Get(fmt.Sprintf(r.pattern, m.ID)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, err // redis.Nil
		}
		return nil, fmt.Errorf("r.client.Get(fmt.Sprintf(r.pattern, m.ID)).Result(): %w", err)
	}
	resp := new(domain.Monitor)
	err = json.Unmarshal(ret, resp)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal(ret, resp): %w", err)
	}
	return resp, nil
}

func (r *redisMonitorCache) GetCnt(_ context.Context, m domain.Monitor) (int64, error) {
	return r.client.Get(fmt.Sprintf(r.patternINT, 1)).Int64()
}
