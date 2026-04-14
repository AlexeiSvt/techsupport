package logger

import (
    "go.uber.org/zap"
)

type ZapWrapper struct {
    sugar *zap.SugaredLogger
}

func NewZapLogger() *ZapWrapper {
    logger, _ := zap.NewDevelopment()
    return &ZapWrapper{
        sugar: logger.Sugar(), 
    }
}

func (z *ZapWrapper) Infow(msg string, args ...any) {
    z.sugar.Infow(msg, args...)
}

func (z *ZapWrapper) Errorw(msg string, args ...any) {
    z.sugar.Errorw(msg, args...)
}

func (z *ZapWrapper) Debugw(msg string, args ...any) {
	z.sugar.Debugw(msg, args...)
}

func (z *ZapWrapper) Fatalw(msg string, args ...any) {
	z.sugar.Fatalw(msg, args...)
}

func (z *ZapWrapper) Warnw(msg string, args ...any) {
	z.sugar.Warnw(msg, args...)
}