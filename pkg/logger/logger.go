package logger

import (
	"fmt"

	"github.com/Myakun/personal-secretary-user-api/pkg/env"
	"go.uber.org/zap"
)

type Logger interface {
	Debug(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Info(args ...interface{})
	Warning(args ...interface{})

	DebugW(msg string, keysAndValues ...interface{})
	ErrorW(msg string, keysAndValues ...interface{})
	FatalW(msg string, keysAndValues ...interface{})
	InfoW(msg string, keysAndValues ...interface{})
	WarningW(msg string, keysAndValues ...interface{})

	DebugWithTag(msg string, tag string)
	ErrorWithTag(msg string, tag string)
	FatalWithTag(msg string, tag string)
	InfoWithTag(msg string, tag string)
	WarningWithTag(msg string, tag string)

	DebugWithTagW(msg string, tag string, keysAndValues ...interface{})
	ErrorWithTagW(msg string, tag string, keysAndValues ...interface{})
	FatalWithTagW(msg string, tag string, keysAndValues ...interface{})
	InfoWithTagW(msg string, tag string, keysAndValues ...interface{})
	WarningWithTagW(msg string, tag string, keysAndValues ...interface{})
}

type zapLogger struct {
	zapLogger *zap.SugaredLogger
}

func NewLogger(appEnv *env.Env) (Logger, error) {
	var base *zap.Logger
	var err error

	if appEnv.IsProd() || appEnv.IsStage() {
		base, err = zap.NewProduction()
	} else {
		base, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	return &zapLogger{zapLogger: base.Sugar()}, nil
}

func (l *zapLogger) Debug(args ...interface{}) {
	l.zapLogger.Debug(args...)
}

func (l *zapLogger) Error(args ...interface{}) {
	l.zapLogger.Error(args...)
}

func (l *zapLogger) Fatal(args ...interface{}) {
	l.zapLogger.Fatal(args...)
}

func (l *zapLogger) Info(args ...interface{}) {
	l.zapLogger.Info(args...)
}

func (l *zapLogger) Warning(args ...interface{}) {
	l.zapLogger.Warn(args...)
}

func (l *zapLogger) CriticalW(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Errorw(msg, keysAndValues...)
}

func (l *zapLogger) DebugW(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Debugw(msg, keysAndValues...)
}

func (l *zapLogger) ErrorW(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Errorw(msg, keysAndValues...)
}

func (l *zapLogger) FatalW(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) InfoW(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Infow(msg, keysAndValues...)
}

func (l *zapLogger) WarningW(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Warnw(msg, keysAndValues...)
}

func (l *zapLogger) DebugWithTag(msg string, tag string) {
	l.zapLogger.Debugw(msg, "tag", tag)
}

func (l *zapLogger) ErrorWithTag(msg string, tag string) {
	l.zapLogger.Errorw(msg, "tag", tag)
}

func (l *zapLogger) FatalWithTag(msg string, tag string) {
	l.zapLogger.Fatalw(msg, "tag", tag)
}

func (l *zapLogger) InfoWithTag(msg string, tag string) {
	l.zapLogger.Infow(msg, "tag", tag)
}

func (l *zapLogger) WarningWithTag(msg string, tag string) {
	l.zapLogger.Warnw(msg, "tag", tag)
}

func (l *zapLogger) DebugWithTagW(msg string, tag string, keysAndValues ...interface{}) {
	args := append([]interface{}{"tag", tag}, keysAndValues...)
	l.zapLogger.Debugw(msg, args...)
}

func (l *zapLogger) EmergencyWithTagW(msg string, tag string, keysAndValues ...interface{}) {
	args := append([]interface{}{"tag", tag}, keysAndValues...)
	l.zapLogger.Fatalw(msg, args...)
}

func (l *zapLogger) ErrorWithTagW(msg string, tag string, keysAndValues ...interface{}) {
	args := append([]interface{}{"tag", tag}, keysAndValues...)
	l.zapLogger.Errorw(msg, args...)
}

func (l *zapLogger) FatalWithTagW(msg string, tag string, keysAndValues ...interface{}) {
	args := append([]interface{}{"tag", tag}, keysAndValues...)
	l.zapLogger.Fatalw(msg, args...)
}

func (l *zapLogger) InfoWithTagW(msg string, tag string, keysAndValues ...interface{}) {
	args := append([]interface{}{"tag", tag}, keysAndValues...)
	l.zapLogger.Infow(msg, args...)
}

func (l *zapLogger) WarningWithTagW(msg string, tag string, keysAndValues ...interface{}) {
	args := append([]interface{}{"tag", tag}, keysAndValues...)
	l.zapLogger.Warnw(msg, args...)
}
