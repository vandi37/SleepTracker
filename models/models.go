package models

import (
	"regexp"
	"time"
)

type (
	User struct {
		ID int64 `json:"id"`
		// 3-40 chars
		Username     string    `json:"username"`
		Nickname     string    `json:"nickname"`
		PasswordHash []byte    `json:"-"`
		Birth        time.Time `json:"birth"`
		CreatedAt    time.Time `json:"created_at,omitempty"`
	}
	Friend struct {
		ID         int64     `json:"friendship_id"`
		User1      User      `json:"user1"`
		User2      User      `json:"user2"`
		IsAccepted bool      `json:"is_accepted"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}
	Sleep struct {
		ID     int64 `json:"id"`
		UserID int64 `json:"user_id"`
		// 0-2160
		SleepTime NullInt16 `json:"sleep_time"`
		// 0-2160
		WakeTime NullInt16 `json:"wake_time"`
		// 0-100
		Score     int8      `json:"score"`
		EnterDate time.Time `json:"enter_date"`
		CreatedAt time.Time `json:"created_at"`
	}
	SleepScore struct {
		// 0-100
		Score     int8      `json:"score"`
		EnterDate time.Time `json:"enter_date"`
	}
)

func ValidUsername(u string) bool {
	if len(u) < 3 || len(u) > 40 {
		return false
	}
	matched, err := regexp.MatchString("^[a-zA-Z0-9_]+$", u)
	return err == nil && matched
}
func ValidSleepWake(sleep, wake NullInt16) bool {
	return (sleep.Valid && wake.Valid && sleep.Int16 < wake.Int16) || (!sleep.Valid && !wake.Valid)
}
func ValidTime(time int16) bool {
	return time >= 0 && time < 2160
}

func ValidPercent(percent int16) bool {
	return percent >= 0 && percent <= 100
}
