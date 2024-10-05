package service

import (
	"github.com/Zrossiz/gophermart/internal/model"
	"github.com/shopspring/decimal"
)

type UserService struct {
	db *UserStorage
}

type UserStorage interface {
	Create(name string, password string) (bool, error)
	GetUserByName(name string) (*model.User, error)
	UpdateUserBalance(userID int64, balance decimal.Decimal) (bool, error)
}

func NewUserService(userStorage UserStorage) *UserService {
	return &UserService{
		db: &userStorage,
	}
}

func (u *UserService) Registration() {

}

func (u *UserService) Login() {

}
