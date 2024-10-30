package service

import (
	"errors"
	"testing"

	"github.com/Zrossiz/gophermart/internal/apperrors"
	"github.com/Zrossiz/gophermart/internal/config"
	"github.com/Zrossiz/gophermart/internal/dto"
	"github.com/Zrossiz/gophermart/internal/model"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type MockUserStorage struct {
	mock.Mock
}

func (m *MockUserStorage) Create(name string, password string) (bool, error) {
	args := m.Called(name, password)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserStorage) GetUserByName(name string) (*model.User, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), nil
}

func (m *MockUserStorage) UpdateUserBalance(userID int64, balance decimal.Decimal) (bool, error) {
	args := m.Called(userID, balance)
	return args.Bool(0), args.Error(1)
}

type MockTokenStorage struct {
	mock.Mock
}

func (m *MockTokenStorage) Create(userID int64, token string) (bool, error) {
	args := m.Called(userID, token)
	return args.Bool(0), args.Error(1)
}

func (m *MockTokenStorage) DeleteByToken(token string) (bool, error) {
	args := m.Called(token)
	return args.Bool(0), args.Error(1)
}

func (m *MockTokenStorage) DeleteTokenByUser(userID int64) (bool, error) {
	args := m.Called(userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockTokenStorage) GetTokenByToken(token string) (*model.RefreshToken, error) {
	args := m.Called(token)
	if args.Get(0) != nil {
		return args.Get(0).(*model.RefreshToken), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockTokenStorage) GetTokenByUser(userID int64) (*model.RefreshToken, error) {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).(*model.RefreshToken), args.Error(1)
	}
	return nil, args.Error(1)
}

type MockOrderStorage struct {
	mock.Mock
}

func (m *MockOrderStorage) CreateOrder(orderID string, userID int) (bool, error) {
	args := m.Called(orderID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockOrderStorage) GetAllOrdersByUser(userID int64) ([]model.Order, error) {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockOrderStorage) UpdateSumAndStatusOrder(orderID string, status string, sum float64) (bool, error) {
	args := m.Called(orderID, status, sum)
	return args.Bool(0), args.Error(1)
}

func (m *MockOrderStorage) GetOrderByID(orderID string) (*model.Order, error) {
	args := m.Called(orderID)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockOrderStorage) GetAllUnhandlerOrders(unhandledStatus1, unhandledStatus2 int) ([]string, error) {
	args := m.Called(unhandledStatus1, unhandledStatus2)
	if args.Get(0) != nil {
		return args.Get(0).([]string), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockOrderStorage) GetAllWithdrawnByUser(userID int64) (float64, error) {
	args := m.Called(userID)
	return args.Get(0).(float64), args.Error(1)
}

func TestUserServiceRegistration(t *testing.T) {
	mockUserStorage := new(MockUserStorage)
	mockOrderStorage := new(MockOrderStorage)
	mockTokenStorage := new(MockTokenStorage)
	cfg := config.Config{
		AccessTokenSecret:  "access",
		RefreshTokenSecret: "refresh",
		Cost:               10,
	}
	log := zap.NewExample()
	service := NewUserService(mockUserStorage, mockTokenStorage, mockOrderStorage, &cfg, log)

	registrationDTO := dto.Registration{
		Login:    "testuser",
		Password: "password123",
	}

	existingUser := &model.User{ID: 1, Name: "testuser"}
	mockUserStorage.On("GetUserByName", registrationDTO.Login).Return(existingUser, nil).Once()

	accessToken, refreshToken, err := service.Registration(registrationDTO)
	if err == nil || err != apperrors.ErrUserAlreadyExists {
		t.Errorf("expected error ErrUserAlreadyExists, got: %v", err)
	}
	if accessToken != "" || refreshToken != "" {
		t.Errorf("tokens must be empty when user already exists")
	}

	mockUserStorage.ExpectedCalls = nil
	mockTokenStorage.ExpectedCalls = nil

	mockUserStorage.On("GetUserByName", registrationDTO.Login).Return(nil, nil).Once()
	mockUserStorage.On("Create", registrationDTO.Login, mock.Anything).Return(true, nil).Once()

	newUser := &model.User{
		ID:       2,
		Name:     "testuser",
		Password: "hashed",
	}

	mockUserStorage.On("GetUserByName", registrationDTO.Login).Return(newUser, nil).Once()
	mockTokenStorage.On("Create", int64(2), mock.Anything).Return(true, nil).Once()

	accessToken, refreshToken, err = service.Registration(registrationDTO)
	if err != nil {
		t.Fatalf("expected err == nil, got: %v", err)
	}
	if accessToken == "" || refreshToken == "" {
		t.Fatalf("expected tokens not to be empty")
	}

	mockUserStorage.AssertExpectations(t)
	mockTokenStorage.AssertExpectations(t)
}

func TestUserServiceLogin(t *testing.T) {
	mockUserStorage := new(MockUserStorage)
	mockOrderStorage := new(MockOrderStorage)
	mockTokenStorage := new(MockTokenStorage)
	cfg := config.Config{
		AccessTokenSecret:  "access_secret",
		RefreshTokenSecret: "refresh_secret",
		Cost:               10,
	}
	log := zap.NewExample()
	service := NewUserService(mockUserStorage, mockTokenStorage, mockOrderStorage, &cfg, log)

	loginDTO := dto.Registration{
		Login:    "testuser",
		Password: "password123",
	}

	mockUserStorage.On("GetUserByName", loginDTO.Login).Return(nil, nil).Once()
	_, _, err := service.Login(loginDTO)
	if err != apperrors.ErrUserNotFound {
		t.Errorf("expected error ErrUserNotFound, got: %v", err)
	}
	mockUserStorage.ExpectedCalls = nil

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("differentpassword"), bcrypt.DefaultCost)
	existingUser := &model.User{ID: 1, Name: "testuser", Password: string(hashedPassword)}
	mockUserStorage.On("GetUserByName", loginDTO.Login).Return(existingUser, nil).Once()
	_, _, err = service.Login(loginDTO)
	if err != apperrors.ErrInvalidPassword {
		t.Errorf("expected error ErrInvalidPassword, got: %v", err)
	}
	mockUserStorage.ExpectedCalls = nil

	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte(loginDTO.Password), bcrypt.DefaultCost)
	existingUser.Password = string(hashedPassword)
	mockUserStorage.On("GetUserByName", loginDTO.Login).Return(existingUser, nil).Once()

	accessToken, refreshToken, err := service.Login(loginDTO)
	if err != nil {
		t.Fatalf("expected err == nil, got: %v", err)
	}
	if accessToken == "" || refreshToken == "" {
		t.Fatalf("expected tokens not to be empty")
	}

	mockUserStorage.AssertExpectations(t)
}

