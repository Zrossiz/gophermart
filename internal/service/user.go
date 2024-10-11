package service

import (
	"time"

	"github.com/Zrossiz/gophermart/internal/apperrors"
	"github.com/Zrossiz/gophermart/internal/config"
	"github.com/Zrossiz/gophermart/internal/dto"
	"github.com/Zrossiz/gophermart/internal/model"
	"github.com/Zrossiz/gophermart/internal/utils"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	dbUser  UserStorage
	dbToken TokenStorage
	cfg     config.Config
	log     zap.Logger
}

type UserStorage interface {
	Create(name string, password string) (bool, error)
	GetUserByName(name string) (*model.User, error)
	UpdateUserBalance(userID int64, balance decimal.Decimal) (bool, error)
}

func NewUserService(userStorage UserStorage, tokenStorage TokenStorage, cfg *config.Config, log *zap.Logger) *UserService {
	return &UserService{
		dbUser:  userStorage,
		dbToken: tokenStorage,
		cfg:     *cfg,
		log:     *log,
	}
}

func (u *UserService) Registration(registrationDTO dto.Registration) (string, string, error) {
	existUser, err := u.dbUser.GetUserByName(registrationDTO.Login)
	if err != nil {
		u.log.Sugar().Error("db query error: %v", err)
		return "", "", apperrors.ErrDBQuery
	}
	if existUser != nil {
		return "", "", apperrors.ErrUserAlreadyExists
	}

	hashedPassword, err := hashPassword(registrationDTO.Password, u.cfg.Cost)
	if err != nil {
		return "", "", apperrors.ErrHashPassword
	}

	_, err = u.dbUser.Create(registrationDTO.Login, hashedPassword)
	if err != nil {
		u.log.Error(err.Error())
		return "", "", apperrors.ErrDBQuery
	}

	createdUser, err := u.dbUser.GetUserByName(registrationDTO.Login)
	if err != nil {
		u.log.Error(err.Error())
		return "", "", apperrors.ErrDBQuery
	}

	JWTAccessProps := utils.GenerateJWTProps{
		Secret:   []byte(u.cfg.AccessTokenSecret),
		Exprires: time.Now().Add(15 * time.Minute),
		UserID:   int64(createdUser.ID),
		Username: createdUser.Name,
	}

	accessToken, err := utils.GenerateJWT(JWTAccessProps)
	if err != nil {
		u.log.Error(err.Error())
		return "", "", apperrors.ErrJWTGeneration
	}

	JWTRefreshProps := utils.GenerateJWTProps{
		Secret:   []byte(u.cfg.RefreshTokenSecret),
		Exprires: time.Now().Add(24 * 30 * time.Hour),
		UserID:   int64(createdUser.ID),
		Username: createdUser.Name,
	}

	refreshToken, err := utils.GenerateJWT(JWTRefreshProps)
	if err != nil {
		u.log.Error(err.Error())
		return "", "", apperrors.ErrJWTGeneration
	}

	_, err = u.dbToken.Create(int64(createdUser.ID), refreshToken)
	if err != nil {
		u.log.Error(err.Error())
		return "", "", apperrors.ErrDBQuery
	}

	return accessToken, refreshToken, nil
}

func (u *UserService) Login() {

}

func hashPassword(password string, cost int) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
