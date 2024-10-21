package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Zrossiz/gophermart/internal/apperrors"
	"github.com/Zrossiz/gophermart/internal/dto"
	"github.com/Zrossiz/gophermart/internal/middleware"
	"github.com/Zrossiz/gophermart/internal/utils"
)

type UserHandler struct {
	userService           UserService
	orderService          OrderService
	balanceHistoryService BalanceHistoryService
}

type UserService interface {
	Registration(registrationDTO dto.Registration) (string, string, error)
	Login(loginDTO dto.Registration) (string, string, error)
	GetUserBalance(username string) (float64, float64, error)
}

func NewUserHandler(
	userService UserService,
	orderSerice OrderService,
	balanceHistoryService BalanceHistoryService,
) *UserHandler {
	return &UserHandler{
		userService:           userService,
		orderService:          orderSerice,
		balanceHistoryService: balanceHistoryService,
	}
}

func (u *UserHandler) Registration(rw http.ResponseWriter, r *http.Request) {
	var registrationDTO dto.Registration

	err := json.NewDecoder(r.Body).Decode(&registrationDTO)
	if err != nil {
		http.Error(rw, "invalID request body", http.StatusBadRequest)
		return
	}

	if registrationDTO.Login == "" {
		http.Error(rw, "login can not be empty", http.StatusBadRequest)
		return
	}

	if registrationDTO.Password == "" {
		http.Error(rw, "password can not be empty", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := u.userService.Registration(registrationDTO)
	if err != nil {
		switch err {
		case apperrors.ErrUserAlreadyExists:
			http.Error(rw, err.Error(), http.StatusConflict)
		case apperrors.ErrDBQuery:
			http.Error(rw, "internal server error", http.StatusInternalServerError)
		case apperrors.ErrHashPassword, apperrors.ErrJWTGeneration:
			http.Error(rw, "error processing request", http.StatusInternalServerError)
		default:
			http.Error(rw, "unknown error", http.StatusInternalServerError)
		}
		return
	}

	refreshTokenCokie := http.Cookie{
		Name:     "refreshtoken",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(2 * time.Hour * 24 * 30),
		HttpOnly: true,
		Secure:   false,
	}

	accessTokenCookie := http.Cookie{
		Name:     "accesstoken",
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   false,
	}

	http.SetCookie(rw, &refreshTokenCokie)
	http.SetCookie(rw, &accessTokenCookie)
	response := map[string]string{
		"message": "registration successful",
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(response)
}

func (u *UserHandler) Login(rw http.ResponseWriter, r *http.Request) {
	var loginDTO dto.Registration

	err := json.NewDecoder(r.Body).Decode(&loginDTO)
	if err != nil {
		http.Error(rw, "invalID request body", http.StatusBadRequest)
		return
	}

	if loginDTO.Login == "" {
		http.Error(rw, "login can not be empty", http.StatusBadRequest)
		return
	}

	if loginDTO.Password == "" {
		http.Error(rw, "password can not be empty", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := u.userService.Login(loginDTO)
	if err != nil {
		switch err {
		case apperrors.ErrInvalIDPassword:
			http.Error(rw, "unauthorized", http.StatusUnauthorized)
		case apperrors.ErrUserAlreadyExists:
			http.Error(rw, "user not found", http.StatusBadRequest)
		case apperrors.ErrDBQuery:
			http.Error(rw, "internal server error", http.StatusInternalServerError)
		case apperrors.ErrHashPassword, apperrors.ErrJWTGeneration:
			http.Error(rw, "error processing request", http.StatusInternalServerError)
		default:
			fmt.Println(err)
			http.Error(rw, "unknown error", http.StatusInternalServerError)
		}
		return
	}

	refreshTokenCookie := http.Cookie{
		Name:     "refreshtoken",
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(2 * time.Hour * 24 * 30),
		HttpOnly: true,
		Secure:   false,
	}

	accessTokenCookie := http.Cookie{
		Name:     "accesstoken",
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Secure:   false,
	}

	http.SetCookie(rw, &refreshTokenCookie)
	http.SetCookie(rw, &accessTokenCookie)
	response := map[string]string{
		"message": "login successful",
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		http.Error(rw, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (u *UserHandler) Withdraw(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(int)
	if !ok {
		http.Error(rw, "could not get user ID", http.StatusUnauthorized)
		return
	}

	var withdrawDTO dto.Withdraw

	err := json.NewDecoder(r.Body).Decode(&withdrawDTO)
	if err != nil {
		http.Error(rw, "invalid request body", http.StatusBadRequest)
		return
	}

	err = u.balanceHistoryService.Withdraw(userID, withdrawDTO.Order, withdrawDTO.Sum)
	if err != nil {
		switch err {
		case apperrors.ErrNotEnoughMoney:
			http.Error(rw, "not enough money on account", http.StatusPaymentRequired)
		case apperrors.ErrInvalIDOrderID:
			http.Error(rw, "invalid order ID", http.StatusUnprocessableEntity)
		default:
			fmt.Println(err)
			http.Error(rw, "unknown error", http.StatusInternalServerError)
		}
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (u *UserHandler) UploadOrder(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(int)
	if !ok {
		http.Error(rw, "could not get user ID", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = u.orderService.UploadOrder(string(body), userID)
	if err != nil {
		switch err {
		case apperrors.ErrOrderAlreadyUploaded:
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
		case apperrors.ErrOrderAlreadyUploadedByAnotherUser:
			http.Error(rw, "order already uploaded by another user", http.StatusConflict)
		case apperrors.ErrDBQuery:
			http.Error(rw, "internal server error", http.StatusInternalServerError)
		case apperrors.ErrInvalIDOrderID:
			http.Error(rw, "invalID order ID", http.StatusUnprocessableEntity)
		default:
			fmt.Println(err)
			http.Error(rw, "unknown error", http.StatusInternalServerError)
		}
		return
	}

	rw.WriteHeader(http.StatusAccepted)
}

func (u *UserHandler) GetAllOrdersByUser(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(int)
	if !ok {
		http.Error(rw, "could not get user ID", http.StatusUnauthorized)
		return
	}

	orders, err := u.orderService.GetAllOrdersByUser(userID)
	if err != nil {
		switch err {
		case apperrors.ErrOrdersNotFound:
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte("[]"))
		default:
			http.Error(rw, "failed to get orders", http.StatusInternalServerError)
		}
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	err = json.NewEncoder(rw).Encode(orders)
	if err != nil {
		http.Error(rw, "unable to encode response", http.StatusInternalServerError)
		return
	}
}

func (u *UserHandler) GetUserBalance(rw http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value(middleware.UserNameContextKey).(string)
	if !ok {
		http.Error(rw, "could not get user ID", http.StatusUnauthorized)
		return
	}

	current, withdrawn, err := u.userService.GetUserBalance(username)
	if err != nil {
		http.Error(rw, "unknown error", http.StatusInternalServerError)
		return
	}

	response := dto.BalanceUser{
		Current:   utils.Round(current, 5),
		Withdrawn: utils.Round(withdrawn, 5),
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		http.Error(rw, "unable to encode response", http.StatusInternalServerError)
	}
}

func (u *UserHandler) Withdrawals(rw http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDContextKey).(int)
	if !ok {
		http.Error(rw, "could not get user ID", http.StatusUnauthorized)
		return
	}

	withdrwals, err := u.balanceHistoryService.GetAllWithdrawlsByUser(userID)
	if err != nil {
		switch err {
		case apperrors.ErrWithdrawlsNotFound:
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			return
		default:
			fmt.Println(err)
			http.Error(rw, "failed to get withdrawls", http.StatusInternalServerError)
			return
		}
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	err = json.NewEncoder(rw).Encode(withdrwals)
	if err != nil {
		http.Error(rw, "unable to encode response", http.StatusInternalServerError)
	}
}
