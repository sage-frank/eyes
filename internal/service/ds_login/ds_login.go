package ds_login

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"eyes/internal/web/common"

	"eyes/internal/domain"
	"eyes/internal/repository"
)

type DsLoginServer interface {
	LoginByPass(ctx context.Context, user domain.DSUser) (string, error)
	LoginByMobile(ctx context.Context, user domain.DSUser) (string, error)
	Select(ctx context.Context, user domain.DSUser) (domain.DSUser, error)
	Save(ctx context.Context, user domain.DSUser) (string, error)
}

type dsLoginService struct {
	dsRepo repository.DsUserRepository
}

func (l dsLoginService) Login(ctx context.Context, user domain.DSUser) (string, error) {
	if user.Phone != "" {
		return l.LoginByMobile(ctx, user)
	}
	return l.LoginByPass(ctx, user)
}

func (l dsLoginService) LoginByPass(ctx context.Context, user domain.DSUser) (string, error) {
	u, err := l.dsRepo.Select(ctx, user)
	if err != nil {
		return "", fmt.Errorf(".dsRepo.Select(ctx, user): %w", err)
	}
	if verifyPassword(user.Password, u) {
		return u.ID, nil
	} else {
		return "", common.ErrPassword
	}
}

func (l dsLoginService) Select(ctx context.Context, user domain.DSUser) (domain.DSUser, error) {
	return l.dsRepo.Select(ctx, user)
}

func (l dsLoginService) Save(ctx context.Context, user domain.DSUser) (string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", err
	}
	user.Salt = base64.StdEncoding.EncodeToString(salt)
	user.Password = hashPassword(user.Password, salt)
	return l.dsRepo.Save(ctx, user)
}

func (l dsLoginService) LoginByMobile(ctx context.Context, user domain.DSUser) (string, error) {
	// TODO implement me
	panic("implement me")
}

var _ DsLoginServer = &dsLoginService{}

func NewLoginService(dsRepo repository.DsUserRepository) DsLoginServer {
	return &dsLoginService{
		dsRepo: dsRepo,
	}
}

func verifyPassword(password string, u domain.DSUser) bool {
	salt, err := base64.StdEncoding.DecodeString(u.Salt) // 解码 Base64 编码的盐值
	if err != nil {
		return false
	}
	hashedPassword := hashPassword(password, salt)
	return hashedPassword == u.Password
}

func hashPassword(password string, salt []byte) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	hash.Write(salt)
	hashedPassword := hash.Sum(nil)
	return base64.StdEncoding.EncodeToString(hashedPassword)
}

func generateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}
