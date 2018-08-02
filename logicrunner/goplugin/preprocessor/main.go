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
	for _, fn := range flag.Args() {
		generateForFile(fn)
	}
}

func generateForFile(fn string) {
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, fn, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Can't parse %s : %s", fn, err)
	}
	//log.Print(node)
	getMethods(node)
	code := generateWrappers()
	code += "\n" + generateExports()
	log.Printf("%s", code)
}

var methods = make(map[string][]*ast.FuncDecl)
var contracts = make([]string, 1)

func getMethods(F *ast.File) {
	for _, d := range F.Decls {
		switch td := d.(type) {
		case *ast.GenDecl:
			if td.Tok != token.TYPE {
				continue
			}
			log.Printf(td.Doc.Text())
			if !strings.Contains(td.Doc.Text(), "@inscontract") {
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

func generateWrappers() string {
	text := ""
	log.Printf("%+v", contracts)
	for _, class := range contracts {
		for _, method := range methods[class] {
			text += generateMethodWrapper(method, class) + "\n\n"
		}
	}
	return text
}

func generateMethodWrapper(method *ast.FuncDecl, class string) string {
	text := ""
	text += "func __INSMETHOD__" + method.Name.Name + "(_self *" + class + ","
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
		text += "var __INSEXPORT" + m + " " + m
	}
	return text
}
