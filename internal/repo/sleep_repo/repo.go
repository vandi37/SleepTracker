package sleep_repo

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

// Delete implements repo.Sleep.
func (Repo) Delete(ctx context.Context, tx *sql.Tx, id int64, user_id int64) models.Error {
	if res, err := tx.ExecContext(ctx, `delete from sleeps where id = $1 and user_id = $2`, id, user_id); err != nil {
		logger.Error(ctx, "got an internal error while deleting sleep record", repo.Namespace, zap.Int64("id", id), zap.Error(err), zap.Int64("user_id", user_id))
		return models.Internal(err)
	} else if ok, err := repo.CheckRes(res, repo.Equals(1)); err != nil {
		logger.Error(ctx, "got an internal error while checking result of deleting a sleep record", repo.Namespace, zap.Error(err), zap.Int64("id", id), zap.Int64("user_id", user_id))
		return models.Internal(err)
	} else if ok {
		return nil
	}
	if err := tx.QueryRowContext(ctx, `select 1 from sleeps where id = $1`, id).Err(); err == sql.ErrNoRows {
		return SleepRecordNotFound(id)
	} else if err != nil {
		logger.Error(ctx, "got an internal error while getting cause of failed delete", repo.Namespace, zap.Error(err), zap.Int64("id", id), zap.Int64("user_id", user_id))
		return models.Internal(err)
	}
	return models.NotAllowed{Id: id, User: user_id, Thing: "delete a sleep record"}
}

// Enter implements repo.Sleep.
func (Repo) Enter(ctx context.Context, tx *sql.Tx, user_id int64, sleep_time, wake_time models.NullInt16, score int16, date time.Time) (int64, models.Error) {
	if !models.ValidSleepWake(sleep_time, wake_time) {
		return 0, InvalidSleepWake{}
	} else if sleep_time.Valid && !models.ValidTime(sleep_time.Int16) {
		return 0, models.Invalid("sleep_time", sleep_time)
	} else if wake_time.Valid && !models.ValidTime(wake_time.Int16) {
		return 0, models.Invalid("wake_time", wake_time)
	} else if !models.ValidPercent(score) {
		return 0, models.Invalid("score", score)
	}
	var id int64
	err := tx.QueryRowContext(ctx, `insert into sleeps (user_id, sleep_time, wake_time, score) values ($1, $2, $3, $4) returning id`,
		user_id, sleep_time, wake_time, score).Scan(&id)
	if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
		return 0, SleepRecordAlreadyWritten(user_id)
	} else if err != nil {
		logger.Error(ctx, "got an internal error while entering a sleep record",
			repo.Namespace, zap.Error(err), zap.Int64("user_id", user_id))
		return 0, models.Internal(err)
	}
	return id, nil
}

// Update implements repo.Sleep.
func (Repo) Update(ctx context.Context, tx *sql.Tx, id, user_id int64, sleep_time, wake_time models.NullInt16, score int16) models.Error {
	if !models.ValidSleepWake(sleep_time, wake_time) {
		return InvalidSleepWake{}
	} else if sleep_time.Valid && !models.ValidTime(sleep_time.Int16) {
		return models.Invalid("sleep_time", sleep_time)
	} else if wake_time.Valid && !models.ValidTime(wake_time.Int16) {
		return models.Invalid("wake_time", wake_time)
	} else if !models.ValidPercent(score) {
		return models.Invalid("score", score)
	}
	if res, err := tx.ExecContext(ctx, `update sleeps set sleep_time = $3, wake_time = $4, score = $5 where id = $1 and user_id = $2`,
		id, user_id, sleep_time, wake_time, score); err != nil {
		logger.Error(ctx, "got an internal error while updating a sleep record",
			repo.Namespace, zap.Error(err), zap.Int64("id", id), zap.Int64("user_id", user_id))
		return models.Internal(err)
	} else if ok, err := repo.CheckRes(res, repo.Equals(1)); err != nil {
		logger.Error(ctx, "got an internal error while checking result of updating a sleep record",
			repo.Namespace, zap.Error(err), zap.Int64("id", id), zap.Int64("user_id", user_id))
		return models.Internal(err)
	} else if ok {
		return nil
	}
	if err := tx.QueryRowContext(ctx, `select 1 from sleeps where id = $1`, id).Err(); err == sql.ErrNoRows {
		return SleepRecordNotFound(id)
	} else if err != nil {
		logger.Error(ctx, "got an internal error while getting cause of failed update", repo.Namespace, zap.Error(err), zap.Int64("id", id), zap.Int64("user_id", user_id))
		return models.Internal(err)
	}
	return models.NotAllowed{Id: id, User: user_id, Thing: "update a sleep record"}
}

