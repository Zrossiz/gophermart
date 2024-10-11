package apperrors

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrDBQuery           = errors.New("database query error")
	ErrHashPassword      = errors.New("error hashing password")
	ErrJWTGeneration     = errors.New("error generating JWT")
	ErrSaveToken         = errors.New("error saving token")
)
