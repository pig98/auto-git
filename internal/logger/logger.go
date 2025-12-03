package logger

import (
	"log"
	"os"
	"strings"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelError
)

var (
	currentLevel LogLevel = LogLevelInfo
)

// SetLevelFromEnv reads LOG_LEVEL from environment and sets the log level
func SetLevelFromEnv() {
	level := strings.ToUpper(strings.TrimSpace(os.Getenv("LOG_LEVEL")))
	switch level {
	case "DEBUG":
		currentLevel = LogLevelDebug
	case "ERROR":
		currentLevel = LogLevelError
	case "INFO", "":
		fallthrough
	default:
		currentLevel = LogLevelInfo
	}
}

// Debug logs a debug message
func Debug(format string, args ...interface{}) {
	if currentLevel <= LogLevelDebug {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// Info logs an info message
func Info(format string, args ...interface{}) {
	if currentLevel <= LogLevelInfo {
		log.Printf("[INFO] "+format, args...)
	}
}

// Error logs an error message (always logged)
func Error(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

// Fatal logs a fatal error and exits
func Fatal(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
