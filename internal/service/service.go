package service

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/vandi37/SleepTracker/internal/repo"
	"github.com/vandi37/SleepTracker/internal/repo/user_repo"
	"github.com/vandi37/SleepTracker/models"
	"github.com/vandi37/SleepTracker/pkg/logger"
	"github.com/vandi37/SleepTracker/pkg/score"
	"github.com/vandi37/SleepTracker/pkg/tokens"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	UserRepo              repo.User
	FriendRepo            repo.Friend
	SleepRepo             repo.Sleep
	AccessJwt, RefreshJwt *tokens.TokenService
	Database              *sql.DB
}

var Namespace = zap.Namespace("service")

func (s *Service) tokens(ctx context.Context, id int64) (models.UserWithToken, models.Error) {
	expires := time.Now().Add(-s.AccessJwt.GetExpiration())
	access, err := s.AccessJwt.Generate(strconv.FormatInt(id, 16))
	if err != nil {
		logger.Error(ctx, "got an internal error while generating access token", Namespace, zap.Error(err), zap.Int64("id", id))
		return models.UserWithToken{}, models.Internal(err)
	}
	refresh, err := s.RefreshJwt.Generate(strconv.FormatInt(id, 16))
	if err != nil {
		logger.Error(ctx, "got an internal error while generating refresh token", Namespace, zap.Error(err), zap.Int64("id", id))
		return models.UserWithToken{}, models.Internal(err)
	}
	return models.UserWithToken{
		Id:      id,
		Access:  access,
		Expires: expires,
		Refresh: refresh,
	}, nil
}

func (s *Service) Register(ctx context.Context, username, nickname, password string, birth time.Time) (models.UserWithToken, models.Error) {
	if len(password) > 72 {
		logger.Debug(ctx, "creation of account failed, invalid password", Namespace, zap.String("username", username), zap.String("nickname", nickname))
		return models.UserWithToken{}, models.Invalid("password", len(password))
	}
	password_hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		logger.Error(ctx, "got an internal error while hashing password", Namespace, zap.Error(err), zap.String("username", username), zap.String("nickname", nickname))
		return models.UserWithToken{}, models.Internal(err)
	}
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(ctx, "got an internal error while beginning transaction", Namespace, zap.Error(err))
		return models.UserWithToken{}, models.Internal(err)
	}
	id, rErr := s.UserRepo.Create(ctx, tx, username, nickname, password_hash, birth)
	if rErr != nil {
		tx.Rollback()
		logger.Debug(ctx, "error creating user", Namespace, zap.Error(rErr), zap.String("username", username), zap.String("nickname", nickname))
		return models.UserWithToken{}, rErr
	}
	tx.Commit()
	logger.Debug(ctx, "created user", Namespace, zap.Int64("id", id), zap.String("username", username), zap.String("nickname", nickname))
	return s.tokens(ctx, id)
}

func (s *Service) Login(ctx context.Context, username, password string) (models.UserWithToken, models.Error) {
	if len(password) > 72 {
		logger.Debug(ctx, "login failed, invalid password", Namespace, zap.String("username", username))
		return models.UserWithToken{}, models.Invalid("password", len(password))
	}
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(ctx, "got an internal error while beginning transaction", Namespace, zap.Error(err))
		return models.UserWithToken{}, models.Internal(err)
	}
	id, passwordHash, rErr := s.UserRepo.GetByUsername(ctx, tx, username)
	if rErr != nil {
		tx.Rollback()
		logger.Debug(ctx, "error getting user", Namespace, zap.Error(rErr), zap.String("username", username))
		return models.UserWithToken{}, rErr
	}
	tx.Commit()
	err = bcrypt.CompareHashAndPassword(passwordHash, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		logger.Debug(ctx, "mismatch password", Namespace, zap.Error(rErr), zap.Int64("id", id), zap.String("username", username))
		return models.UserWithToken{}, user_repo.InvalidCredentials{}
	} else if err != nil {
		logger.Error(ctx, "got an internal error while comparing password", Namespace, zap.Error(err))
		return models.UserWithToken{}, models.Internal(err)
	}
	logger.Debug(ctx, "login into user", Namespace, zap.Int64("id", id), zap.String("username", username))
	return s.tokens(ctx, id)
}

