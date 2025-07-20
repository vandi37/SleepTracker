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
	Password struct {
		Password string `json:"password"`
	}
	Id struct {
		Id int64 `json:"id"`
	}
	Friendships struct {
		Friendships []Friend `json:"friendships"`
	}
	Sleeps struct {
		Sleeps []Sleep `json:"sleeps"`
	}
	Scores struct {
		Scores []SleepScore `json:"scores"`
	}
)
