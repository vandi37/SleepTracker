package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/vandi37/SleepTracker/internal/config"
	"github.com/vandi37/SleepTracker/internal/repo/friend_repo"
	"github.com/vandi37/SleepTracker/internal/repo/sleep_repo"
	"github.com/vandi37/SleepTracker/internal/repo/user_repo"
	"github.com/vandi37/SleepTracker/internal/service"
	"github.com/vandi37/SleepTracker/internal/transport/handler"
	"github.com/vandi37/SleepTracker/pkg/logger"
	"github.com/vandi37/SleepTracker/pkg/tokens"
	"go.uber.org/zap"
)

var logPath = fmt.Sprintf("sleep-tracker-logs:%s.log", time.Now().UTC().Format(time.DateOnly))

func Run(ctx context.Context) {
	logger, file := logger.ConsoleAndFile(logPath)

	config, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("error loading config", zap.Error(err))
	}

	db, err := sql.Open("postgres", config.ConnString)
	if err != nil {
		logger.Fatal("error opening database", zap.Error(err))
	}
	err = db.PingContext(ctx)
	if err != nil {
		logger.Fatal("error pinging database", zap.Error(err))
	}
	accessExp, err := time.ParseDuration(config.Access.Expires)
	if err != nil {
		logger.Fatal("error loading access expr", zap.Error(err))
	}
	access := tokens.New(config.Access.Secret, accessExp)
	refreshExp, err := time.ParseDuration(config.Refresh.Expires)
	if err != nil {
		logger.Fatal("error loading refresh expr", zap.Error(err))
	}
	refresh := tokens.New(config.Refresh.Secret, refreshExp)

	service := service.Service{
		UserRepo:   user_repo.Repo{},
		FriendRepo: friend_repo.Repo{},
		SleepRepo:  sleep_repo.Repo{},
		AccessJwt:  access,
		RefreshJwt: refresh,
		Database:   db,
	}

	server := http.Server{Addr: fmt.Sprintf(":%d", config.Port), Handler: handler.RegisterHandler([]string{"*"}, handler.NewHandler(&service), logger).Handler()}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Fatal("error running router", zap.Error(err))
		}
	}()
	logger.Info("started server", zap.Int("port", config.Port))
	<-ctx.Done()
	ctx, close := context.WithTimeout(context.Background(), time.Minute)
	defer close()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("error shutting down the server", zap.Error(err))
	}

	if err := db.Close(); err != nil {
		logger.Fatal("error closing the database", zap.Error(err))
	}

	if err := file.Close(); err != nil {
		logger.Fatal("error closing the logging file", zap.Error(err))
	}
}
