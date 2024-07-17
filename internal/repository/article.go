package repository

import (
	"context"
	"fmt"
	"time"

	"eyes/internal/domain"
	"eyes/internal/repository/cache"
	"eyes/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
	Save(ctx context.Context, article domain.Article) error
}

func NewArticleRepository(dao dao.ArticleDAO, cache cache.Cache) ArticleRepository {
	return &articleRepository{
		dao:   dao,
		cache: cache,
	}
}

type articleRepository struct {
	dao   dao.ArticleDAO
	cache cache.Cache
}

func (a *articleRepository) Create(ctx context.Context, article domain.Article) (int64, error) {
	return a.dao.Insert(ctx, dao.Article{
		Author:  article.Author,
		Title:   article.Title,
		Content: article.Content,
	})
}

func (a *articleRepository) Save(ctx context.Context, article domain.Article) error {
	err := a.dao.Save(ctx, dao.Article{
		Id:      article.ID,
		Author:  article.Author,
		Title:   article.Title,
		Content: article.Content,
	})
	if err != nil {
		return err
	}
	return a.cache.Set(ctx, fmt.Sprintf("article-%d", article.ID), article, time.Second*60)
}

func (a *articleRepository) Update(ctx context.Context, article domain.Article) error {
	return a.dao.Update(ctx, dao.Article{
		Id:      article.ID,
		Author:  article.Author,
		Title:   article.Title,
		Content: article.Content,
	})
}
