package apperrors

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrDBQuery           = errors.New("database query error")
	ErrHashPassword      = errors.New("error hashing password")
	ErrJWTGeneration     = errors.New("error generating JWT")
	ErrSaveToken         = errors.New("error saving token")
)
