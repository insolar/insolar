package log

import (
	"fmt"
	"strings"
	"sync"

	"github.com/andreyromancev/belt"
)

type ConsoleLogger struct {
	lock   sync.RWMutex
	fields map[string]string
}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{fields: make(map[string]string)}
}

func (l *ConsoleLogger) Debug(a ...interface{}) {
	l.withLevel("DEBUG", a...)
}

func (l *ConsoleLogger) Debugf(s string, a ...interface{}) {
	l.Debug(fmt.Sprintf(s, a...))
}

func (l *ConsoleLogger) Info(a ...interface{}) {
	l.withLevel("INFO", a...)
}

func (l *ConsoleLogger) Infof(s string, a ...interface{}) {
	l.Info(fmt.Sprintf(s, a...))
}

func (l *ConsoleLogger) Warn(a ...interface{}) {
	l.withLevel("WARN", a...)
}

func (l *ConsoleLogger) Warnf(s string, a ...interface{}) {
	l.Warn(fmt.Sprintf(s, a...))
}

func (l *ConsoleLogger) Error(a ...interface{}) {
	l.withLevel("ERROR", a...)
}

func (l *ConsoleLogger) Errorf(s string, a ...interface{}) {
	l.Error(fmt.Sprintf(s, a...))
}

func (l *ConsoleLogger) Fatal(a ...interface{}) {
	l.withLevel("FATAL", a...)
}

func (l *ConsoleLogger) Fatalf(s string, a ...interface{}) {
	l.Fatal(fmt.Sprintf(s, a...))
}

func (l *ConsoleLogger) Panic(a ...interface{}) {
	l.withLevel("PANIC", a...)
}

func (l *ConsoleLogger) Panicf(s string, a ...interface{}) {
	l.Panic(fmt.Sprintf(s, a...))
}

func (l *ConsoleLogger) WithFields(fields map[string]string) belt.Logger {
	l.lock.Lock()
	res := make(map[string]string, len(l.fields))
	for k, v := range l.fields {
		res[k] = v
	}
	l.lock.Unlock()
	for k, v := range fields {
		res[k] = v
	}
	return &ConsoleLogger{fields: res}
}

func (l *ConsoleLogger) WithField(k string, v string) belt.Logger {
	return l.WithFields(map[string]string{k: v})
}

func (l *ConsoleLogger) withLevel(level string, a ...interface{}) {
	fields := make([]string, 0, len(l.fields))
	for k, v := range l.fields {
		fields = append(fields, fmt.Sprintf("%s=%s", k, v))
	}
	a = append([]interface{}{fmt.Sprintf("[%s] ", level)}, a...)
	a = append(a, fmt.Sprintf("(%s)", strings.Join(fields, " ")))
	fmt.Println(a...)
}
