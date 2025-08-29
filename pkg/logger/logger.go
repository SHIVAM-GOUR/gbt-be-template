package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger wraps logrus logger
type Logger struct {
	*logrus.Logger
}

// New creates a new logger instance
func New(level, format string) *Logger {
	log := logrus.New()

	// Set log level
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	// Set log format
	if format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	log.SetOutput(os.Stdout)

	return &Logger{log}
}

// WithFields creates a new logger entry with fields
func (l *Logger) WithFields(fields map[string]interface{}) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

// WithField creates a new logger entry with a single field
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithError creates a new logger entry with an error field
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

// HTTP creates a logger entry for HTTP requests
func (l *Logger) HTTP(method, path, userAgent, ip string, statusCode int, duration int64) *logrus.Entry {
	return l.WithFields(map[string]interface{}{
		"method":      method,
		"path":        path,
		"user_agent":  userAgent,
		"ip":          ip,
		"status_code": statusCode,
		"duration_ms": duration,
		"type":        "http_request",
	})
}

// Database creates a logger entry for database operations
func (l *Logger) Database(operation, table string, duration int64) *logrus.Entry {
	return l.WithFields(map[string]interface{}{
		"operation":   operation,
		"table":       table,
		"duration_ms": duration,
		"type":        "database",
	})
}

// Auth creates a logger entry for authentication events
func (l *Logger) Auth(userID uint, action, ip string) *logrus.Entry {
	return l.WithFields(map[string]interface{}{
		"user_id": userID,
		"action":  action,
		"ip":      ip,
		"type":    "auth",
	})
}
