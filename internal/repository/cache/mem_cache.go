package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type MemcachedCache struct {
	client *memcache.Client
}

func NewMemcachedCache(server string) Cache {
	client := memcache.New(server)
	return &MemcachedCache{client: client}
}

func (m *MemcachedCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return m.client.Set(&memcache.Item{
		Key:        key,
		Value:      bytes,
		Expiration: int32(expiration.Seconds()),
	})
}

func (m *MemcachedCache) Get(ctx context.Context, key string, dest interface{}) error {
	item, err := m.client.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(item.Value, dest)
}

func (m *MemcachedCache) SetCnt(ctx context.Context, key string, cnt int64, expiration time.Duration) error {
	value := fmt.Sprintf("%d", cnt)
	return m.client.Set(&memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: int32(expiration.Seconds()),
	})
}

func (m *MemcachedCache) GetCnt(ctx context.Context, key string) (int64, error) {
	item, err := m.client.Get(key)
	if err != nil {
		return 0, err
	}
	var cnt int64
	err = json.Unmarshal(item.Value, &cnt)
	return cnt, err
}
