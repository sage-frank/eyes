package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AzureDAO interface {
	Insert(context.Context, Azure) (int64, error)
	Update(context.Context, Azure) error
	Save(context.Context, Azure) error
	Select(context.Context) (UserInfo, error)
}

func NewAzureDAO(db *gorm.DB) AzureDAO {
	return &azureDAO{
		db: db,
	}
}

type azureDAO struct {
	db *gorm.DB
}

func (a *azureDAO) Insert(ctx context.Context, Azure Azure) (int64, error) {
	now := time.Now().UnixMilli()
	Azure.Ctime = now
	Azure.Utime = now
	err := a.db.WithContext(ctx).Create(&Azure).Error
	return Azure.Id, err
}

func (a *azureDAO) Save(ctx context.Context, Azure Azure) error {
	Azure.Utime = time.Now().UnixMilli()
	Azure.Ctime = time.Now().UnixMilli()
	return a.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{"utime", "content", "title"}),
		}).Create(Azure).Error
}

func (a *azureDAO) Update(ctx context.Context, Azure Azure) error {
	Azure.Utime = time.Now().UnixMilli()
	return a.db.WithContext(ctx).Updates(&Azure).Error
}

func (a *azureDAO) Select(ctx context.Context) (UserInfo, error) {
	userInfo := UserInfo{}
	email := ctx.Value("email")
	err := a.db.WithContext(ctx).Where("email=?", email).First(&userInfo).Error
	return userInfo, err
}

type Azure struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 存的是作者的用户 ID
	Author  int64  `gorm:"not null"`
	Title   string `form:"title"`
	Content string `form:"content"`
	Ctime   int64
	Utime   int64
}

type UserInfo struct {
	ID         int64     `gorm:"primaryKey,autoIncrement"`
	RealName   string    `gorm:"real_name"`
	Email      string    `gorm:"email"`
	openID     string    `gorm:"open_id"`
	valid      int64     `gorm:"valid"`
	CreateDate time.Time `gorm:"create_date"`
}
