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
