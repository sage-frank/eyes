package alert

import (
	"context"
	"eyes/internal/domain"
	"eyes/internal/repository"
)

type MonitorService interface {
	Count(ctx context.Context, monitor domain.Monitor) (int64, error)
	Detail(ctx context.Context, monitor domain.Monitor) (*domain.Monitor, error)
	List(ctx context.Context, page int64, size int64, monitor domain.Monitor) ([]*domain.Monitor, int64, error)
	Query(ctx context.Context, page, size int64, args string, monitor domain.Monitor) ([]*domain.Monitor, int64, error)
}

func (a articleService) Count(ctx context.Context, monitor domain.Monitor) (int64, error) {
	return a.repo.Count(ctx, monitor)
}

func (a articleService) Detail(ctx context.Context, monitor domain.Monitor) (*domain.Monitor, error) {
	return a.repo.Detail(ctx, monitor)
}

func (a articleService) List(ctx context.Context, page int64, size int64, monitor domain.Monitor) ([]*domain.Monitor, int64, error) {
	return a.repo.List(ctx, page, size, monitor)
}

func (a articleService) Query(ctx context.Context, page, size int64, args string, monitor domain.Monitor) ([]*domain.Monitor, int64, error) {
	return a.repo.Query(ctx, page, size, args, monitor)
}

func NewMonitorService(repo repository.MonitorRepository) MonitorService {
	return &articleService{
		repo: repo,
	}
}

type articleService struct {
	// 代表制作库
	repo repository.MonitorRepository
}
