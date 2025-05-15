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

	DebugWith(msg string, keysAndValues ...interface{})
	ErrorWith(msg string, keysAndValues ...interface{})
	FatalWith(msg string, keysAndValues ...interface{})
	InfoWith(msg string, keysAndValues ...interface{})
	WarningWith(msg string, keysAndValues ...interface{})

	DebugLogTag(msg string, tag string)
	ErrorLogTag(msg string, tag string)
	FatalLogTag(msg string, tag string)
	InfoLogTag(msg string, tag string)
	WarningLogTag(msg string, tag string)

	DebugLogTagWith(msg string, tag string, keysAndValues ...interface{})
	ErrorLogTagWith(msg string, tag string, keysAndValues ...interface{})
	FatalLogTagWith(msg string, tag string, keysAndValues ...interface{})
	InfoLogTagWith(msg string, tag string, keysAndValues ...interface{})
	WarningLogTagWith(msg string, tag string, keysAndValues ...interface{})
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

func (l *zapLogger) CriticalWith(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Errorw(msg, keysAndValues...)
}

func (l *zapLogger) DebugWith(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Debugw(msg, keysAndValues...)
}

func (l *zapLogger) ErrorWith(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Errorw(msg, keysAndValues...)
}

func (l *zapLogger) FatalWith(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Fatalw(msg, keysAndValues...)
}

func (l *zapLogger) InfoWith(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Infow(msg, keysAndValues...)
}

func (l *zapLogger) WarningWith(msg string, keysAndValues ...interface{}) {
	l.zapLogger.Warnw(msg, keysAndValues...)
}

func (l *zapLogger) DebugLogTag(msg string, tag string) {
	l.zapLogger.Debugw(msg, "tag", tag)
}

func (l *zapLogger) ErrorLogTag(msg string, tag string) {
	l.zapLogger.Errorw(msg, "tag", tag)
}

func (l *zapLogger) FatalLogTag(msg string, tag string) {
	l.zapLogger.Fatalw(msg, "tag", tag)
}

func (l *zapLogger) InfoLogTag(msg string, tag string) {
	l.zapLogger.Infow(msg, "tag", tag)
}

func (l *zapLogger) WarningLogTag(msg string, tag string) {
	l.zapLogger.Warnw(msg, "tag", tag)
}

func (l *zapLogger) DebugLogTagWith(msg string, tag string, keysAndValues ...interface{}) {
	args := append([]interface{}{"tag", tag}, keysAndValues...)
	l.zapLogger.Debugw(msg, args...)
}

func (l *zapLogger) EmergencyLogTagWith(msg string, tag string, keysAndValues ...interface{}) {
	args := append([]interface{}{"tag", tag}, keysAndValues...)
	l.zapLogger.Fatalw(msg, args...)
}

func (l *zapLogger) ErrorLogTagWith(msg string, tag string, keysAndValues ...interface{}) {
	args := append([]interface{}{"tag", tag}, keysAndValues...)
	l.zapLogger.Errorw(msg, args...)
}

func (l *zapLogger) FatalLogTagWith(msg string, tag string, keysAndValues ...interface{}) {
	args := append([]interface{}{"tag", tag}, keysAndValues...)
	l.zapLogger.Fatalw(msg, args...)
}

func (l *zapLogger) InfoLogTagWith(msg string, tag string, keysAndValues ...interface{}) {
	args := append([]interface{}{"tag", tag}, keysAndValues...)
	l.zapLogger.Infow(msg, args...)
}

func (l *zapLogger) WarningLogTagWith(msg string, tag string, keysAndValues ...interface{}) {
	args := append([]interface{}{"tag", tag}, keysAndValues...)
	l.zapLogger.Warnw(msg, args...)
}
