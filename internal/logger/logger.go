package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init initialises the global Uber Zap logger.
// Call this once at application startup.
func Init() {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	log, err = cfg.Build()
	if err != nil {
		panic("failed to initialise logger: " + err.Error())
	}
}

// Info logs an informational message.
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Error logs an error message.
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// Fatal logs a fatal message then calls os.Exit(1).
func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// With returns a logger with pre-set fields (useful for request-scoped logging).
func With(fields ...zap.Field) *zap.Logger {
	return log.With(fields...)
}

// Sync flushes buffered log entries. Call on application shutdown.
func Sync() {
	_ = log.Sync()
}
