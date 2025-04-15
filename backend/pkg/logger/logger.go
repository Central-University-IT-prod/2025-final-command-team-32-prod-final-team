package logger

import (
	"context"

	"go.uber.org/zap"
)

const (
	contextKey = "logger"
)

type Logger struct {
	l *zap.Logger
}

func CtxWithLogger(ctx context.Context) context.Context {
	l := Logger{}
	zaplogger, _ := zap.NewProduction()
	l.l = zaplogger
	return context.WithValue(ctx, contextKey, &l)
}

func FromCtx(ctx context.Context) *Logger {
	return ctx.Value(contextKey).(*Logger)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.l.Info(msg, fields...)

}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.l.Error(msg, fields...)
}
