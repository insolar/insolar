/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package log

import (
	"path"
	"runtime"
	"strings"
)

// callInfo bundles the info about the call environment
// when a logging statement occurred.
type callInfo struct {
	packageName string
	fileName    string
	funcName    string
	line        int
}

func getCallInfo(skipCallNumber int) *callInfo {
	pc, file, line, _ := runtime.Caller(skipCallNumber)
	_, fileName := path.Split(file)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), ".")
	pl := len(parts)
	packageName := ""
	funcName := parts[pl-1]

	if pl > 1 && parts[pl-2][0] == '(' {
		funcName = parts[pl-2] + "." + funcName
		packageName = strings.Join(parts[0:pl-2], ".")
	} else {
		packageName = strings.Join(parts[0:pl-1], ".")
	}

	return &callInfo{
		packageName: stripPackageName(packageName),
		fileName:    fileName,
		funcName:    funcName,
		line:        line,
	}
}

func stripPackageName(packageName string) string {
	result := strings.TrimPrefix(packageName, "github.com/insolar/insolar/")
	i := strings.Index(result, ".")
	if result == packageName || i == -1 {
		return result
	}
	return result[:i]
}
