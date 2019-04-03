package log

import (
	"context"

	"github.com/andreyromancev/belt"
)

type NoLogger struct{}

func (NoLogger) Debug(...interface{})                     {}
func (NoLogger) Debugf(string, ...interface{})            {}
func (NoLogger) Info(...interface{})                      {}
func (NoLogger) Infof(string, ...interface{})             {}
func (NoLogger) Warn(...interface{})                      {}
func (NoLogger) Warnf(string, ...interface{})             {}
func (NoLogger) Error(...interface{})                     {}
func (NoLogger) Errorf(string, ...interface{})            {}
func (NoLogger) Fatal(...interface{})                     {}
func (NoLogger) Fatalf(string, ...interface{})            {}
func (NoLogger) Panic(...interface{})                     {}
func (NoLogger) Panicf(string, ...interface{})            {}
func (NoLogger) WithFields(map[string]string) belt.Logger { return NoLogger{} }
func (NoLogger) WithField(string, string) belt.Logger     { return NoLogger{} }

type loggerKey struct{}

func FromContext(ctx context.Context) belt.Logger {
	if log, ok := ctx.Value(loggerKey{}).(belt.Logger); ok {
		return log
	}

	return NoLogger{}
}

func WithLogger(ctx context.Context, log belt.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, log)
}
