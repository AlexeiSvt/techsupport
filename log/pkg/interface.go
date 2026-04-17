package pkg

// Logger defines a structured logging interface used across the application.
// It abstracts the underlying logging implementation (e.g., Zap, Logrus)
// to allow loose coupling and easier testing/mocking.
type Logger interface {
	// Infow logs an informational message with structured key-value pairs.
	// Used for general application flow and state tracking.
	Infow(msg string, args ...any)

	// Errorw logs an error-level message with structured context.
	// Used for reporting failures and unexpected errors.
	Errorw(msg string, args ...any)

	// Debugw logs detailed debug information.
	// Intended for development and troubleshooting only.
	Debugw(msg string, args ...any)

	// Fatalw logs a fatal error and terminates the application.
	// Should be used only for unrecoverable situations.
	Fatalw(msg string, args ...any)

	// Warnw logs warning-level messages indicating potential issues.
	// Execution can continue after warnings.
	Warnw(msg string, args ...any)
}