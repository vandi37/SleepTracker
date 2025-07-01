package friend_repo

import (
	"fmt"
	"net/http"

	"github.com/vandi37/SleepTracker/models"
)

type FriendshipAlreadyExists struct {
	from, to int64
}

func (f FriendshipAlreadyExists) Error() string {
	return fmt.Sprintf("friendship or friendship request between %x and %x already exists", f.from, f.to)
}
func (FriendshipAlreadyExists) Code() int { return http.StatusConflict }
func (f FriendshipAlreadyExists) JsonError() models.JsonError {
	return models.JsonError{
		Status:    http.StatusConflict,
		Message: "friendship or friendship request already exists",
		Context: map[string]any{
			"user1_id": f.from,
			"user2_id": f.to,
		},
	}
}

type SelfRequest int64

func (s SelfRequest) Error() string {
	return fmt.Sprintf("can't requests yourself for a friendship (%x)", int64(s))
}
func (SelfRequest) Code() int { return http.StatusUnprocessableEntity }
func (s SelfRequest) JsonError() models.JsonError {
	return models.JsonError{
		Status:    http.StatusUnprocessableEntity,
		Message: "can't requests yourself for a friendship",
		Context: map[string]any{"user_id": s},
	}
}

type FriendshipNotFound int64

func (f FriendshipNotFound) Error() string { return fmt.Sprintf("friendship %x not found", int64(f)) }
func (FriendshipNotFound) Code() int       { return http.StatusNotFound }
func (u FriendshipNotFound) JsonError() models.JsonError {
	return models.JsonError{
		Status:    http.StatusNotFound,
		Message: "friendship not found",
		Context: map[string]any{"id": u},
	}
}

type FriendshipAlreadyAccepted int64

func (f FriendshipAlreadyAccepted) Error() string {
	return fmt.Sprintf("friendship %x already accepted", int64(f))
}
func (FriendshipAlreadyAccepted) Code() int { return http.StatusConflict }
func (f FriendshipAlreadyAccepted) JsonError() models.JsonError {
	return models.JsonError{
		Status:    http.StatusConflict,
		Message: "friendship already accepted",
		Context: map[string]any{"id": f},
	}
}

