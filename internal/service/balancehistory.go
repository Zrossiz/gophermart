package service

import (
	"github.com/Zrossiz/gophermart/internal/dto"
	"github.com/Zrossiz/gophermart/internal/model"
)

type BalanceHistoryService struct {
	db BalanceHistoryStorage
}

type BalanceHistoryStorage interface {
	Create(balanceHistoryDTO dto.CreateBalanceHistory) (bool, error)
	GetAllDebits(userID int64) ([]model.BalanceHistory, error)
	Withdraw(userId, orderId, sum int) error
}

func NewBalanceHistoryService(balanceHistoryStorage BalanceHistoryStorage) *BalanceHistoryService {
	return &BalanceHistoryService{
		db: balanceHistoryStorage,
	}
}

func (b *BalanceHistoryService) Withdraw(userId, orderId, sum int) error {
	err := b.db.Withdraw(userId, orderId, sum)
	if err != nil {
		return err
	}

	return nil
}
