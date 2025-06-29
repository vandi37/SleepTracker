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
		CreatedAt    time.Time `json:"created_at"`
	}
	Friend struct {
		FriendshipID int64            `json:"friendship_id"`
		User1ID      int64            `json:"user1_id"`
		User2ID      int64            `json:"user2_id"`
		Status       FriendshipStatus `json:"status"`
		CreatedAt    time.Time        `json:"created_at"`
	}
	Sleep struct {
		ID     int64 `json:"id"`
		UserID int64 `json:"user_id"`
		// 0-1440
		SleepTime NullInt16 `json:"sleep_time"`
		// 0-1440
		WakeTime NullInt16 `json:"wake_time"`
		// 0-100
		Score     int8      `json:"score"`
		EnterDate time.Time `json:"enter_date"`
		CreatedAt time.Time `json:"created_at"`
	}
)

func ValidUsername(u string) bool {
	if len(u) < 3 || len(u) > 40 {
		return false
	}
	matched, err := regexp.MatchString("^[a-zA-Z0-9_]+$", u)
	return err == nil && matched
}

func ValidStatus(s FriendshipStatus) bool {
	switch s {
	case FriendshipRequested, FriendshipAccepted, FriendshipRejected:
		return true
	}
	return false
}
func ValidSleepWake(sleep, wake NullInt16) bool {
	if !sleep.Valid && !wake.Valid {
		return true
	}
	return false
}
func ValidTime(time int64) bool {
	return time >= 0 && time < 1440
}

func ValidPercent(percent int16) bool {
	return percent >= 0 && percent <= 100
}
