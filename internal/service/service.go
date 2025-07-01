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
		return models.UserWithToken{}, models.Invalid("password", len(password))
	}
	password_hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		logger.Error(ctx, "got an internal error while hashing password", Namespace, zap.Error(err))
		return models.UserWithToken{}, models.Internal(err)
	}
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(ctx, "got an internal error while beginning transaction", Namespace, zap.Error(err))
		return models.UserWithToken{}, models.Internal(err)
	}
	id, crErr := s.UserRepo.Create(ctx, tx, username, nickname, password_hash, birth)
	if crErr != nil {
		tx.Rollback()
		logger.Debug(ctx, "error creating user", Namespace, zap.Error(crErr))
		return models.UserWithToken{}, crErr
	}
	tx.Commit()

	return s.tokens(ctx, id)
}

func (s *Service) Login(ctx context.Context, username, password string) (models.UserWithToken, models.Error) {
	if len(password) > 72 {
		return models.UserWithToken{}, models.Invalid("password", len(password))
	}
	tx, err := s.Database.BeginTx(ctx, nil)
	if err != nil {
		logger.Error(ctx, "got an internal error while beginning transaction", Namespace, zap.Error(err))
		return models.UserWithToken{}, models.Internal(err)
	}
	id, passwordHash, gtErr := s.UserRepo.GetByUsername(ctx, tx, username)
	if gtErr != nil {
		tx.Rollback()
		return models.UserWithToken{}, gtErr
	}
	tx.Commit()
	err = bcrypt.CompareHashAndPassword(passwordHash, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return models.UserWithToken{}, user_repo.InvalidCredentials{}
	} else if err != nil {
		logger.Error(ctx, "got an internal error while comparing password", Namespace, zap.Error(err))
		return models.UserWithToken{}, models.Internal(err)
	}
	return s.tokens(ctx, id)
}

func (s *Service) Refresh(ctx context.Context, refresh string) (models.UserWithToken, models.Error) {
	token, err := s.RefreshJwt.Parse(refresh)
	if err != nil {
		return models.UserWithToken{}, handleJWTError(err)
	}
	if !token.Valid {
		return models.UserWithToken{}, models.JsonError{Status: http.StatusUnauthorized, Message: "invalid token"}
	}
	strId, err := token.Claims.GetSubject()
	if err != nil {
		return models.UserWithToken{}, models.JsonError{Status: http.StatusUnauthorized, Message: "invalid token"}
	}
	id, err := strconv.ParseInt(strId, 16, 64)
	if err != nil {
		return models.UserWithToken{}, models.Invalid("id", strId)
	}
	expires := time.Now().Add(-s.AccessJwt.GetExpiration())
	access, err := s.AccessJwt.Generate(strconv.FormatInt(id, 16))
	if err != nil {
		logger.Error(ctx, "got an internal error while generating access token", Namespace, zap.Error(err), zap.Int64("id", id))
		return models.UserWithToken{}, models.Internal(err)
	}
	return models.UserWithToken{
		Id:      id,
		Expires: expires,
		Access:  access,
	}, nil
}

func (s *Service) Valid(ctx context.Context, access string) (int64, models.Error) {
	token, err := s.AccessJwt.Parse(access)
	if err != nil {
		return 0, handleJWTError(err)
	}
	if !token.Valid {
		return 0, models.JsonError{Status: http.StatusUnauthorized, Message: "invalid token"}
	}
	strId, err := token.Claims.GetSubject()
	if err != nil {
		return 0,  models.JsonError{Status: http.StatusUnauthorized, Message: "invalid token"}
	}
	id, err := strconv.ParseInt(strId, 16, 64)
	if err != nil {
		return 0, models.Invalid("id", strId)
	}
	return id, nil
}
