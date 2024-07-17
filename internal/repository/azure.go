package repository

import (
	"context"

	"eyes/internal/domain"
	"eyes/internal/repository/dao"
)

type AzureRepository interface {
	Create(ctx context.Context, Azure domain.Azure) (int64, error)
	Update(ctx context.Context, Azure domain.Azure) error
	Save(ctx context.Context, Azure domain.Azure) error
	Select(ctx context.Context) (domain.UserInfo, error)
}

func NewAzureRepository(dao dao.AzureDAO) AzureRepository {
	return &azureRepository{
		dao: dao,
	}
}

type azureRepository struct {
	dao dao.AzureDAO
}

func (a *azureRepository) Create(ctx context.Context, Azure domain.Azure) (int64, error) {
	return a.dao.Insert(ctx, dao.Azure{
		Author:  Azure.Author,
		Title:   Azure.Title,
		Content: Azure.Content,
	})
}

func (a *azureRepository) Save(ctx context.Context, Azure domain.Azure) error {
	return a.dao.Save(ctx, dao.Azure{
		Id:      Azure.ID,
		Author:  Azure.Author,
		Title:   Azure.Title,
		Content: Azure.Content,
	})
}

func (a *azureRepository) Update(ctx context.Context, Azure domain.Azure) error {
	return a.dao.Update(ctx, dao.Azure{
		Id:      Azure.ID,
		Author:  Azure.Author,
		Title:   Azure.Title,
		Content: Azure.Content,
	})
}

func (a *azureRepository) Select(ctx context.Context) (domain.UserInfo, error) {
	u, err := a.dao.Select(ctx)
	return domain.UserInfo{
		ID:         u.ID,
		RealName:   u.RealName,
		Email:      u.Email,
		CreateDate: u.CreateDate,
	}, err
}
