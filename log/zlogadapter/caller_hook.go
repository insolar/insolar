// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package zlogadapter

import (
	"runtime"
	"strings"

	"github.com/rs/zerolog"
)

// FuncFieldName is the field name used for func field.
var FuncFieldName = "func"

type callerHook struct {
	callerSkipFrameCount int
}

func newCallerHook(skipFrameCount int) *callerHook {
	return &callerHook{callerSkipFrameCount: skipFrameCount}
}

// Run implements zerolog.Hook.
func (ch *callerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level == zerolog.NoLevel {
		return
	}
	info := getCallInfo(ch.callerSkipFrameCount)
	e.Str(zerolog.CallerFieldName, info.fileName)
	e.Str(FuncFieldName, info.funcName)
}

// callInfo bundles the info about the call environment
// when a logging statement occurred.
type callInfo struct {
	fileName string
	funcName string
	line     int
}

func getCallInfo(skipCallNumber int) *callInfo {
	pc, file, line, _ := runtime.Caller(skipCallNumber)

	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	funcName := parts[pl-1]

	if pl > 1 && strings.HasPrefix(parts[pl-2], "(") {
		funcName = parts[pl-2] + "." + funcName
	}

	return &callInfo{
		fileName: trimInsolarPrefix(file, line),
		funcName: funcName,
		line:     line,
	}
}
