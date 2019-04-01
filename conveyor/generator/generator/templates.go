/*
 *    Copyright 2019 Insolar Technologies
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

package generator

import (
	"go/parser"
	"go/token"
	"path"
	"strings"
	"text/template"
)

func getDirFromPath(p string) string {
	dir, _ := path.Split(p)
	if strings.HasSuffix(dir, "/") {
		return dir[:len(dir)-1]
	}
	return dir
}

var (
	templateFuncs = template.FuncMap{
		"fileToImport": func(f string) string {
			if idx := strings.Index(f, "github.com/insolar/insolar"); idx >= 0 {
				return getDirFromPath(f[idx:])
			}
			return getDirFromPath(f)
		},
		"unPackage": func(t string, p string) string {
			if idx := strings.Index(t, p); idx == 0 || (idx == 1 && t[0] == '*') {
				return strings.Replace(t, p+".", "", 1)
			}
			return t
		},
		"isNull": func(i int) bool {
			return i == 0
		},
		"handlerExists": func(x *handler) bool {
			return x != nil
		},
	}
)

func getPackage(file string) string {
	set := token.NewFileSet()
	node, err := parser.ParseFile(set, file, nil, parser.ParseComments)
	checkErr(err)
	return node.Name.Name
}
