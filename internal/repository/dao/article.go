package dao

import (
	"context"
	sf "eyes/utility"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, article Article) (int64, error)
	Update(ctx context.Context, article Article) error
	Save(ctx context.Context, article Article) error
}

func NewArticleDAO(db *gorm.DB, node sf.ISFNode) ArticleDAO {
	return &articleDAO{
		db:   db,
		node: node,
	}
}

type articleDAO struct {
	db   *gorm.DB
	node sf.ISFNode
}

func (a *articleDAO) Insert(ctx context.Context, article Article) (int64, error) {
	now := time.Now().UnixMilli()

	article.Id = a.node.GenID()
	article.Ctime = now
	article.Utime = now
	err := a.db.WithContext(ctx).Create(&article).Error
	return article.Id, err
}

func (a *articleDAO) Save(ctx context.Context, article Article) error {
	article.Utime = time.Now().UnixMilli()
	article.Ctime = time.Now().UnixMilli()
	return a.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{"utime", "content", "title"}),
		}).Create(article).Error
}

func (a *articleDAO) Update(ctx context.Context, article Article) error {
	article.Utime = time.Now().UnixMilli()
	return a.db.WithContext(ctx).Updates(&article).Error
}

type Article struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 存的是作者的用户 ID
	Author  int64  `gorm:"not null"`
	Title   string `form:"title"`
	Content string `form:"content"`
	Ctime   int64
	Utime   int64
}