func TestUserServiceGetUserBalance(t *testing.T) {
	mockUserStorage := new(MockUserStorage)
	mockOrderStorage := new(MockOrderStorage)
	mockTokenStorage := new(MockTokenStorage)
	cfg := config.Config{}
	log := zap.NewExample()
	service := NewUserService(mockUserStorage, mockTokenStorage, mockOrderStorage, &cfg, log)

	mockUserStorage.On("GetUserByName", "testuser").Return(nil, nil).Once()
	balance, withdrawn, err := service.GetUserBalance("testuser")
	if err != apperrors.ErrUserNotFound {
		t.Errorf("expected error ErrUserNotFound, got: %v", err)
	}
	if balance != 0.00 || withdrawn != 0.00 {
		t.Errorf("expected balance and withdrawn to be 0.00, got: balance = %v, withdrawn = %v", balance, withdrawn)
	}
	mockUserStorage.ExpectedCalls = nil

	mockUserStorage.On("GetUserByName", "testuser").Return(nil, errors.New("db error")).Once()
	balance, withdrawn, err = service.GetUserBalance("testuser")
	if err != apperrors.ErrDBQuery {
		t.Errorf("expected error ErrDBQuery, got: %v", err)
	}
	if balance != 0.00 || withdrawn != 0.00 {
		t.Errorf("expected balance and withdrawn to be 0.00 on error, got: balance = %v, withdrawn = %v", balance, withdrawn)
	}
	mockUserStorage.ExpectedCalls = nil

	curUser := &model.User{ID: 1, Name: "testuser", Account: 0}
	mockUserStorage.On("GetUserByName", "testuser").Return(curUser, nil).Once()
	mockOrderStorage.On("GetAllWithdrawnByUser", int64(curUser.ID)).Return(0.00, errors.New("db error")).Once()
	balance, withdrawn, err = service.GetUserBalance("testuser")
	if err != nil {
		t.Errorf("expected error = nil, got: %v", err)
	}
	if balance != curUser.Account || withdrawn != 0.00 {
		t.Errorf("expected balance = %v and withdrawn = 0.00 on error, got: balance = %v, withdrawn = %v", curUser.Account, balance, withdrawn)
	}
	mockUserStorage.ExpectedCalls = nil
	mockOrderStorage.ExpectedCalls = nil

	curUser = &model.User{ID: 1, Name: "testuser", Account: 100.50}
	withdrawnAmount := 50.25
	mockUserStorage.On("GetUserByName", "testuser").Return(curUser, nil).Once()
	mockOrderStorage.On("GetAllWithdrawnByUser", int64(curUser.ID)).Return(withdrawnAmount, nil).Once()
	balance, withdrawn, err = service.GetUserBalance("testuser")
	if err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
	if balance != curUser.Account || withdrawn != withdrawnAmount {
		t.Errorf("expected balance = %v and withdrawn = %v, got: balance = %v, withdrawn = %v", curUser.Account, withdrawnAmount, balance, withdrawn)
	}

	mockUserStorage.AssertExpectations(t)
	mockOrderStorage.AssertExpectations(t)
}
