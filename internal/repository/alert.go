package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"eyes/internal/domain"
	"eyes/internal/repository/cache"
	"eyes/internal/repository/es"
	sf "eyes/utility"

	"github.com/go-redis/redis"
)

var COUNT_PREFIX = "CNT"

type MonitorRepository interface {
	Count(ctx context.Context, m domain.Monitor) (int64, error)
	Detail(ctx context.Context, m domain.Monitor) (*domain.Monitor, error)
	List(ctx context.Context, page, size int64, m domain.Monitor) ([]*domain.Monitor, int64, error)
	Query(ctx context.Context, page, size int64, args string, m domain.Monitor) ([]*domain.Monitor, int64, error)
}

func (r *monitorRepository) Count(ctx context.Context, m domain.Monitor) (int64, error) {
	cnt, err1 := r.cache.GetCnt(ctx, m.ID)
	if err1 != nil {

		cnt, err := r.es.Count(ctx, m)
		if err != nil {
			return 0, fmt.Errorf("r.es.Count(ctx, m): %w", errors.Join(err1, err))
		}

		err = r.cache.SetCnt(ctx, fmt.Sprintf("%s-%s", COUNT_PREFIX, m.ID), cnt, time.Second*60)
		if err != nil {
			return 0, fmt.Errorf("r.cache.SetCnt(ctx, cnt): %w", err)
		}
		return cnt, nil
	} else {
		return cnt, nil
	}
}

func (r *monitorRepository) Detail(ctx context.Context, m domain.Monitor) (*domain.Monitor, error) {
	var val domain.Monitor
	err := r.cache.Get(ctx, m.ID, &val)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("r.cache.Get(ctx, m): %w", err)
		} else {
			val, err := r.es.Detail(ctx, m)
			if err != nil {
				return nil, fmt.Errorf("r.es.Detail(ctx, m): %w", err)
			}
			err1 := r.cache.Set(ctx, val.ID, *val, time.Second*60)
			if err1 != nil {
				return nil, fmt.Errorf(" r.cache.Set(ctx, *v) %w", err1)
			}
			return val, nil
		}
	}
	return &val, nil
}

func (r *monitorRepository) List(ctx context.Context, page, size int64, m domain.Monitor) ([]*domain.Monitor, int64, error) {
	return r.es.List(ctx, page, size, m)
}

func (r *monitorRepository) Query(ctx context.Context, page, size int64, args string, m domain.Monitor) ([]*domain.Monitor, int64, error) {
	return r.es.Query(ctx, page, size, args, m)
}

func NewMonitorRepository(es es.MonitorESDAO, cache cache.Cache, node sf.ISFNode) MonitorRepository {
	return &monitorRepository{
		es:    es,
		cache: cache,
		node:  node,
	}
}

type monitorRepository struct {
	es    es.MonitorESDAO // interface
	cache cache.Cache     // interface
	node  sf.ISFNode      // pointer
}
