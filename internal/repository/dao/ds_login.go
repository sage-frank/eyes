package dao

import (
	"context"
	"database/sql"
	"fmt"

	"eyes/internal/domain"

	"eyes/utility"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DSUserDAO interface {
	Insert(context.Context, DSUser) (string, error)
	Update(context.Context, DSUser) error
	Save(context.Context, DSUser) error
	Select(context.Context, DSUser) (domain.DSUser, error)
}

type dsUserDAO struct {
	db *gorm.DB
}

func NewDsUserDAO(db *gorm.DB) DSUserDAO {
	return &dsUserDAO{
		db: db,
	}
}

func (ud dsUserDAO) Insert(ctx context.Context, u DSUser) (string, error) {
	u.ID = utility.GenID()
	err := ud.db.WithContext(ctx).Create(&u).Error
	return u.ID, err
}

func (ud dsUserDAO) Update(ctx context.Context, u DSUser) error {
	return ud.db.WithContext(ctx).Updates(&u).Error
}

func (ud dsUserDAO) Save(ctx context.Context, u DSUser) error {
	return ud.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{"username", "phone", "email"}),
		}).Create(u).Error
}

func (ud dsUserDAO) Select(ctx context.Context, u DSUser) (domain.DSUser, error) {
	userInfo := domain.DSUser{}
	var err error
	if u.Email != "" {
		err = ud.db.WithContext(ctx).
			Where("email = ?", u.Email).
			First(&userInfo).Error
	} else if u.Phone != "" {
		err = ud.db.WithContext(ctx).
			Where("phone = ?", u.Phone).
			First(&userInfo).Error
	} else if u.Username != "" {
		err = ud.db.WithContext(ctx).
			Where("username = ?", u.Username).
			First(&userInfo).Error
	} else {
		err = fmt.Errorf("no login method provided")
	}
	return userInfo, err
}

type DSUser struct {
	ID              string       `json:"id"`                         // 用户唯一标识符
	Nickname        string       `json:"nickname,omitempty"`         // 用户昵称
	Username        string       `json:"username,omitempty"`         // 用户名
	Password        string       `json:"password,omitempty"`         // 密码
	Salt            string       `json:"salt,omitempty"`             // 盐
	Phone           string       `json:"phone,omitempty"`            // 用户电话号码
	Email           string       `json:"email,omitempty"`            // 用户电子邮箱
	Avatar          string       `json:"avatar,omitempty"`           // 用户头像URL
	Description     string       `json:"description,omitempty"`      // 用户描述或简介
	IsHot           int          `json:"is_hot,omitempty"`           // 是否是热门用户
	IsPro           int          `json:"is_pro,omitempty"`           // 假设这是专业用户标识
	CreatedIP       string       `json:"created_ip,omitempty"`       // 创建时的IP地址
	CreatedLocation string       `json:"created_location,omitempty"` // 创建时的位置
	CreatedAt       sql.NullTime `json:"created_at"`                 // 创建时间
	UpdatedAt       sql.NullTime `json:"updated_at"`                 // 最后更新时间
	Status          int          `json:"status,omitempty"`           // 用户状态
}
