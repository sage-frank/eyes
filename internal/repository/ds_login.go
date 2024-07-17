package repository

import (
	"context"
	"errors"
	"reflect"

	"eyes/internal/repository/cache"

	"eyes/internal/domain"
	"eyes/internal/repository/dao"
	"gorm.io/gorm"
)

type DsUserRepository interface {
	Create(context.Context, domain.DSUser) (string, error)
	Update(context.Context, domain.DSUser) error
	Save(context.Context, domain.DSUser) (string, error)
	Select(context.Context, domain.DSUser) (domain.DSUser, error)
}

type dsLoginRepository struct {
	dao   dao.DSUserDAO
	cache cache.Cache
}

func (d dsLoginRepository) Create(ctx context.Context, user domain.DSUser) (string, error) {
	u := dao.DSUser{}
	copyStructFields(user, u)
	return d.dao.Insert(ctx, u)
}

func (d dsLoginRepository) Update(ctx context.Context, user domain.DSUser) error {
	u := dao.DSUser{}
	copyStructFields(user, u)
	return d.dao.Update(ctx, u)
}

func (d dsLoginRepository) Save(ctx context.Context, ds domain.DSUser) (string, error) {
	return d.dao.Insert(ctx, dao.DSUser{
		Nickname:    ds.Nickname,
		Username:    ds.Username,
		Password:    ds.Password,
		Salt:        ds.Salt,
		Phone:       ds.Phone,
		Email:       ds.Email,
		Avatar:      ds.Avatar,
		Description: ds.Description,
	})
}

func (d dsLoginRepository) Select(ctx context.Context, user domain.DSUser) (domain.DSUser, error) {
	u, err := d.dao.Select(ctx, dao.DSUser{
		Phone:    user.Phone,
		Email:    user.Email,
		Username: user.Username,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.DSUser{}, nil
		} else {
			return domain.DSUser{}, err
		}
	}
	return u, nil
}

func NewDsLoginRepository(dao dao.DSUserDAO, cache cache.Cache) DsUserRepository {
	return &dsLoginRepository{
		dao:   dao,
		cache: cache,
	}
}

var _ DsUserRepository = &dsLoginRepository{}

func copyStructFields(src interface{}, dst interface{}) {
	srcVal := reflect.ValueOf(src).Elem()
	dstVal := reflect.ValueOf(dst).Elem()
	for i := 0; i < srcVal.NumField(); i++ {
		field := srcVal.Field(i)
		dstVal.Field(i).Set(field)
	}
}
