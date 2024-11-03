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
	Withdraw(userID int, orderID string, sum float64) error
}

func NewBalanceHistoryService(balanceHistoryStorage BalanceHistoryStorage) *BalanceHistoryService {
	return &BalanceHistoryService{
		db: balanceHistoryStorage,
	}
}

func (b *BalanceHistoryService) Withdraw(userID int, orderID string, sum float64) error {
	err := b.db.Withdraw(userID, orderID, sum)
	if err != nil {
		return err
	}

	return nil
}

func (b *BalanceHistoryService) GetAllWithdrawlsByUser(userID int) ([]model.BalanceHistory, error) {
	withdrawls, err := b.db.GetAllDebits(int64(userID))
	if err != nil {
		return nil, err
	}

	return withdrawls, nil
}
