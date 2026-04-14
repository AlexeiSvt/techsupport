package logger

import (
	"techsupport/log/pkg"
	"go.uber.org/zap"
)

type zapWrapper struct {
	sugar *zap.SugaredLogger
}

func NewZapLogger() pkg.Logger {
	z, _ := zap.NewDevelopment()
	return &zapWrapper{
		sugar: z.Sugar(),
	}
}

func (z *zapWrapper) Debugw(msg string, keysAndValues ...any) { z.sugar.Debugw(msg, keysAndValues...) }
func (z *zapWrapper) Infow(msg string, keysAndValues ...any)   { z.sugar.Infow(msg, keysAndValues...) }
func (z *zapWrapper) Warnw(msg string, keysAndValues ...any)   { z.sugar.Warnw(msg, keysAndValues...) }
func (z *zapWrapper) Errorw(msg string, keysAndValues ...any)  { z.sugar.Errorw(msg, keysAndValues...) }
func (z *zapWrapper) Fatalw(msg string, keysAndValues ...any)  { z.sugar.Fatalw(msg, keysAndValues...) }