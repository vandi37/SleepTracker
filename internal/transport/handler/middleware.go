package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vandi37/SleepTracker/pkg/logger"
	"go.uber.org/zap"
)

func AddLogger(l *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(logger.ContextLogger, l)
		ctx.Next()
	}
}

func Logging(ctx *gin.Context) {
	start := time.Now()
	path := ctx.FullPath()
	if path == "" {
		path = ctx.Request.URL.String()
	}

	ctx.Next()

	latency := time.Since(start)

	clientIP := ctx.ClientIP()
	if clientIP == "" {
		clientIP = "-"
	}

	statusCode := ctx.Writer.Status()
	method := ctx.Request.Method

	fields := []zap.Field{
		zap.Int("status", statusCode),
		zap.Duration("latency", latency),
		zap.String("method", method),
		zap.String("path", path),
	}

	switch {
	case statusCode >= http.StatusInternalServerError:
		logger.Error(ctx, "request", fields...)
	default:
		logger.Info(ctx, "request", fields...)
	}

}

func (h *Handler) LoginMiddleware(ctx *gin.Context) {
	token := strings.TrimPrefix(ctx.Request.Header.Get("Authorization"), "Bearer ")
	id, err := h.service.Valid(ctx, token)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code(), err.JsonError())
		return
	}
	ctx.Set(ID_KEY, id)
	ctx.Next()
}

const ID_KEY = "login-id"
