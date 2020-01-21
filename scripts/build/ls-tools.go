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

//usr/bin/env go run -mod=vendor "$0" "$@"; exit "$?"
// ls-tools.go - go script extract tools imports from tools.go

// +build tools

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	f, err := parser.ParseFile(token.NewFileSet(), "tools.go", nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}
	v := &visitor{}
	ast.Walk(v, f)
	for _, imp := range v.imports {
		fmt.Println(imp)
	}
}

type visitor struct {
	imports []string
}

func (v *visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	if imp, ok := n.(*ast.ImportSpec); ok {
		if imp.Name.String() == "_" {
			v.imports = append(v.imports, imp.Path.Value)
		}
	}
	return v
}
