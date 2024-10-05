package service

import (
	"github.com/Zrossiz/gophermart/internal/model"
)

type StatusService struct {
	db *StatusStorage
}

type StatusStorage interface {
	Create(status string) (bool, error)
	GetAll() ([]model.Status, error)
}

func NewStatusService(statusStorage StatusStorage) *StatusService {
	return &StatusService{
		db: &statusStorage,
	}
}
