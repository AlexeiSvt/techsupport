package pkg

type Logger interface {
	Infow(msg string, args ...any)
	Errorw(msg string, args ...any)
	Debugw(msg string, args ...any)
	Fatalw(msg string, args ...any)
	Warnw(msg string, args ...any)
}
