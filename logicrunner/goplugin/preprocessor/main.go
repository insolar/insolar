package main

import (
	"go/parser"
	"go/token"
	"log"

	"os"

	"go/ast"

	"strings"

	flag "github.com/spf13/pflag"
)

func init() {
	flag.Parse()
}

func main() {
	log.Println(os.Getwd())
	for i, fn := range flag.Args() {
		log.Println(i, fn)
		fs := token.NewFileSet()
		node, err := parser.ParseFile(fs, fn, nil, parser.ParseComments)
		if err != nil {
			log.Fatalf("Can't parse %s : %s", fn, err)
		}
		//log.Print(node)
		GenerateWrappers(node)
		log.Printf("%+v", methods)
	}
}

var methods = make(map[string][]*ast.FuncDecl)
var contracts = make([]string, 1)

func GenerateWrappers(F *ast.File) {
	for _, d := range F.Decls {
		switch td := d.(type) {
		case *ast.GenDecl:
			if td.Tok != token.TYPE {
				continue
			}
			if !strings.Contains(td.Doc.Text(), "@contract") {
				continue
			}
			typename := td.Specs[0].(*ast.TypeSpec).Name.Name
			contracts = append(contracts, typename)

		case *ast.FuncDecl:
			if td.Recv.NumFields() == 0 { // not a method
				continue
			}
			r := td.Recv.List[0].Type
			if tr, ok := r.(*ast.StarExpr); ok { // *type
				r = tr.X
			}
			typename := r.(*ast.Ident).Name
			methods[typename] = append(methods[typename], td)
		}
	}
}