// Week implements repo.Sleep.
func (Repo) Week(ctx context.Context, tx *sql.Tx, user_id int64, page int) ([]models.Sleep, models.Error) {
	rows, err := tx.QueryContext(ctx, `select id, user_id, sleep_time, wake_time, score, enter_date from sleeps where 
  		user_id = $1 and enter_date between 
    	(date_trunc('week', current_date) - interval '$2 weeks') 
    	and 
    	(date_trunc('week', current_date) - interval '$2 weeks' + interval '6 days') order by enter_date desc`, user_id, page)
	if err != nil {
		logger.Error(ctx, "got an internal error while getting week stats", repo.Namespace, zap.Error(err), zap.Int64("user_id", user_id))
		return nil, models.Internal(err)
	}
	defer rows.Close()
	sleeps := make([]models.Sleep, 0, 7)
	for rows.Next() {
		var sleep models.Sleep
		if err := rows.Scan(&sleep.Id, &sleep.UserId, &sleep.SleepTime, &sleep.WakeTime, &sleep.Score, &sleep.EnterDate); err != nil {
			logger.Error(ctx, "got an internal error while scanning sleep record row",
				repo.Namespace,
				zap.Error(err),
				zap.Int64("user_id", user_id))
			return nil, models.Internal(err)
		}
		sleeps = append(sleeps, sleep)
	}
	if err := rows.Err(); err != nil {
		logger.Error(ctx, "got an internal error after scanning sleep record rows",
			repo.Namespace,
			zap.Error(err),
			zap.Int64("user_id", user_id))
		return nil, models.Internal(err)
	}
	return sleeps, nil
}

// Year implements repo.Sleep.
func (Repo) Year(ctx context.Context, tx *sql.Tx, user_id int64, page int) ([]models.SleepScore, models.Error) {
	rows, err := tx.QueryContext(ctx, `select score, enter_date from sleeps where 
  		user_id = $1 and enter_date between 
    	(date_trunc('year', current_date) - interval '$2 years') 
		and
    	(date_trunc('year', current_date) - interval '$2 years' + interval '1 year' - interval '1 day')  order by enter_date desc`, user_id, page)
	if err != nil {
		logger.Error(ctx, "got an internal error while getting week stats", repo.Namespace, zap.Error(err), zap.Int64("user_id", user_id))
		return nil, models.Internal(err)
	}
	defer rows.Close()
	scores := make([]models.SleepScore, 0, 365)
	for rows.Next() {
		var score models.SleepScore
		if err := rows.Scan(&score.Score, &score.EnterDate); err != nil {
			logger.Error(ctx, "got an internal error while scanning score row",
				repo.Namespace,
				zap.Error(err),
				zap.Int64("user_id", user_id))
			return nil, models.Internal(err)
		}
		scores = append(scores, score)
	}
	if err := rows.Err(); err != nil {
		logger.Error(ctx, "got an internal error after scanning scores rows",
			repo.Namespace,
			zap.Error(err),
			zap.Int64("user_id", user_id))
		return nil, models.Internal(err)
	}
	return scores, nil
}

var _ repo.Sleep = Repo{}
