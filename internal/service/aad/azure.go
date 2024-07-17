package aad

import (
	"context"

	"eyes/internal/domain"
	"eyes/internal/repository"
)

type AzureService interface {
	Save(ctx context.Context, azure domain.Azure) (int64, error)
	Select(ctx context.Context) (domain.UserInfo, error)
}

func NewAzureService(bRepo repository.AzureRepository) AzureService {
	return &azureService{
		bRepo: bRepo,
	}
}

type azureService struct {
	// 代表制作库
	bRepo repository.AzureRepository
}

// Save 真正的业务逻辑，在这里
func (a *azureService) Save(ctx context.Context, azure domain.Azure) (int64, error) {
	// 这是更新
	if azure.ID > 0 {
		return azure.ID, a.bRepo.Update(ctx, azure)
	}
	// 这是新建
	return a.bRepo.Create(ctx, azure)
}

func (a *azureService) Select(ctx context.Context) (domain.UserInfo, error) {
	return a.bRepo.Select(ctx)
}
