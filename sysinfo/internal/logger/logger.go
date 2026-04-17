package logger

import (
	"techsupport/log/pkg"

	"go.uber.org/zap"
)

// zapWrapper is a thin adapter over zap.SugaredLogger.
// It implements the pkg.Logger interface and decouples the application
// from direct dependency on Uber Zap.
type zapWrapper struct {
	// sugar is the underlying Zap SugaredLogger used for structured logging.
	sugar *zap.SugaredLogger
}

// NewZapLogger initializes a new development-mode Zap logger
// and wraps it into a pkg.Logger implementation.
// It is intended for local development and debugging purposes.
func NewZapLogger() pkg.Logger {
	z, _ := zap.NewDevelopment()

	return &zapWrapper{
		sugar: z.Sugar(),
	}
}

// Debugw logs a debug-level message with structured key-value pairs.
// Used for detailed diagnostic information during development.
func (z *zapWrapper) Debugw(msg string, keysAndValues ...any) {
	z.sugar.Debugw(msg, keysAndValues...)
}

// Infow logs an informational message with structured context.
// Used for tracking normal application flow.
func (z *zapWrapper) Infow(msg string, keysAndValues ...any) {
	z.sugar.Infow(msg, keysAndValues...)
}

// Warnw logs a warning-level message indicating a potential issue.
// Execution can continue after warning events.
func (z *zapWrapper) Warnw(msg string, keysAndValues ...any) {
	z.sugar.Warnw(msg, keysAndValues...)
}

// Errorw logs an error-level message with structured context.
// Used for failures that do not immediately terminate the application.
func (z *zapWrapper) Errorw(msg string, keysAndValues ...any) {
	z.sugar.Errorw(msg, keysAndValues...)
}

// Fatalw logs a fatal error message and terminates the application.
// Should be used only for unrecoverable situations.
func (z *zapWrapper) Fatalw(msg string, keysAndValues ...any) {
	z.sugar.Fatalw(msg, keysAndValues...)
}