/*
 *    Copyright 2018 INS Ecosystem
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

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"

	"strconv"
)

var clientFoundation = "github.com/insolar/insolar/toolkit/go/foundation"
var foundationPath = "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type parsedFile struct {
	name    string
	code    []byte
	fileSet *token.FileSet
	node    *ast.File

	types    map[string]*ast.TypeSpec
	methods  map[string][]*ast.FuncDecl
	contract string
}

func printUsage() {
	fmt.Println("usage: preprocessor <command> [<args>]")
	fmt.Println("Commands: ")
	fmt.Println(" wrapper   generate contract's wrapper")
	fmt.Println(" proxy     generate contract's proxy")
	fmt.Println(" imports   rewrite imports")
}

type outputFlag struct {
	path   string
	writer io.Writer
}

func newOutputFlag() *outputFlag {
	return &outputFlag{path: "-", writer: os.Stdout}
}

func (r *outputFlag) String() string {
	return r.path
}
func (r *outputFlag) Set(arg string) error {
	var res io.Writer
	if arg == "-" {
		res = os.Stdout
	} else {
		var err error
		res, err = os.OpenFile(arg, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return errors.Wrap(err, "couldn't open file for writing")
		}
	}
	r.path = arg
	r.writer = res
	return nil
}
func (r *outputFlag) Type() string {
	return "file"
}

func main() {

	if len(os.Args) == 1 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "wrapper":
		fs := flag.NewFlagSet("wrapper", flag.ExitOnError)
		output := newOutputFlag()
		fs.VarP(output, "output", "o", "output file (use - for STDOUT)")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}

		for _, fn := range fs.Args() {
			err := generateContractWrapper(fn, output.writer)
			if err != nil {
				panic(err)
			}
		}
	case "proxy":
		fs := flag.NewFlagSet("proxy", flag.ExitOnError)
		output := newOutputFlag()
		fs.VarP(output, "output", "o", "output file (use - for STDOUT)")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}

		if fs.NArg() != 1 {
			panic(errors.New("proxy command should be followed by exactly one file name to process"))
		}

		err = generateContractProxy(fs.Arg(0), output.writer)
		if err != nil {
			panic(err)
		}
	case "imports":
		fs := flag.NewFlagSet("imports", flag.ExitOnError)
		output := newOutputFlag()
		fs.VarP(output, "output", "o", "output file (use - for STDOUT)")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}

		if fs.NArg() != 1 {
			panic(errors.New("imports command should be followed by exactly one file name to process"))
		}

		err = cmdRewriteImports(fs.Arg(0), output.writer)
		if err != nil {
			panic(err)
		}
	default:
		printUsage()
		fmt.Printf("\n\n%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
}

func slurpFile(fileName string) ([]byte, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		return nil, errors.Wrap(err, "Can't open file '"+fileName+"'")
	}
	defer file.Close() //nolint: errcheck

	res, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "Can't read file '"+fileName+"'")
	}
	return res, nil
}

func parseFile(fileName string) (*parsedFile, error) {
	res := parsedFile{
		name: fileName,
	}
	sourceCode, err := slurpFile(fileName)
	if err != nil {
		return &res, errors.Wrap(err, "Can't read slurp file")
	}
	res.code = sourceCode

	res.fileSet = token.NewFileSet()
	node, err := parser.ParseFile(res.fileSet, res.name, res.code, parser.ParseComments)
	if err != nil {
		return &res, errors.Wrapf(err, "Can't parse %s", fileName)
	}
	res.node = node

	getMethods(&res)

	return &res, nil
}

func openTemplate(fileName string) (*template.Template, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, errors.Wrap(nil, "couldn't find info about current file")
	}
	templateDir := filepath.Join(filepath.Dir(currentFile), fileName)
	tmpl, err := template.ParseFiles(templateDir)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't parse template for output")
	}
	return tmpl, nil
}

func generateContractMethodsInfo(parsed *parsedFile) []map[string]interface{} {
	var methodsInfo []map[string]interface{}
	for _, method := range parsed.methods[parsed.contract] {
		argsInit, argsList := generateZeroListOfTypes(parsed, "args", method.Type.Params)

		rets := []string{}
		for i := range method.Type.Results.List {
			rets = append(rets, fmt.Sprintf("ret%d", i))
		}
		resultList := strings.Join(rets, ", ")

		info := map[string]interface{}{
			"Name":              method.Name.Name,
			"ArgumentsZeroList": argsInit,
			"Results":           resultList,
			"Arguments":         argsList,
		}
		methodsInfo = append(methodsInfo, info)
	}
	return methodsInfo
}

func generateContractWrapper(fileName string, out io.Writer) error {
	parsed, err := parseFile(fileName)
	if err != nil {
		return errors.Wrap(err, "couldn't parse")
	}

	packageName := parsed.node.Name.Name
	if packageName != "main" {
		panic("Contract must be in main package")
	}

	tmpl, err := openTemplate("templates/wrapper.go.tpl")
	if err != nil {
		return errors.Wrap(err, "couldn't open template file for wrapper")
	}

	data := struct {
		PackageName    string
		ContractType   string
		Methods        []map[string]interface{}
		ParsedCode     []byte
		FoundationPath string
	}{
		packageName,
		parsed.contract,
		generateContractMethodsInfo(parsed),
		parsed.code,
		foundationPath,
	}
	err = tmpl.Execute(out, data)
	if err != nil {
		return errors.Wrap(err, "couldn't write code output handle")
	}

	return nil
}

func generateContractProxy(fileName string, out io.Writer) error {
	parsed, err := parseFile(fileName)
	if err != nil {
		return errors.Wrap(err, "couldn't parse")
	}

	match := regexp.MustCompile("([^/]+)/([^/]+).go$").FindStringSubmatch(fileName)
	if match == nil {
		return errors.Wrap(err, "couldn't match filename without extension and path")
	}

	packageName := parsed.node.Name.Name
	if packageName != "main" {
		panic("Contract must be in main package")
	}

	proxyPackageName := match[2]
	if proxyPackageName == "main" {
		proxyPackageName = match[1]
	}

	types := generateTypes(parsed)

	methodsProxies := generateMethodsProxies(parsed)

	tmpl, err := openTemplate("templates/proxy.go.tpl")
	if err != nil {
		return errors.Wrap(err, "couldn't open template file for proxy")
	}

	data := struct {
		PackageName    string
		Types          []string
		ContractType   string
		MethodsProxies []map[string]interface{}
	}{
		proxyPackageName,
		types,
		parsed.contract,
		methodsProxies,
	}
	err = tmpl.Execute(out, data)
	// _, err = out.Write([]byte(code))
	if err != nil {
		return errors.Wrap(err, "couldn't write code output handle")
	}

	return nil
}

func getMethods(parsed *parsedFile) {
	parsed.types = make(map[string]*ast.TypeSpec)
	parsed.methods = make(map[string][]*ast.FuncDecl)
	for _, d := range parsed.node.Decls {
		switch td := d.(type) {
		case *ast.GenDecl:
			if td.Tok != token.TYPE {
				continue
			}

			for _, e := range td.Specs {
				typeNode := e.(*ast.TypeSpec)

				if strings.Contains(td.Doc.Text(), "@inscontract") {
					if parsed.contract != "" {
						panic("more than one contract in a file")
					}
					parsed.contract = typeNode.Name.Name
				} else {
					parsed.types[typeNode.Name.Name] = typeNode
				}
			}
		case *ast.FuncDecl:
			if td.Recv == nil || td.Recv.NumFields() == 0 {
				continue // todo we must store it and use, it may be a constructor
			}

			r := td.Recv.List[0].Type
			if tr, ok := r.(*ast.StarExpr); ok { // *type
				r = tr.X
			}
			typename := r.(*ast.Ident).Name
			parsed.methods[typename] = append(parsed.methods[typename], td)
		}
	}
}

// nolint
func generateTypes(parsed *parsedFile) []string {
	var types []string
	for _, t := range parsed.types {
		types = append(types, "type "+string(parsed.code[t.Pos()-1:t.End()-1]))
	}

	return types
}

func generateZeroListOfTypes(parsed *parsedFile, name string, list *ast.FieldList) (string, string) {
	text := fmt.Sprintf("%s := [%d]interface{}{}\n", name, list.NumFields())

	if list == nil {
		return text, ""
	}

	for i, arg := range list.List {
		initializer := ""
		tname := string(parsed.code[arg.Type.Pos()-1 : arg.Type.End()-1])
		switch tname {
		case "uint", "int", "int8", "uint8", "int32", "uint32", "int64", "uint64":
			initializer = tname + "(0)"
		case "string":
			initializer = `""`
		default:
			switch td := arg.Type.(type) {
			case *ast.StarExpr:
				initializer = "&" + string(parsed.code[td.X.Pos()-1:td.X.End()-1]) + "{}"
			default:
				initializer = tname + "{}"
			}
		}
		text += fmt.Sprintf("\t%s[%d] = %s\n", name, i, initializer)
	}

	listCode := ""
	for i, arg := range list.List {
		if i > 0 {
			listCode += ", "
		}
		listCode += fmt.Sprintf("%s[%d].(%s)", name, i, string(parsed.code[arg.Type.Pos()-1:arg.Type.End()-1]))
	}

	return text, listCode
}

func generateArguments(parsed *parsedFile, params *ast.FieldList) string {
	args := ""
	for i, arg := range params.List {
		if i > 0 {
			args += ", "
		}
		args += arg.Names[0].Name
		args += " " + string(parsed.code[arg.Type.Pos()-1:arg.Type.End()-1])
	}
	return args
}

func generateResultsTypes(parsed *parsedFile, results *ast.FieldList) string {
	resultsTypes := ""
	if results.NumFields() > 0 {
		for i, arg := range results.List {
			if i > 0 {
				resultsTypes += ", "
			}
			resultsTypes += string(parsed.code[arg.Type.Pos()-1 : arg.Type.End()-1])
		}
	}
	return resultsTypes
}

func generateInitArguments(list *ast.FieldList) string {
	initArgs := ""
	initArgs += fmt.Sprintf("var args [%d]interface{}\n", list.NumFields())
	for i, arg := range list.List {
		initArgs += fmt.Sprintf("\targs[%d] = %s\n", i, arg.Names[0].Name)
	}
	return initArgs
}

func generateMethodProxyInfo(parsed *parsedFile, method *ast.FuncDecl) map[string]interface{} {

	args := generateArguments(parsed, method.Type.Params)

	resultsTypes := generateResultsTypes(parsed, method.Type.Results)

	initArgs := generateInitArguments(method.Type.Params)

	resInit, resList := generateZeroListOfTypes(parsed, "resList", method.Type.Results)

	info := map[string]interface{}{
		"Name":           method.Name.Name,
		"ResultZeroList": resInit,
		"Results":        resList,
		"Arguments":      args,
		"ResultsTypes":   resultsTypes,
		"InitArgs":       initArgs,
	}
	return info
}

func generateMethodsProxies(parsed *parsedFile) []map[string]interface{} {
	var methodsProxies []map[string]interface{}

	for _, method := range parsed.methods[parsed.contract] {
		methodsProxies = append(methodsProxies, generateMethodProxyInfo(parsed, method))
	}
	return methodsProxies
}

func cmdRewriteImports(fname string, w io.Writer) error {
	parsed, err := parseFile(fname)
	if err != nil {
		return errors.Wrap(err, "couldn't parse")
	}
	if err := rewriteImports(parsed); err != nil {
		return errors.Wrap(err, "couldn't process")
	}
	if err := printer.Fprint(w, parsed.fileSet, parsed.node); err != nil {
		return errors.Wrap(err, "couldn't save")
	}
	return nil
}

func rewriteImports(p *parsedFile) error {
	quoted := strconv.Quote(clientFoundation)
	for _, d := range p.node.Decls {
		td, ok := d.(*ast.GenDecl)
		if !ok {
			continue
		}
		if td.Tok != token.IMPORT {
			continue
		}
		for _, s := range td.Specs {
			is, ok := s.(*ast.ImportSpec)
			if !ok {
				continue
			}
			if is.Path.Value == quoted {
				is.Path = &ast.BasicLit{Value: strconv.Quote(foundationPath)}
			}
		}
	}
	return nil
}
