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

// ls-imports - go tool which extract go imports from provided file.

package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

var (
	underscoresOnly = flag.Bool("u", false, "print only underscore imports")
	parseFile       = flag.String("f", "tools.go", "path to *.go file with imports to extract")
)

func main() {
	flag.Parse()
	f, err := parser.ParseFile(token.NewFileSet(), *parseFile, nil, parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}
	v := &visitor{}
	ast.Walk(v, f)
	for _, imp := range v.imports {
		if *underscoresOnly && imp.Name != "_" {
			continue
		}
		fmt.Println(imp.Value)
	}
}

type importItem struct {
	Name  string
	Value string
}

type visitor struct {
	imports []importItem
}

func (v *visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}

	if imp, ok := n.(*ast.ImportSpec); ok {
		v.imports = append(v.imports, importItem{
			Name:  imp.Name.String(),
			Value: imp.Path.Value,
		})
	}
	return v
}
