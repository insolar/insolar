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

	"io/ioutil"

	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

func init() {
	flag.Parse()
}

func main() {
	log.Println(os.Getwd())
	for _, fn := range flag.Args() {
		w, err := generateForFile(fn)
		if err != nil {
			panic(err)
		}
		_, _ = io.Copy(os.Stdout, w)
	}
}

func generateForFile(fn string) (io.Reader, error) {
	fs := token.NewFileSet()

	F, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		return nil, errors.Wrap(err, "Can't open file "+fn)
	}
	defer F.Close()

	buff, err := ioutil.ReadAll(F)
	if err != nil {
		return nil, errors.Wrap(err, "Can't read file "+fn)
	}

	node, err := parser.ParseFile(fs, fn, buff, parser.ParseComments)
	if err != nil {
		log.Fatalf("Can't parse %s : %s", fn, err)
	}
	if node.Name.Name != "main" {
		panic("Contract must be in main package")
	}
	getMethods(node, buff)
	code := generateWrappers()
	code += "\n" + generateExports() + "\n"
	return bytes.NewBuffer([]byte(code)), nil
}

var types = make(map[string]string)
var methods = make(map[string][]*ast.FuncDecl)
var contract string

func getMethods(F *ast.File, text []byte) {
	for _, d := range F.Decls {
		switch td := d.(type) {
		case *ast.GenDecl:
			if td.Tok != token.TYPE {
				continue
			}
			typeNode := td.Specs[0].(*ast.TypeSpec)
			if strings.Contains(td.Doc.Text(), "@inscontract") {
				if contract != "" {
					panic("more than one contract in a file")
				}
				contract = typeNode.Name.Name
				continue
			}
			types[typeNode.Name.Name] = string(text[typeNode.Pos()-1 : typeNode.End()])
			continue
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
	for _, t := range types {
		text += "type " + t + "\n"
	}
	for _, method := range methods[contract] {
		text += generateMethodWrapper(method, contract) + "\n\n"
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
	text += "var INSEXPORT " + contract
	return text
}
