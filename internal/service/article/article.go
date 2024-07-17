package article

import (
	"context"

	"eyes/internal/domain"
	"eyes/internal/repository"
)

type ArtService interface {
	Save(ctx context.Context, article domain.Article) (int64, error)
	Publish(ctx context.Context, article domain.Article) error
}

func NewArticleService(bRepo repository.ArticleRepository, cRepo repository.ArticleRepository) ArtService {
	return &articleService{
		bRepo: bRepo,
		cRepo: cRepo,
	}
}

type articleService struct {
	// 代表制作库
	bRepo repository.ArticleRepository
	cRepo repository.ArticleRepository
}

type cartService struct {
	cRepo repository.ArticleRepository
}

// Publish 真正的业务逻辑，在这里
func (a *articleService) Publish(ctx context.Context, article domain.Article) error {
	// 这是更新
	id, err := a.Save(ctx, article)
	if err != nil {
		return err
	}
	article.ID = id
	// 在实际场景中，你这里要加监控和重试
	return a.cRepo.Save(ctx, article)
}

// Save 真正的业务逻辑，在这里
func (a *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
	// 这是更新
	if article.ID > 0 {
		return article.ID, a.bRepo.Update(ctx, article)
	}
	// 这是新建
	return a.bRepo.Create(ctx, article)
}

// Save cartService 的实现
func (c cartService) Save(ctx context.Context, article domain.Article) (int64, error) {
	return c.cRepo.Create(ctx, article)
}

// Publish cartService 的实现
func (c cartService) Publish(ctx context.Context, article domain.Article) error {
	return c.cRepo.Save(ctx, article)
}
