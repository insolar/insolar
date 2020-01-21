// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
