//
// Copyright 2019 Insolar Technologies GmbH
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
//

package log

import (
	"runtime"
	"strings"
)

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

func stripPackageName(packageName string) string {
	result := strings.TrimPrefix(packageName, insolarPrefix)
	i := strings.Index(result, ".")
	if result == packageName || i == -1 {
		return result
	}
	return result[:i]
}
