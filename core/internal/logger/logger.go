// Package logger provides a concrete implementation of the system's logging interface.
// It uses Uber's Zap library under the hood for high-performance structured logging.
package logger

import (
	"techsupport/log/pkg"
	"go.uber.org/zap"
)

// zapWrapper implements the pkg.Logger interface by wrapping a zap.SugaredLogger.
// This abstraction allows the rest of the system to remain decoupled from the specific 
// logging library used, facilitating easier maintenance or library swaps in the future.
type zapWrapper struct {
	sugar *zap.SugaredLogger
}

// NewZapLogger initializes a new logger instance using Zap's development configuration.
// In development mode, logs are human-readable, include stack traces on errors, 
// and use color-coded output.
func NewZapLogger() pkg.Logger {
	// For production, this should be replaced with zap.NewProduction() 
	// or loaded from a configuration file.
	z, _ := zap.NewDevelopment()
	return &zapWrapper{
		sugar: z.Sugar(),
	}
}

// Debugw logs a message and a set of key-value pairs at the Debug level.
// Use this for detailed information useful during development and troubleshooting.
func (z *zapWrapper) Debugw(msg string, keysAndValues ...any) { z.sugar.Debugw(msg, keysAndValues...) }

// Infow logs a message and a set of key-value pairs at the Info level.
// This is the standard level for tracking general application flow and significant events.
func (z *zapWrapper) Infow(msg string, keysAndValues ...any)   { z.sugar.Infow(msg, keysAndValues...) }

// Warnw logs a message and a set of key-value pairs at the Warn level.
// Use this for unexpected events that do not stop the application but require attention.
func (z *zapWrapper) Warnw(msg string, keysAndValues ...any)   { z.sugar.Warnw(msg, keysAndValues...) }

// Errorw logs a message and a set of key-value pairs at the Error level.
// Use this for critical issues that impact specific operations or data integrity.
func (z *zapWrapper) Errorw(msg string, keysAndValues ...any)  { z.sugar.Errorw(msg, keysAndValues...) }

// Fatalw logs a message and a set of key-value pairs at the Fatal level, 
// then calls os.Exit(1). Use this for unrecoverable system-wide failures.
func (z *zapWrapper) Fatalw(msg string, keysAndValues ...any)  { z.sugar.Fatalw(msg, keysAndValues...) }