package user_repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/vandi37/SleepTracker/internal/repo"
	"github.com/vandi37/SleepTracker/models"
	"github.com/vandi37/SleepTracker/pkg/logger"
	"go.uber.org/zap"
)

type Repo struct{}

// GetBirth implements repo.User.
func (Repo) GetBirth(ctx context.Context, tx *sql.Tx, id int64) (time.Time, models.Error) {
	var time time.Time
	err := tx.QueryRowContext(ctx, `select birth from users where id = $1`, id).
		Scan(&time)
	if err == sql.ErrNoRows {
		return time, UserNotFound(id)
	}
	if err != nil {
		logger.Error(ctx, "got an internal error while getting user birth", repo.Namespace, zap.Error(err), zap.Int64("id", id))
		return time, models.Internal(err)
	}
	return time, nil
}

// Create implements repo.User.
func (Repo) Create(ctx context.Context, tx *sql.Tx, username string, nickname string, password []byte, birth time.Time) (int64, models.Error) {
	if !models.ValidUsername(username) {
		return 0, models.Invalid("username", username)
	}
	if nickname == "" {
		return 0, models.Invalid("nickname", nickname)
	}
	var id int64
	err := tx.QueryRowContext(
		ctx,
		`insert into users (username, nickname, password_hash, birth) values ($1, $2, $3, $4) returning id`,
		username,
		nickname,
		password,
		birth,
	).Scan(&id)
	if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
		return 0, UsernameTakenError(username)
	} else if err != nil {
		logger.Error(ctx, "got an internal error while inserting a new user",
			repo.Namespace,
			zap.Error(err),
			zap.String("username", username),
			zap.String("nickname", nickname),
		)
		return 0, models.Internal(err)
	}
	return id, nil
}

// Delete implements repo.User.
func (Repo) Delete(ctx context.Context, tx *sql.Tx, id int64) models.Error {
	res, err := tx.ExecContext(ctx, `delete from users where id = $1`, id)
	if err != nil {
		logger.Error(ctx, "got an internal error while deleting user", repo.Namespace, zap.Error(err), zap.Int64("id", id))
		return models.Internal(err)
	}

	if ok, err := repo.CheckRes(res, repo.Equals(1)); err != nil {
		logger.Error(ctx, "got an internal error while checking result of deleting user", repo.Namespace, zap.Error(err), zap.Int64("id", id))
		return models.Internal(err)
	} else if !ok {
		return UserNotFound(id)
	}
	return nil
}

// Get implements repo.User.
func (Repo) Get(ctx context.Context, tx *sql.Tx, id int64) (models.User, models.Error) {
	var user models.User
	err := tx.QueryRowContext(ctx, `select id, username, nickname, birth, created_at from users where id = $1`, id).
		Scan(&user.Id, &user.Username, &user.Nickname, &user.Birth, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return user, UserNotFound(id)
	}
	if err != nil {
		logger.Error(ctx, "got an internal error while getting user", repo.Namespace, zap.Error(err), zap.Int64("id", id))
		return user, models.Internal(err)
	}
	return user, nil
}

// GetWithPassword implements repo.User.
func (Repo) GetByUsername(ctx context.Context, tx *sql.Tx, username string) (int64, []byte, models.Error) {
	var id int64
	var passwordHash []byte
	err := tx.QueryRowContext(ctx, `select id, password_hash, from users where username = $1`, username).
		Scan(&id, &passwordHash)
	if err == sql.ErrNoRows {
		return id, passwordHash, InvalidCredentials{}
	}
	if err != nil {
		logger.Error(ctx, "got an internal error while getting user", repo.Namespace, zap.Error(err), zap.String("username", username))
		return id, passwordHash, models.Internal(err)
	}
	return id, passwordHash, nil
}

// Update implements repo.User.
func (Repo) Update(ctx context.Context, tx *sql.Tx, id int64, username string, nickname string, birth time.Time) models.Error {
	if !models.ValidUsername(username) {
		return models.Invalid("username", username)
	}
	if nickname == "" {
		return models.Invalid("nickname", nickname)
	}
	res, err := tx.ExecContext(ctx, `update users set username = $2, nickname = $3, birth = $4 where id = $1`, id, username, nickname, birth)
	if err, ok := err.(*pq.Error); ok && err.Code == "23505" {
		return UsernameTakenError(username)
	} else if err != nil {
		logger.Error(ctx, "got an internal error while updating user",
			repo.Namespace,
			zap.Error(err),
			zap.Int64("id", id),
			zap.String("username", username),
			zap.String("nickname", nickname),
		)
		return models.Internal(err)
	}
	if ok, err := repo.CheckRes(res, repo.Equals(1)); err != nil {
		logger.Error(ctx, "got an internal error while checking result of updating user",
			repo.Namespace,
			zap.Error(err),
			zap.Int64("id", id),
			zap.String("username", username),
			zap.String("nickname", nickname),
		)
		return models.Internal(err)
	} else if !ok {
		return UserNotFound(id)
	}
	return nil
}

// UpdatePassword implements repo.User.
func (Repo) UpdatePassword(ctx context.Context, tx *sql.Tx, id int64, password []byte) models.Error {
	res, err := tx.ExecContext(ctx, `update users set password_hash = $2 where id = $1`, id, password)
	if err != nil {
		logger.Error(ctx, "got an internal error while updating user password", repo.Namespace, zap.Error(err), zap.Int64("id", id))
		return models.Internal(err)
	}
	if ok, err := repo.CheckRes(res, repo.Equals(1)); err != nil {
		logger.Error(ctx, "got an internal error while checking result of updating user", repo.Namespace, zap.Error(err), zap.Int64("id", id))
		return models.Internal(err)
	} else if !ok {
		return UserNotFound(id)
	}
	return nil
}

var _ repo.User = (*Repo)(nil)
