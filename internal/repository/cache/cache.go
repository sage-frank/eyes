package cache

import (
	"context"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	SetCnt(ctx context.Context, key string, cnt int64, expiration time.Duration) error
	GetCnt(ctx context.Context, key string) (int64, error)
}
