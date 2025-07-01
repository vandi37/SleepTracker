package service

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vandi37/SleepTracker/models"
)

func handleJWTError(err error) models.Error {
	switch {
	case errors.Is(err, jwt.ErrTokenMalformed):
		return models.JsonError{Status: http.StatusBadRequest, Message: "invalid token format"}
	case errors.Is(err, jwt.ErrTokenExpired):
		return models.JsonError{Status: http.StatusUnauthorized, Message: "token expired"}
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return models.JsonError{Status: http.StatusUnauthorized, Message: "token not valid yet"}
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return models.JsonError{Status: http.StatusBadRequest, Message: "invalid token signature"}
	default:
		return models.JsonError{Status: http.StatusUnauthorized, Message: "invalid token"}
	}
}
