package service

import (
	"github.com/Zrossiz/gophermart/internal/config"
	"go.uber.org/zap"
)

type Service struct {
	UserService           *UserService
	BalanceHistoryService *BalanceHistoryService
	RefreshTokenService   *RefreshTokenService
	OrderService          *OrderService
	StatusService         *StatusService
}

type Storage struct {
	BalanceHistoryStorage BalanceHistoryStorage
	UserStorage           UserStorage
	OrderStorage          OrderStorage
	TokenStorage          TokenStorage
	StatusStorage         StatusStorage
}

func New(db Storage, cfg *config.Config, log *zap.Logger) *Service {
	return &Service{
		UserService: NewUserService(
			db.UserStorage,
			db.TokenStorage,
			db.OrderStorage,
			cfg,
			log,
		),
		BalanceHistoryService: NewBalanceHistoryService(db.BalanceHistoryStorage),
		RefreshTokenService:   NewRefreshTokenService(db.TokenStorage),
		StatusService:         NewStatusService(db.StatusStorage),
		OrderService:          NewOrderService(db.OrderStorage, log),
	}
}
