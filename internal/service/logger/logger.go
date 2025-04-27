package logger

import (
	"fmt"
	"sync"
)

var (
	loggerServiceInstance         *Logger
	initLoggerServiceInstanceOnce sync.Once
)

type Logger struct {
}

func (logger *Logger) Fatal(msg string) {
	fmt.Println("\033[31mFatal\033[0m: " + msg)
}

func (logger *Logger) Critical(msg string) {
	fmt.Println("\033[31mCRITICAL\033[0m: " + msg)
}

func (logger *Logger) Debug(msg string) {
	fmt.Println("DEBUG: " + msg)
}

func (logger *Logger) DebugWithLogTag(msg string, logTag string) {
	fmt.Println("DEBUG [" + logTag + "]: " + msg)
}

func (logger *Logger) Emergency(msg string) {
	fmt.Println("EMERGENCY: " + msg)
}

func (logger *Logger) Error(msg string) {
	fmt.Println("\033[31mERROR\033[0m: " + msg)
}

func (logger *Logger) Info(msg string) {
	fmt.Println("INFO: " + msg)
}

func InitLoggerService() (*Logger, error) {
	initLoggerServiceInstanceOnce.Do(func() {
		loggerServiceInstance = &Logger{}
	})

	return loggerServiceInstance, nil
}

func GetLoggerService() *Logger {
	if loggerServiceInstance == nil {
		panic("Logger service not initialized")
	}
	return loggerServiceInstance
}
