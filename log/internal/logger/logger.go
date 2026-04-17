// Package logger provides a thin wrapper over Uber Zap SugaredLogger.
// It exposes structured logging methods used across the application.
package logger

import (
	"go.uber.org/zap"
)

// ZapWrapper wraps zap.SugaredLogger to provide a simplified logging interface.
// It is used to decouple the application from direct dependency on zap.
type ZapWrapper struct {
	// sugar is the underlying Zap sugared logger instance.
	sugar *zap.SugaredLogger
}

// NewZapLogger initializes a new development-mode Zap logger
// and wraps it into ZapWrapper.
// It is intended for local development and debugging usage.
func NewZapLogger() *ZapWrapper {
	logger, _ := zap.NewDevelopment()

	return &ZapWrapper{
		sugar: logger.Sugar(),
	}
}

// Infow logs an informational message with structured key-value pairs.
// It is used for general application flow tracking.
func (z *ZapWrapper) Infow(msg string, args ...any) {
	z.sugar.Infow(msg, args...)
}

// Errorw logs an error-level message with structured context.
// It is used for reporting failures and critical issues.
func (z *ZapWrapper) Errorw(msg string, args ...any) {
	z.sugar.Errorw(msg, args...)
}

// Debugw logs debug-level messages with structured context.
// It is intended for development and troubleshooting purposes.
func (z *ZapWrapper) Debugw(msg string, args ...any) {
	z.sugar.Debugw(msg, args...)
}

// Fatalw logs a fatal-level message and terminates the application.
// It should be used only for unrecoverable errors.
func (z *ZapWrapper) Fatalw(msg string, args ...any) {
	z.sugar.Fatalw(msg, args...)
}

// Warnw logs warning-level messages indicating potential issues.
// It is used when something unexpected happens but execution can continue.
func (z *ZapWrapper) Warnw(msg string, args ...any) {
	z.sugar.Warnw(msg, args...)
}