func (s *Service) Refresh(ctx context.Context, refresh string) (models.UserWithToken, models.Error) {
	token, err := s.RefreshJwt.Parse(refresh)
	if err != nil {
		logger.Debug(ctx, "error parsing jwt while refreshing a token", Namespace, zap.Error(err))
		return models.UserWithToken{}, handleJWTError(err)
	}
	if !token.Valid {
		logger.Debug(ctx, "got an invalid refresh token", Namespace)
		return models.UserWithToken{}, models.JsonError{Status: http.StatusUnauthorized, Message: "invalid token"}
	}
	strId, err := token.Claims.GetSubject()
	if err != nil {
		logger.Debug(ctx, "got an invalid refresh token without subject", Namespace)
		return models.UserWithToken{}, models.JsonError{Status: http.StatusUnauthorized, Message: "invalid token"}
	}
	id, err := strconv.ParseInt(strId, 16, 64)
	if err != nil {
		logger.Debug(ctx, "got an invalid refresh token with a non int subject", Namespace)
		return models.UserWithToken{}, models.Invalid("id", strId)
	}
	expires := time.Now().Add(-s.AccessJwt.GetExpiration())
	access, err := s.AccessJwt.Generate(strconv.FormatInt(id, 16))
	if err != nil {
		logger.Error(ctx, "got an internal error while generating access token", Namespace, zap.Error(err), zap.Int64("id", id))
		return models.UserWithToken{}, models.Internal(err)
	}
	logger.Debug(ctx, "refreshed a token", Namespace, zap.Int64("id", id))
	return models.UserWithToken{
		Id:      id,
		Expires: expires,
		Access:  access,
	}, nil
}

func (s *Service) Valid(ctx context.Context, access string) (int64, models.Error) {
	token, err := s.AccessJwt.Parse(access)
	if err != nil {
		logger.Debug(ctx, "error parsing jwt while validating an access token", Namespace, zap.Error(err))
		return 0, handleJWTError(err)
	}
	if !token.Valid {
		logger.Debug(ctx, "got an invalid access token", Namespace)
		return 0, models.JsonError{Status: http.StatusUnauthorized, Message: "invalid token"}
	}
	strId, err := token.Claims.GetSubject()
	if err != nil {
		logger.Debug(ctx, "got an invalid access token without subject", Namespace)
		return 0, models.JsonError{Status: http.StatusUnauthorized, Message: "invalid token"}
	}
	id, err := strconv.ParseInt(strId, 16, 64)
	if err != nil {
		logger.Debug(ctx, "got an invalid access token with a non int subject", Namespace)
		return 0, models.Invalid("id", strId)
	}
	logger.Debug(ctx, "validated an access token", Namespace, zap.Int64("id", id))
	return id, nil
}
func (s *Service) GetUser(ctx context.Context, id int64) (models.User, models.Error) {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(ctx, "got an internal error while beginning transaction", Namespace, zap.Error(err))
		return models.User{}, models.Internal(err)
	}
	user, rErr := s.UserRepo.Get(ctx, tx, id)
	if rErr != nil {
		tx.Rollback()
		logger.Debug(ctx, "error getting user", Namespace, zap.Error(rErr), zap.Int64("id", id))
		return models.User{}, rErr
	}
	tx.Commit()
	logger.Debug(ctx, "got a user", Namespace, zap.Int64("id", id))
	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, upd models.User) models.Error {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(ctx, "got an internal error while beginning transaction", Namespace, zap.Error(err))
		return models.Internal(err)
	}
	rErr := s.UserRepo.Update(ctx, tx, upd.Id, upd.Username, upd.Nickname, time.Time(upd.Birth))
	if rErr != nil {
		tx.Rollback()
		logger.Debug(ctx, "error updating user", Namespace, zap.Error(rErr), zap.Int64("id", upd.Id))
		return rErr
	}
	tx.Commit()
	logger.Debug(ctx, "updated a user", Namespace, zap.Int64("id", upd.Id))
	return nil
}

