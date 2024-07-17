package dao

import (
	"context"
	"time"

	sf "eyes/utility"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserDAO interface {
	Insert(context.Context, User) (int64, error)
	Update(context.Context, User) error
	Save(context.Context, User) error
}

type userDAO struct {
	db   *gorm.DB
	node sf.ISFNode
}

func NewUserDAO(db *gorm.DB, node sf.ISFNode) UserDAO {
	return &userDAO{
		db:   db,
		node: node,
	}
}

func (ud userDAO) Insert(ctx context.Context, u User) (int64, error) {
	now := time.Now().UnixMilli()
	u.ID = ud.node.GenID()
	u.Ctime = now
	u.Utime = now
	err := ud.db.WithContext(ctx).Create(&u).Error
	return u.ID, err
}

func (ud userDAO) Update(ctx context.Context, u User) error {
	u.Utime = time.Now().UnixMilli()
	return ud.db.WithContext(ctx).Updates(&u).Error
}

func (ud userDAO) Save(ctx context.Context, u User) error {
	u.Utime = time.Now().UnixMilli()
	u.Ctime = time.Now().UnixMilli()
	return ud.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{"utime", "name", "content"}),
		}).Create(u).Error
}

type User struct {
	ID int64 `gorm:"primaryKey,autoIncrement"`
	// 存的是作者的用户 ID
	Name     string `gorm:"not null" form:"name"`
	PassWord string `form:"name"`
	Content  string `form:"content"`
	Ctime    int64
	Utime    int64
}
