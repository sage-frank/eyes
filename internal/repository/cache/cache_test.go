package cache

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"eyes/internal/domain"

	"github.com/go-redis/redis/v8"
)

func TestRedisMonitorCache_Get(t *testing.T) {
	// 使用 Redis 缓存
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:26379",
	})
	redisCache := NewRedisCache(redisClient)
	useCache(redisCache)

	// 使用 Memcached 缓存
	memcachedCache := NewMemcachedCache("localhost:11211")
	useCache(memcachedCache)
}

func useCache(c Cache) {
	ctx := context.Background()

	// 设置 Monitor 缓存
	monitor := domain.Monitor{ID: "123", FlowID: 100}
	err := c.Set(ctx, "monitor_123", monitor, time.Minute*15)
	if err != nil {
		fmt.Println("Error setting cache:", err)
		return
	}

	// 获取 Monitor 缓存
	var cachedMonitor domain.Monitor
	err = c.Get(ctx, "monitor_123", &cachedMonitor)
	if err != nil {
		fmt.Println("Error getting cache:", err)
		return
	}
	fmt.Println("Cached Monitor:", cachedMonitor)
	fmt.Println(cachedMonitor.ID == "123" && cachedMonitor.FlowID == 100)
	if cachedMonitor.ID != "123" || cachedMonitor.FlowID != 100 {
		log.Fatalln("id != 123 or FlowID != 100")
	}

	// 设置计数缓存
	err = c.SetCnt(ctx, "monitor_count", 42, time.Minute*15)
	if err != nil {
		fmt.Println("Error setting cache count:", err)
		return
	}

	// 获取计数缓存
	cnt, err := c.GetCnt(ctx, "monitor_count")
	if err != nil {
		fmt.Println("Error getting cache count:", err)
		return
	}
	fmt.Println("Cached Count:", cnt)
}
