/*
 *    Copyright 2018 Insolar
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

package preprocessor

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"text/template"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/log"
)

var clientFoundation = "github.com/insolar/insolar/toolkit/go/foundation"
var foundationPath = "github.com/insolar/insolar/logicrunner/goplugin/foundation"
var proxyctxPath = "github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
var corePath = "github.com/insolar/insolar/core"

type ParsedFile struct {
	name    string
	code    []byte
	fileSet *token.FileSet
	node    *ast.File

	types        map[string]*ast.TypeSpec
	methods      map[string][]*ast.FuncDecl
	constructors map[string][]*ast.FuncDecl
	contract     string
}

func (r *ParsedFile) codeOfNode(n ast.Node) string {
	return string(r.code[n.Pos()-1 : n.End()-1])
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

func ParseFile(fileName string) (*ParsedFile, error) {
	res := ParsedFile{
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

	err = getMethods(&res)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if res.contract == "" {
		return nil, errors.New("Only one smart contract must exist")
	}

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

func numberedVars(list *ast.FieldList, name string) string {
	if list == nil || list.NumFields() == 0 {
		return ""
	}

	rets := make([]string, list.NumFields())
	for i := range list.List {
		rets[i] = fmt.Sprintf("%s%d", name, i)
	}
	return strings.Join(rets, ", ")
}

func typeIndexes(parsed *ParsedFile, list *ast.FieldList, t string) []int {
	if list == nil || list.NumFields() == 0 {
		return []int{}
	}

	rets := []int{}
	for i, e := range list.List {
		if parsed.codeOfNode(e.Type) == t {
			rets = append(rets, i)
		}
	}
	return rets
}

func generateContractMethodsInfo(parsed *ParsedFile) ([]map[string]interface{}, []map[string]interface{}, map[string]bool) {
	imports := make(map[string]bool)
	imports[fmt.Sprintf(`"%s"`, proxyctxPath)] = true
	imports[fmt.Sprintf(`"%s"`, corePath)] = true

	var methodsInfo []map[string]interface{}
	for _, method := range parsed.methods[parsed.contract] {
		argsInit, argsList := generateZeroListOfTypes(parsed, "args", method.Type.Params)
		extendImportsMap(parsed, method.Type.Params, imports)

		info := map[string]interface{}{
			"Name":                method.Name.Name,
			"ArgumentsZeroList":   argsInit,
			"Arguments":           argsList,
			"Results":             numberedVars(method.Type.Results, "ret"),
			"ErrorInterfaceInRes": typeIndexes(parsed, method.Type.Results, "error"),
		}
		methodsInfo = append(methodsInfo, info)
	}

	var funcInfo []map[string]interface{}
	for _, f := range parsed.constructors[parsed.contract] {
		argsInit, argsList := generateZeroListOfTypes(parsed, "args", f.Type.Params)
		extendImportsMap(parsed, f.Type.Params, imports)

		info := map[string]interface{}{
			"Name":              f.Name.Name,
			"ArgumentsZeroList": argsInit,
			"Arguments":         argsList,
			"Results":           numberedVars(f.Type.Results, "ret"),
		}
		funcInfo = append(funcInfo, info)
	}

	return methodsInfo, funcInfo, imports
}

func GenerateContractWrapper(parsed *ParsedFile, out io.Writer) error {
	packageName := parsed.node.Name.Name

	tmpl, err := openTemplate("templates/wrapper.go.tpl")
	if err != nil {
		return errors.Wrap(err, "couldn't open template file for wrapper")
	}

	methodsInfo, funcInfo, imports := generateContractMethodsInfo(parsed)

	data := struct {
		PackageName    string
		ContractType   string
		Methods        []map[string]interface{}
		Functions      []map[string]interface{}
		ParsedCode     []byte
		FoundationPath string
		Imports        map[string]bool
	}{
		packageName,
		parsed.contract,
		methodsInfo,
		funcInfo,
		parsed.code,
		foundationPath,
		imports,
	}
	err = tmpl.Execute(out, data)
	if err != nil {
		return errors.Wrap(err, "couldn't write code output handle")
	}

	return nil
}

func ProxyPackageName(parsed *ParsedFile) (string, error) {
	match := regexp.MustCompile("([^/]+)/([^/]+).(go|insgoc)$").FindStringSubmatch(parsed.name)
	if match == nil {
		return "", errors.New("couldn't match filename without extension and path")
	}

	packageName := parsed.node.Name.Name

	proxyPackageName := packageName
	if proxyPackageName == "main" {
		proxyPackageName = match[2]
	}
	if proxyPackageName == "main" {
		proxyPackageName = match[1]
	}
	return proxyPackageName, nil
}

func GenerateContractProxy(parsed *ParsedFile, classReference string, out io.Writer) error {

	proxyPackageName, err := ProxyPackageName(parsed)
	if err != nil {
		return err
	}

	types := generateTypes(parsed)

	methodsProxies := generateMethodsProxies(parsed)

	constructorProxies := generateConstructorProxies(parsed)

	imports := generateImports(parsed)

	tmpl, err := openTemplate("templates/proxy.go.tpl")
	if err != nil {
		return errors.Wrap(err, "couldn't open template file for proxy")
	}

	data := struct {
		PackageName         string
		Types               []string
		ContractType        string
		MethodsProxies      []map[string]interface{}
		ConstructorsProxies []map[string]string
		ClassReference      string
		Imports             map[string]bool
	}{
		proxyPackageName,
		types,
		parsed.contract,
		methodsProxies,
		constructorProxies,
		classReference,
		imports,
	}
	err = tmpl.Execute(out, data)
	if err != nil {
		return errors.Wrap(err, "couldn't write code output handle")
	}

	return nil
}

func typeName(t ast.Expr) string {
	if tmp, ok := t.(*ast.StarExpr); ok { // *type
		t = tmp.X
	}
	return t.(*ast.Ident).Name
}

func IsContract(typeNode *ast.TypeSpec) bool {
	baseContract := "foundation.BaseContract"
	switch st := typeNode.Type.(type) {
	case *ast.StructType:
		if st.Fields == nil {
			return false
		}
		for _, fd := range st.Fields.List {
			selectField, ok := fd.Type.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			pack := selectField.X.(*ast.Ident).Name
			class := selectField.Sel.Name
			if (baseContract == (pack + "." + class)) && len(fd.Names) == 0 {
				return true
			}
		}
	}

	return false
}

func getMethods(parsed *ParsedFile) error {
	parsed.types = make(map[string]*ast.TypeSpec)
	parsed.methods = make(map[string][]*ast.FuncDecl)
	parsed.constructors = make(map[string][]*ast.FuncDecl)
	for _, d := range parsed.node.Decls {
		switch td := d.(type) {
		case *ast.GenDecl:
			if td.Tok != token.TYPE {
				continue
			}

			for _, e := range td.Specs {
				typeNode := e.(*ast.TypeSpec)

				if IsContract(typeNode) {
					if parsed.contract != "" {
						return errors.New("more than one contract in a file")
					}
					parsed.contract = typeNode.Name.Name
				} else {
					parsed.types[typeNode.Name.Name] = typeNode
				}
			}
		case *ast.FuncDecl:
			if td.Recv == nil || td.Recv.NumFields() == 0 {
				if !strings.HasPrefix(td.Name.Name, "New") {
					continue // doesn't look like a constructor
				}

				if td.Type.Results.NumFields() < 1 {
					log.Infof("Ignored %q as constructor, not enought returned values", td.Name.Name)
					continue
				}

				if td.Type.Results.NumFields() > 1 {
					log.Errorf("Constructor %q returns more than one argument, not supported at the moment", td.Name.Name)
					continue
				}

				typename := typeName(td.Type.Results.List[0].Type)
				parsed.constructors[typename] = append(parsed.constructors[typename], td)
			} else {
				typename := typeName(td.Recv.List[0].Type)
				parsed.methods[typename] = append(parsed.methods[typename], td)
			}
		}
	}

	return nil
}

// nolint
func generateTypes(parsed *ParsedFile) []string {
	var types []string
	for _, t := range parsed.types {
		types = append(types, "type "+parsed.codeOfNode(t))
	}

	return types
}

func extendImportsMap(parsed *ParsedFile, params *ast.FieldList, imports map[string]bool) {
	if params == nil || params.NumFields() == 0 {
		return
	}

	for _, e := range params.List {
		if parsed.codeOfNode(e.Type) == "error" {
			imports[fmt.Sprintf(`"%s"`, foundationPath)] = true
		}
	}

	for _, e := range params.List {
		tname := parsed.codeOfNode(e.Type)
		tname = strings.Trim(tname, "*")
		tnameFrom := strings.Split(tname, ".")

		if len(tnameFrom) < 2 {
			continue
		}

		for _, imp := range parsed.node.Imports {
			var importAlias string
			var impValue string

			if imp.Name != nil {
				importAlias = imp.Name.Name
				impValue = fmt.Sprintf(`%s %s`, importAlias, imp.Path.Value)
			} else {
				impValue = imp.Path.Value
				importString := strings.Trim(impValue, `"`)
				importAlias = filepath.Base(importString)
			}

			if importAlias == tnameFrom[0] {
				imports[impValue] = true
				break
			}
		}
	}
}

func generateZeroListOfTypes(parsed *ParsedFile, name string, list *ast.FieldList) (string, string) {
	text := fmt.Sprintf("%s := [%d]interface{}{}\n", name, list.NumFields())

	if list == nil {
		return text, ""
	}

	for i, arg := range list.List {
		tname := parsed.codeOfNode(arg.Type)
		if tname == "error" {
			tname = "*foundation.Error"
		}

		text += fmt.Sprintf("\tvar a%d %s\n", i, tname)
		text += fmt.Sprintf("\t%s[%d] = a%d\n", name, i, i)
	}

	listCode := ""
	for i, arg := range list.List {
		if i > 0 {
			listCode += ", "
		}
		listCode += fmt.Sprintf("%s[%d].(%s)", name, i, parsed.codeOfNode(arg.Type))
	}

	return text, listCode
}

func genFieldList(parsed *ParsedFile, params *ast.FieldList, withNames bool) string {
	res := ""
	if params == nil {
		return res
	}
	for i, e := range params.List {
		if i > 0 {
			res += ", "
		}
		if withNames {
			res += e.Names[0].Name + " "
		}
		res += parsed.codeOfNode(e.Type)
	}
	return res
}

func generateInitArguments(list *ast.FieldList) string {
	initArgs := ""
	initArgs += fmt.Sprintf("var args [%d]interface{}\n", list.NumFields())
	for i, arg := range list.List {
		initArgs += fmt.Sprintf("\targs[%d] = %s\n", i, arg.Names[0].Name)
	}
	return initArgs
}

func generateMethodProxyInfo(parsed *ParsedFile, method *ast.FuncDecl) map[string]interface{} {

	resInit, resList := generateZeroListOfTypes(parsed, "resList", method.Type.Results)

	return map[string]interface{}{
		"Name":           method.Name.Name,
		"ResultZeroList": resInit,
		"Results":        resList,
		"Arguments":      genFieldList(parsed, method.Type.Params, true),
		"ResultsTypes":   genFieldList(parsed, method.Type.Results, false),
		"InitArgs":       generateInitArguments(method.Type.Params),
	}
}

func generateMethodsProxies(parsed *ParsedFile) []map[string]interface{} {
	var methodsProxies []map[string]interface{}

	for _, method := range parsed.methods[parsed.contract] {
		if method.Name.IsExported() {
			methodsProxies = append(methodsProxies, generateMethodProxyInfo(parsed, method))
		}
	}

	return methodsProxies
}

func generateImports(parsed *ParsedFile) map[string]bool {
	imports := make(map[string]bool)
	imports[fmt.Sprintf(`"%s"`, proxyctxPath)] = true
	imports[fmt.Sprintf(`"%s"`, corePath)] = true
	for _, method := range parsed.methods[parsed.contract] {
		extendImportsMap(parsed, method.Type.Params, imports)
		extendImportsMap(parsed, method.Type.Results, imports)
	}
	for _, fun := range parsed.constructors[parsed.contract] {
		extendImportsMap(parsed, fun.Type.Params, imports)
		extendImportsMap(parsed, fun.Type.Results, imports)
	}
	return imports
}

func generateConstructorProxies(parsed *ParsedFile) []map[string]string {
	var res []map[string]string

	for _, e := range parsed.constructors[parsed.contract] {
		info := map[string]string{
			"Name":      e.Name.Name,
			"Arguments": genFieldList(parsed, e.Type.Params, true),
			"InitArgs":  generateInitArguments(e.Type.Params),
		}
		res = append(res, info)
	}
	return res
}

func CmdRewriteImports(parsed *ParsedFile, w io.Writer) error {
	if err := rewriteImports(parsed); err != nil {
		return errors.Wrap(err, "couldn't process")
	}
	if err := printer.Fprint(w, parsed.fileSet, parsed.node); err != nil {
		return errors.Wrap(err, "couldn't save")
	}
	return nil
}

func rewriteImports(p *ParsedFile) error {
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

func GetContractName(p *ParsedFile) string {
	return p.node.Name.Name
}

func RewriteContractPackage(p *ParsedFile, w io.Writer) {
	p.node.Name.Name = "main"
	printer.Fprint(w, p.fileSet, p.node)
}

// GetRealGenesisDir return dir under genesis dir
func GetRealGenesisDir(dir string) (string, error) {
	gopath := build.Default.GOPATH
	if gopath == "" {
		return "", errors.Errorf("GOPATH is not set")
	}

	contractsPath := ""
	for _, p := range strings.Split(gopath, ":") {
		contractsPath = path.Join(p, "src/github.com/insolar/insolar/genesis/", dir)
		_, err := os.Stat(contractsPath)
		if err == nil {
			return contractsPath, nil
		}
	}
	return "", errors.Errorf("Not found github.com/insolar/insolar in GOPATH")
}

// GetRealContractsNames returns names of all real smart contracts
func GetRealContractsNames() ([]string, error) {
	pathWithContracts, err := GetRealGenesisDir("experiment")
	if err != nil {
		return nil, errors.Wrap(err, "[ GetContractNames ]")
	}
	if len(pathWithContracts) == 0 {
		return nil, errors.New("[ GetContractNames ] There are contracts dir")
	}
	var result []string
	files, err := ioutil.ReadDir(pathWithContracts)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.IsDir() {
			result = append(result, f.Name())
		}
	}

	return result, nil
}
