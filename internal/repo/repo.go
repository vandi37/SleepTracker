package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/vandi37/SleepTracker/models"
	"go.uber.org/zap"
)

type User interface {
	Create(ctx context.Context, tx *sql.Tx, username, nickname string, password []byte, birth time.Time) (int64, models.Error)
	Update(ctx context.Context, tx *sql.Tx, id int64, username, nickname string, birth time.Time) models.Error
	UpdatePassword(ctx context.Context, tx *sql.Tx, id int64, password []byte) models.Error
	Get(ctx context.Context, tx *sql.Tx, id int64) (models.User, models.Error)
	GetWithPassword(ctx context.Context, tx *sql.Tx, username string, password []byte) (models.User, models.Error)
	Delete(ctx context.Context, tx *sql.Tx, id int64) models.Error
}

var RepoNamespace = zap.Namespace("repository")
