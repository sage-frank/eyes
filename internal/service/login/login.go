package login

import (
	"context"

	"eyes/internal/domain"
	"eyes/internal/repository"
)

type LService interface {
	LoginByPass(ctx context.Context, user domain.User) (int64, error)
	LoginByMobile(ctx context.Context, user domain.User) error
}

type loginService struct {
	uRepo repository.UserRepository
}

func (l loginService) LoginByPass(ctx context.Context, user domain.User) (int64, error) {
	// TODO implement me
	panic("implement me")
}

func (l loginService) LoginByMobile(ctx context.Context, user domain.User) error {
	// TODO implement me
	panic("implement me")
}

func NewLoginService(uRepo repository.UserRepository) LService {
	return &loginService{
		uRepo: uRepo,
	}
}

//
//// Publish 真正的业务逻辑，在这里
//func (a *articleService) Publish(ctx context.Context, article domain.Article) error {
//	// 这是更新
//	id, err := a.Save(ctx, article)
//	if err != nil {
//		return err
//	}
//	article.ID = id
//	// 在实际场景中，你这里要加监控和重试
//	return a.cRepo.Save(ctx, article)
//}
//
//// Save 真正的业务逻辑，在这里
//func (a *articleService) Save(ctx context.Context, article domain.Article) (int64, error) {
//	// 这是更新
//	if article.ID > 0 {
//		return article.ID, a.bRepo.Update(ctx, article)
//	}
//	// 这是新建
//	return a.bRepo.Create(ctx, article)
//}
