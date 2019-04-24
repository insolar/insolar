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
	"go/ast"
	"go/parser"
	"go/token"
	"path"
	"strings"
	"text/template"
)

var (
	templateFuncs = template.FuncMap{
		"getAdapterHelper": func(m stateMachineWithID, helperType *string) string {
			if helperType == nil {
				return ""
			}
			return ", helpers." + m.AdapterHelperCatalog[*helperType]
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

func parseSource(file string) (string, map[string]string) {
	imports := make(map[string]string)
	set := token.NewFileSet()
	node, err := parser.ParseFile(set, file, nil, parser.ParseComments)
	checkErr(err)
	for _, decls := range node.Decls {
		if decl, ok := decls.(*ast.GenDecl); ok {
			if decl.Tok != token.IMPORT {
				continue
			}
			for _, spec := range decl.Specs {
				importPath := spec.(*ast.ImportSpec).Path.Value
				if name := spec.(*ast.ImportSpec).Name; name != nil {
					exitWithError("Import aliases not allowed <%s %s>", name, importPath)
				}
				_, name := path.Split(importPath[1 : len(importPath)-1])
				imports[name] = importPath
			}
		}
	}
	return node.Name.Name, imports
}
