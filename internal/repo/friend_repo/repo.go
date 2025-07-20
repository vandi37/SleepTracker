package friend_repo

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/vandi37/SleepTracker/internal/repo"
	"github.com/vandi37/SleepTracker/internal/repo/user_repo"
	"github.com/vandi37/SleepTracker/models"
	"github.com/vandi37/SleepTracker/pkg/logger"
	"go.uber.org/zap"
)

type Repo struct{}

// GetSecond implements repo.Friend.
func (r Repo) GetSecond(ctx context.Context, tx *sql.Tx, id, user_id int64) (int64, models.Error) {
	var getId int64
	err := tx.QueryRowContext(ctx, `select
	case 
        when user1_id = $2 then user2_id
        when user2_id = $2 then user1_id
    end
	from friends
	where id = $1 
	and (user1_id = $2 or user2_id = $2); `, id).
		Scan(&getId)
	if err == sql.ErrNoRows {
		return getId, FriendshipNotFound(id)
	}
	if err != nil {
		logger.Error(ctx, "got an internal error while getting user", repo.Namespace, zap.Error(err), zap.Int64("id", id))
		return getId, models.Internal(err)
	}
	return getId, nil
}

// Accept implements repo.Friend.
func (Repo) Accept(ctx context.Context, tx *sql.Tx, id int64, user2_id int64) models.Error {
	res, err := tx.ExecContext(ctx, `update friends set is_accepted = true, updated_at = default 
		where id = $1 and user2_id = $2 and not is_accepted`, id, user2_id)
	if err != nil {
		logger.Error(ctx, "got an internal error while accepting friendship request", repo.Namespace, zap.Error(err), zap.Int64("id", id), zap.Int64("user2_id", user2_id))
		return models.Internal(err)
	} else if ok, err := repo.CheckRes(res, repo.Equals(1)); err != nil {
		logger.Error(ctx, "got an internal error while checking result of accepting friendship request",
			zap.Error(err),
			zap.Int64("id", id),
			zap.Int64("user2_id", user2_id),
		)
		return models.Internal(err)
	} else if ok {
		return nil
	}
	var isAccepted bool
	if err := tx.QueryRowContext(ctx, `select is_accepted from friends where id = $1`, id).Scan(&isAccepted); err == sql.ErrNoRows {
		return FriendshipNotFound(id)
	} else if err != nil {
		logger.Error(ctx, "got an internal error while getting cause of failed update", repo.Namespace, zap.Error(err), zap.Int64("id", id), zap.Int64("user2_id", user2_id))
		return models.Internal(err)
	} else if isAccepted {
		return FriendshipAlreadyAccepted(id)
	}
	return models.NotAllowed{Id: id, User: user2_id, Thing: "accept a friendship request"}
}

// Delete implements repo.Friend.
func (Repo) Delete(ctx context.Context, tx *sql.Tx, id int64, user_id int64) models.Error {
	res, err := tx.ExecContext(ctx, `delete from friends where id = $1 and (user1_id = $2 or user2_id = $2) `, id, user_id)
	if err != nil {
		logger.Error(ctx, "got an internal error while deleting friendship or friendship request",
			repo.Namespace, zap.Error(err), zap.Int64("id", id), zap.Int64("user_id", user_id))
		return models.Internal(err)
	} else if ok, err := repo.CheckRes(res, repo.Equals(1)); err != nil {
		logger.Error(ctx, "got an internal error while checking result of deleting friendship or friendship request",
			repo.Namespace,
			zap.Error(err),
			zap.Int64("id", id),
			zap.Int64("user_id", user_id),
		)
		return models.Internal(err)
	} else if ok {
		return nil
	}
	if err := tx.QueryRowContext(ctx, `select 1 from friends where id = $1`, id).Err(); err == sql.ErrNoRows {
		return FriendshipNotFound(id)
	} else if err != nil {
		logger.Error(ctx, "got an internal error while getting cause of failed delete",
			repo.Namespace, zap.Error(err), zap.Int64("id", id), zap.Int64("user_id", user_id))
		return models.Internal(err)
	}
	return models.NotAllowed{Id: id, User: user_id, Thing: "delete a friendship"}
}

// Get implements repo.Friend.
func (Repo) Get(ctx context.Context, tx *sql.Tx, user_id int64, limit, offset int) ([]models.Friend, models.Error) {
	rows, err := tx.QueryContext(ctx, `
        select 
            f.id,
            f.is_accepted,
            f.created_at,
            f.updated_at,
			u1.id,
			u1.username,
			u1.nickname,
			u1.birth,
			u1.created_at,
			u2.id,
			u2.username,
			u2.nickname,
			u2.birth,
			u2.created_at
        from friends f
        inner join users u1 on f.user1_id = u1.id
        inner join users u2 on f.user2_id = u2.id
        where f.user1_id = $1 or f.user2_id = $1
        order by f.updated_at desc limit $2 offset $3`, user_id, limit, offset)
	if err != nil {
		logger.Error(ctx, "got an internal error while getting friendships and friendship requests",
			repo.Namespace,
			zap.Error(err),
			zap.Int64("user_id", user_id),
		)
		return nil, models.Internal(err)
	}
	defer rows.Close()
	var friends []models.Friend
	for rows.Next() {
		var friend models.Friend
		var user1 models.User
		var user2 models.User

		err := rows.Scan(
			&friend.Id,
			&friend.IsAccepted,
			&friend.CreatedAt,
			&friend.UpdatedAt,
			&user1.Id,
			&user1.Username,
			&user1.Nickname,
			&user1.Birth,
			&user1.CreatedAt,
			&user2.Id,
			&user2.Username,
			&user2.Nickname,
			&user2.Birth,
			&user2.CreatedAt,
		)
		if err != nil {
			logger.Error(ctx, "got an internal error while scanning friend row",
				repo.Namespace,
				zap.Error(err),
				zap.Int64("user_id", user_id))
			return nil, models.Internal(err)
		}
		friend.User1 = user1
		friend.User2 = user2
		friends = append(friends, friend)
	}

	if err := rows.Err(); err != nil {
		logger.Error(ctx, "got an internal error after scanning friend rows",
			repo.Namespace,
			zap.Error(err),
			zap.Int64("user_id", user_id))
		return nil, models.Internal(err)
	}
	return friends, nil
}

// Request implements repo.Friend.
func (Repo) Request(ctx context.Context, tx *sql.Tx, user1_id int64, user2_id int64) (int64, models.Error) {
	if user1_id == user2_id {
		return 0, SelfRequest(user1_id)
	}
	var id int64

	err := tx.QueryRowContext(ctx, `insert into friends (user1_id, user2_id) values ($1, $2) returning id`, user1_id, user2_id).
		Scan(&id)
	if err == sql.ErrNoRows {
		logger.Error(ctx, "got no rows after requesting a friendship",
			repo.Namespace,
			zap.Int64("user1_id", user1_id),
			zap.Int64("user2_id", user2_id),
		)
		return 0, models.Internal(repo.NoRows)
	} else if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
		return 0, FriendshipAlreadyExists{user1_id, user2_id}
	} else if ok && pgErr.Code == "23503" {
		return 0, user_repo.UserNotFound(user2_id) // might be also from but usually is to
	} else if err != nil {
		logger.Error(ctx, "got an internal error while requesting a friendship",
			repo.Namespace,
			zap.Error(err),
			zap.Int64("user1_id", user1_id),
			zap.Int64("user2_id", user2_id),
		)
		return 0, models.Internal(err)
	}
	return id, nil
}

var _ repo.Friend = Repo{}
