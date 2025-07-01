package models

import "time"

type (
	UserWithToken struct {
		Id      int64     `json:"id,omitempty"`
		Access  string    `json:"access,omitempty"`
		Expires time.Time `json:"expires,omitempty"`
		Refresh string    `json:"refresh,omitempty"`
	}
)
