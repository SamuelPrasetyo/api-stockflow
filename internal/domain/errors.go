package domain

import "errors"

var (
	// User errors
	ErrUserNotFound        = errors.New("user not found")
	ErrUsernameAlreadyUsed = errors.New("username already used")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserIDInvalid       = errors.New("user id is invalid")

	// Authentication errors
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrInvalidToken         = errors.New("invalid token")
	ErrTokenExpired         = errors.New("token expired")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrForbidden            = errors.New("forbidden - insufficient permissions")

	// General errors
	ErrInternalServer = errors.New("internal server error")
	ErrBadRequest     = errors.New("bad request")
)
