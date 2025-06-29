package repo

import (
	"database/sql"
)

func CheckRes(res sql.Result, check func(int64) bool) (bool, error) {
	n, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return check(n), nil
}

func Equals(n int64) func(int64) bool {
	return func(i int64) bool {
		return i == n
	}
}

func More(n int64) func(int64) bool {
	return func(i int64) bool {
		return i > n
	}
}
