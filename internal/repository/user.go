package repository

import (
	"context"

	"eyes/internal/domain"
	"eyes/internal/repository/cache"
	"eyes/internal/repository/dao"
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) (int64, error)
	Update(ctx context.Context, u domain.User) error
	Save(ctx context.Context, u domain.User) error
}

func NewUserRepository(dao dao.UserDAO, cache cache.Cache) UserRepository {
	return &userRepository{
		dao:   dao,
		cache: cache,
	}
}

type userRepository struct {
	dao   dao.UserDAO
	cache cache.Cache
}

func (u2 userRepository) Create(ctx context.Context, u domain.User) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (u2 userRepository) Update(ctx context.Context, u domain.User) error {
	// TODO implement me
	panic("implement me")
}

func (u2 userRepository) Save(ctx context.Context, u domain.User) error {
	// TODO implement me
	panic("implement me")
}
