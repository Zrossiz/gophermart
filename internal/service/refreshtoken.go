package service

import "github.com/Zrossiz/gophermart/internal/model"

type RefreshTokenService struct {
	db *TokenStorage
}

type TokenStorage interface {
	Create(userID int64, token string) (bool, error)
	DeleteByToken(token string) (bool, error)
	DeleteTokenByUser(userID int64) (bool, error)
	GetTokenByToken(token string) (*model.RefreshToken, error)
	GetTokenByUser(userID int64) (*model.RefreshToken, error)
}

func NewRefreshTokenService(tokenStorage TokenStorage) *RefreshTokenService {
	return &RefreshTokenService{
		db: &tokenStorage,
	}
}
