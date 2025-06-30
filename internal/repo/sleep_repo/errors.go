package sleep_repo

import (
	"fmt"
	"net/http"

	"github.com/vandi37/SleepTracker/models"
)

type SleepRecordNotFound int64

func (s SleepRecordNotFound) Error() string { return fmt.Sprintf("friendship %x not found", int64(s)) }
func (SleepRecordNotFound) Code() int       { return http.StatusNotFound }
func (s SleepRecordNotFound) JsonError() models.JsonError {
	return models.JsonError{
		Code:    http.StatusNotFound,
		Message: "friendship not found",
		Context: map[string]any{"id": s},
	}
}

type SleepRecordAlreadyWritten int64

func (s SleepRecordAlreadyWritten) Error() string {
	return fmt.Sprintf("today's sleep record by user %x already exists", int64(s))
}
func (SleepRecordAlreadyWritten) Code() int { return http.StatusConflict }
func (s SleepRecordAlreadyWritten) JsonError() models.JsonError {
	return models.JsonError{
		Code:    http.StatusConflict,
		Message: "sleep record already exists",
		Context: map[string]any{
			"user_id": s,
		},
	}
}

type InvalidSleepWake struct{}

func (InvalidSleepWake) Error() string {
	return "sleep and wake should be either both null or both not null"
}
func (InvalidSleepWake) Code() int { return http.StatusUnprocessableEntity }
func (i InvalidSleepWake) JsonError() models.JsonError {
	return models.JsonError{
		Code:    http.StatusUnprocessableEntity,
		Message: i.Error(),
	}
}
