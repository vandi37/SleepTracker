package repo

import (
	"errors"
	"fmt"
)

var NoRows = errors.New("no rows")

func ErrScanning(err error) error {
	return fmt.Errorf("error scanning: %w", err)
}
