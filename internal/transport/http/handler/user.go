package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Zrossiz/gophermart/internal/apperrors"
	"github.com/Zrossiz/gophermart/internal/dto"
)

type UserHandler struct {
	service UserService
}

type UserService interface {
	Registration(registrationDTO dto.Registration) (string, string, error)
	Login(loginDTO dto.Registration) (string, string, error)
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (u *UserHandler) Registration(rw http.ResponseWriter, r *http.Request) {
	var registrationDTO dto.Registration

	err := json.NewDecoder(r.Body).Decode(&registrationDTO)
	if err != nil {
		http.Error(rw, "invalid request body", http.StatusBadRequest)
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

	accessToken, refreshToken, err := u.service.Registration(registrationDTO)
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
		http.Error(rw, "invalid request body", http.StatusBadRequest)
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

	accessToken, refreshToken, err := u.service.Login(loginDTO)
	if err != nil {
		switch err {
		case apperrors.ErrInvalidPassword:
			http.Error(rw, "unauthorized", http.StatusUnauthorized)
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
		"message": "login successful",
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(response)
}
