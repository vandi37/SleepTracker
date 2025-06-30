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
		Code:    http.StatusConflict,
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
		Code:    http.StatusUnprocessableEntity,
		Message: "can't requests yourself for a friendship",
		Context: map[string]any{"user_id": s},
	}
}

type FriendshipNotFound int64

func (f FriendshipNotFound) Error() string { return fmt.Sprintf("friendship %x not found", int64(f)) }
func (FriendshipNotFound) Code() int       { return http.StatusNotFound }
func (u FriendshipNotFound) JsonError() models.JsonError {
	return models.JsonError{
		Code:    http.StatusNotFound,
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
		Code:    http.StatusConflict,
		Message: "friendship already accepted",
		Context: map[string]any{"id": f},
	}
}

type NotAllowed struct{ by, id int64 }

func (f NotAllowed) Error() string {
	return fmt.Sprintf("user %x is not allowed to accept friendship request %x", f.by, f.id)
}
func (NotAllowed) Code() int { return http.StatusForbidden }
func (f NotAllowed) JsonError() models.JsonError {
	return models.JsonError{
		Code:    http.StatusForbidden,
		Message: "user is not allowed to accept friendship request",
		Context: map[string]any{"id": f.id, "user2_id": f.by},
	}
}
