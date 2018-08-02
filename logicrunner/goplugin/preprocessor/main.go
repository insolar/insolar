package main

import (
	"go/parser"
	"go/token"
	"log"

	"os"

	"go/ast"

	"strings"

	"bytes"
	"io"

	flag "github.com/spf13/pflag"
)

func init() {
	flag.Parse()
}

func main() {
	log.Println(os.Getwd())
	for _, fn := range flag.Args() {
		w := generateForFile(fn)
		io.Copy(os.Stdout, w)
	}
}

func generateForFile(fn string) io.Reader {
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, fn, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Can't parse %s : %s", fn, err)
	}
	if node.Name.Name != "main" {
		panic("Contract must be in main package")
	}
	getMethods(node)
	code := generateWrappers()
	code += "\n" + generateExports() + "\n"
	return bytes.NewBuffer([]byte(code))
}

var methods = make(map[string][]*ast.FuncDecl)
var contracts []string

func getMethods(F *ast.File) {
	for _, d := range F.Decls {
		switch td := d.(type) {
		case *ast.GenDecl:
			if td.Tok != token.TYPE {
				continue
			}
			if !strings.Contains(td.Doc.Text(), "@inscontract") {
				continue
			}
			typename := td.Specs[0].(*ast.TypeSpec).Name.Name
			contracts = append(contracts, typename)
			if len(contracts) > 1 {
				panic("more than one contract in a file")
			}
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

func generateWrappers() string {
	text := ""
	for _, class := range contracts {
		for _, method := range methods[class] {
			text += generateMethodWrapper(method, class) + "\n\n"
		}
	}
	return text
}

func generateMethodWrapper(method *ast.FuncDecl, class string) string {
	text := ""
	text += "func (_self *" + class + ") INSMETHOD__" + method.Name.Name + "("
	for _, arg := range method.Type.Params.List {
		text += arg.Names[0].Name + " interface{}, "
	}
	text += ") {\n\t_self." + method.Name.Name + "("
	for _, arg := range method.Type.Params.List {
		text += arg.Names[0].Name + ".(" + arg.Type.(*ast.Ident).Name + "), "
	}
	text += ")\n}"
	return text
}

func generateExports() string {
	text := ""
	for _, m := range contracts {
		text += "var INSEXPORT " + m
	}
	return text
}
