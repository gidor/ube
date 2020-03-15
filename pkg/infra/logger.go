package infra

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger represents a logger
type Logger struct {
	Log *zap.Logger
}

var logger *Logger
var err error

// init the singletone logger
func init() {
	level := 0
	switch Getenv("loglevel", "info") {
	case "error":
		level = 2
	case "warning":
		level = 1
	case "info":
		level = 0
	case "debug":
		level = -1
	}

	// level 0 is a info level ,so debug level doesn't show
	// debug level can show in level -1

	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zapcore.Level(level))
	tlogger, err := config.Build()
	if err == nil {
		defer tlogger.Sync()
		logger = &Logger{
			Log: tlogger,
		}
	}

}
func GetLogger() (*Logger, error) {
	return logger, err
}

// CreateLogger creates a logger instance for all components
func CreateLogger(level int) (*Logger, error) {
	return logger, err
	// config := zap.NewProductionConfig()
	// config.Level = zap.NewAtomicLevelAt(zapcore.Level(level))
	// logger, err := config.Build()
	// if err != nil {
	// 	return nil, err
	// }
	// defer logger.Sync()
	// return &Logger{
	// 	Log: logger,
	// }, nil
}

// WithFields creates an entry from the standard logger and adds multiple fields to it
func (l *Logger) WithFields(fields ...zap.Field) *zap.Logger {
	return l.Log.With(fields...)
}
