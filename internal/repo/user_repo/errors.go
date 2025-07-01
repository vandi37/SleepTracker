package user_repo

import (
	"fmt"
	"net/http"

	"github.com/vandi37/SleepTracker/models"
)

type UsernameTakenError string

func (u UsernameTakenError) Error() string {
	return fmt.Sprintf("username %s is taken", string(u))
}
func (UsernameTakenError) Code() int { return http.StatusConflict }
func (u UsernameTakenError) JsonError() models.JsonError {
	return models.JsonError{
		Status:  http.StatusConflict,
		Message: "username taken",
		Context: map[string]any{
			"username": u,
		},
	}
}

type UserNotFound int64

func (u UserNotFound) Error() string { return fmt.Sprintf("user %x not found", int64(u)) }
func (UserNotFound) Code() int       { return http.StatusNotFound }
func (u UserNotFound) JsonError() models.JsonError {
	return models.JsonError{
		Status:  http.StatusNotFound,
		Message: "user not found",
		Context: map[string]any{"id": u},
	}
}

type InvalidCredentials struct{}

func (InvalidCredentials) Error() string { return "invalid credentials" }
func (InvalidCredentials) Code() int     { return http.StatusUnauthorized }
func (i InvalidCredentials) JsonError() models.JsonError {
	return models.JsonError{
		Status:  http.StatusUnauthorized,
		Message: i.Error(),
	}
}
