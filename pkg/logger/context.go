package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	ContextLogger = "logger-value"
)

func Context(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ContextLogger, l)
}

func FromCtx(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(ContextLogger).(*zap.Logger)
	if !ok {
		return nil
	}
	return logger
}

func Log(ctx context.Context, lvl zapcore.Level, msg string, fields ...zap.Field) bool {
	logger := FromCtx(ctx)
	if logger == nil {
		return false
	}

	logger.Log(lvl, msg, fields...)
	return true
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.DebugLevel, msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.InfoLevel, msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.WarnLevel, msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.ErrorLevel, msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.FatalLevel, msg, fields...)
}

func DPanic(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.DPanicLevel, msg, fields...)
}

func Panic(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.PanicLevel, msg, fields...)
}