func (s *Service) UpdatePassword(ctx context.Context, id int64, password string) models.Error {
	if len(password) > 72 {
		logger.Debug(ctx, "updating password failed, invalid password", Namespace, zap.Int64("id", id))
		return models.Invalid("password", len(password))
	}
	password_hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		logger.Error(ctx, "got an internal error while hashing password", Namespace, zap.Error(err), zap.Int64("id", id))
		return models.Internal(err)
	}
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(ctx, "got an internal error while beginning transaction", Namespace, zap.Error(err))
		return models.Internal(err)
	}
	rErr := s.UserRepo.UpdatePassword(ctx, tx, id, password_hash)
	if rErr != nil {
		tx.Rollback()
		logger.Debug(ctx, "error updating user password", Namespace, zap.Error(rErr), zap.Int64("id", id))
		return rErr
	}
	tx.Commit()
	logger.Debug(ctx, "updated user password", Namespace, zap.Int64("id", id))
	return nil
}

func (s *Service) DeleteUser(ctx context.Context, id int64) models.Error {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(ctx, "got an internal error while beginning transaction", Namespace, zap.Error(err))
		return models.Internal(err)
	}
	rErr := s.UserRepo.Delete(ctx, tx, id)
	if rErr != nil {
		tx.Rollback()
		logger.Debug(ctx, "error deleting user", Namespace, zap.Error(rErr), zap.Int64("id", id))
		return rErr
	}
	tx.Commit()
	logger.Debug(ctx, "deleted a user", Namespace, zap.Int64("id", id))
	return nil
}

func (s *Service) Request(ctx context.Context, user1_id, user2_id int64) (int64, models.Error) {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(ctx, "got an internal error while beginning transaction", Namespace, zap.Error(err))
		return 0, models.Internal(err)
	}
	id, rErr := s.FriendRepo.Request(ctx, tx, user1_id, user2_id)
	if rErr != nil {
		tx.Rollback()
		logger.Debug(ctx, "error requesting a friendship", Namespace, zap.Error(rErr), zap.Int64("user1_id", user1_id), zap.Int64("user2_id", user2_id))
		return 0, rErr
	}
	tx.Commit()
	logger.Debug(ctx, "requested a friendship", Namespace, zap.Int64("user1_id", user1_id), zap.Int64("user2_id", user2_id))
	return id, nil
}

func (s *Service) AcceptFriendship(ctx context.Context, id, user2_id int64) models.Error {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(ctx, "got an internal error while beginning transaction", Namespace, zap.Error(err))
		return models.Internal(err)
	}
	rErr := s.FriendRepo.Accept(ctx, tx, id, user2_id)
	if rErr != nil {
		tx.Rollback()
		logger.Debug(ctx, "error accepting a friendship", Namespace, zap.Error(rErr), zap.Int64("user2_id", user2_id), zap.Int64("id", id))
		return rErr
	}
	tx.Commit()
	logger.Debug(ctx, "accepted a friendship", Namespace, zap.Int64("user2_id", user2_id), zap.Int64("id", id))
	return nil
}

func (s *Service) DeleteFriendship(ctx context.Context, id, user_id int64) models.Error {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(ctx, "got an internal error while beginning transaction", Namespace, zap.Error(err))
		return models.Internal(err)
	}
	rErr := s.FriendRepo.Delete(ctx, tx, id, user_id)
	if rErr != nil {
		tx.Rollback()
		logger.Debug(ctx, "error deleting a friendship", Namespace, zap.Error(rErr), zap.Int64("user_id", user_id), zap.Int64("id", id))
		return rErr
	}
	tx.Commit()
	logger.Debug(ctx, "deleted a friendship", Namespace, zap.Int64("user_id", user_id), zap.Int64("id", id))
	return nil
}

