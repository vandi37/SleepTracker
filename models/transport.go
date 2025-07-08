package models

import "time"

type (
	UserWithToken struct {
		Id      int64     `json:"id,omitempty"`
		Access  string    `json:"access,omitempty"`
		Expires time.Time `json:"expires,omitempty"`
		Refresh string    `json:"refresh,omitempty"`
	}
	UpdateUser struct {
		Id       int64     `json:"id,omitempty"`
		Username string    `json:"username,omitempty"`
		Nickname string    `json:"nickname,omitempty"`
		Birth    time.Time `json:"birth,omitempty"`
	}
	GetFriendships struct {
		UserId int64 `json:"user_id"`
		Limit  int   `json:"limit"`
		Offset int   `json:"offset"`
	}
)
