package service

import (
	"github.com/Zrossiz/gophermart/internal/dto"
	"github.com/Zrossiz/gophermart/internal/model"
)

type BalanceHistoryService struct {
	db *BalanceHistoryStorage
}

type BalanceHistoryStorage interface {
	Create(balanceHistoryDTO dto.CreateBalanceHistory) (bool, error)
	GetAllDebits(userID int64) ([]model.BalanceHistory, error)
}

func NewBalanceHistoryService(balanceHistoryStorage BalanceHistoryStorage) *BalanceHistoryService {
	return &BalanceHistoryService{
		db: &balanceHistoryStorage,
	}
}
