package models

import (
	"time"

	"github.com/vandi37/SleepTracker/pkg/date"
)

type (
	UserReq struct {
		Username string    `json:"username"`
		Password string    `json:"password"`
		Nickname string    `json:"nickname"`
		Birth    date.Date `json:"birth"`
	}
	UserWithToken struct {
		Id      int64     `json:"id"`
		Access  string    `json:"access,omitempty"`
		Expires time.Time `json:"expires,omitempty"`
		Refresh string    `json:"refresh,omitempty"`
	}
	GetFriendships struct {
		UserId int64 `json:"user_id"`
		Limit  int   `json:"limit"`
		Offset int   `json:"offset"`
	}
)