func (s *Service) GetFriendships(ctx context.Context, get models.GetFriendships) ([]models.Friend, models.Error) {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(ctx, "got an internal error while beginning transaction", Namespace, zap.Error(err))
		return nil, models.Internal(err)
	}
	fr, rErr := s.FriendRepo.Get(ctx, tx, get.UserId, get.Limit, get.Offset)
	if rErr != nil {
		tx.Rollback()
		logger.Debug(ctx, "error getting friendships", Namespace, zap.Error(rErr), zap.Int64("user_id", get.UserId))
		return nil, rErr
	}
	tx.Commit()
	logger.Debug(ctx, "got friendships", Namespace, zap.Int64("user_id", get.UserId))
	return fr, nil
}

func (s *Service) EnterSleep(ctx context.Context, enter models.Sleep) (int64, models.Error) {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(ctx, "got an internal error while beginning transaction", Namespace, zap.Error(err))
		return 0, models.Internal(err)
	}

	var enterScore int16
	if enter.SleepTime.Valid && enter.WakeTime.Valid {
		birth, rErr := s.UserRepo.GetBirth(ctx, tx, enter.UserId)
		if rErr != nil {
			tx.Rollback()
			logger.Debug(ctx, "error getting birth", Namespace, zap.Error(rErr), zap.Int64("user_id", enter.UserId), zap.String("date", enter.EnterDate.String()))
			return 0, rErr
		}
		enterScore = score.CalculateSleepScore(score.GetAge(birth), enter.SleepTime.Int16, enter.WakeTime.Int16)
	}
	id, rErr := s.SleepRepo.Enter(ctx, tx, enter.UserId, enter.SleepTime, enter.WakeTime, enterScore, time.Time(enter.EnterDate))
	if rErr != nil {
		tx.Rollback()
		logger.Debug(ctx, "error entering sleep", Namespace, zap.Error(rErr), zap.Int64("user_id", enter.UserId), zap.String("date", enter.EnterDate.String()))
		return 0, rErr
	}
	tx.Commit()
	logger.Debug(ctx, "got friendships", Namespace, zap.Int64("user_id", enter.UserId), zap.String("date", enter.EnterDate.String()))
	return id, nil
}

func (s *Service) UpdateSleep(ctx context.Context, upd models.Sleep) models.Error {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(ctx, "got an internal error while beginning transaction", Namespace, zap.Error(err))
		return models.Internal(err)
	}
	var enterScore int16
	if upd.SleepTime.Valid && upd.WakeTime.Valid {
		birth, rErr := s.UserRepo.GetBirth(ctx, tx, upd.UserId)
		if rErr != nil {
			tx.Rollback()
			logger.Debug(ctx, "error getting birth", Namespace, zap.Error(rErr), zap.Int64("user_id", upd.UserId), zap.String("date", upd.EnterDate.String()))
			return rErr
		}
		enterScore = score.CalculateSleepScore(score.GetAge(birth), upd.SleepTime.Int16, upd.WakeTime.Int16)
	}
	rErr := s.SleepRepo.Update(ctx, tx, upd.Id, upd.UserId, upd.SleepTime, upd.WakeTime, enterScore)
	if rErr != nil {
		tx.Rollback()
		logger.Debug(ctx, "error updating sleep record", Namespace, zap.Error(rErr), zap.Int64("id", upd.Id))
		return rErr
	}
	tx.Commit()
	logger.Debug(ctx, "updated a sleep record", Namespace, zap.Int64("id", upd.Id))
	return nil
}

func (s *Service) DeleteSleep(ctx context.Context, id, user_id int64) models.Error {
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(ctx, "got an internal error while beginning transaction", Namespace, zap.Error(err))
		return models.Internal(err)
	}
	rErr := s.SleepRepo.Delete(ctx, tx, id, user_id)
	if rErr != nil {
		tx.Rollback()
		logger.Debug(ctx, "error deleting a sleep record", Namespace, zap.Error(rErr), zap.Int64("id", id))
		return rErr
	}
	tx.Commit()
	logger.Debug(ctx, "deleted a sleep record", Namespace, zap.Int64("id", id))
	return nil
}
