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

package preprocessor

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
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
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/insolar/insolar/insolar/genesisrefs"

	"github.com/insolar/insolar/insolar"

	"github.com/pkg/errors"
)

var foundationPath = "github.com/insolar/insolar/logicrunner/builtin/foundation"
var proxyctxPath = "github.com/insolar/insolar/logicrunner/common"
var corePath = "github.com/insolar/insolar/insolar"

var immutableFlag = "ins:immutable"
var sagaFlagStart = "ins:saga("
var sagaFlagEnd = ")"
var sagaFlagStartLength = len(sagaFlagStart)

const (
	TemplateDirectory = "templates"

	mainPkg   = "main"
	errorType = "error"
)

// SagaInfo stores sagas-related information for given contract method.
// If a method is marked with //ins:saga(Rollback) SagaInfo stores
// `IsSaga: true, RollbackMethodName: "Rollback"`. Also Arguments always stores
// a comma separated list of arguments and NumArguments always stores the number
// of arguments.
type SagaInfo struct {
	IsSaga             bool
	RollbackMethodName string
	// we have to duplicate argument list here because it's not always used in templates
	Arguments    string
	NumArguments int
}

// ParsedFile struct with prepared info we extract from source code
type ParsedFile struct {
	name        string
	code        []byte
	fileSet     *token.FileSet
	node        *ast.File
	machineType insolar.MachineType

	types        map[string]*ast.TypeSpec
	methods      map[string][]*ast.FuncDecl
	constructors map[string][]*ast.FuncDecl
	contract     string
}

// ParseFile parses a file as Go source code of a smart contract
// and returns it as `ParsedFile`
func ParseFile(fileName string, machineType insolar.MachineType) (*ParsedFile, error) {
	res := &ParsedFile{
		name:        fileName,
		machineType: machineType,
	}
	sourceCode, err := slurpFile(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "Can't read file")
	}
	res.code = sourceCode

	res.fileSet = token.NewFileSet()
	node, err := parser.ParseFile(res.fileSet, res.name, res.code, parser.ParseComments)
	if err != nil {
		return nil, errors.Wrapf(err, "Can't parse %s", fileName)
	}
	res.node = node

	err = res.parseTypes()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	err = res.parseFunctionsAndMethods()
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if res.contract == "" {
		return nil, errors.New("Only one smart contract must exist")
	}

	return res, nil
}

func (pf *ParsedFile) parseTypes() error {
	pf.types = make(map[string]*ast.TypeSpec)
	for _, decl := range pf.node.Decls {
		tDecl, ok := decl.(*ast.GenDecl)
		if !ok || tDecl.Tok != token.TYPE {
			continue
		}

		for _, e := range tDecl.Specs {
			typeNode := e.(*ast.TypeSpec)

			err := pf.parseTypeSpec(typeNode)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (pf *ParsedFile) parseTypeSpec(typeSpec *ast.TypeSpec) error {
	if isContractTypeSpec(typeSpec) {
		if pf.contract != "" {
			return errors.New("more than one contract in a file")
		}
		pf.contract = typeSpec.Name.Name
	} else {
		pf.types[typeSpec.Name.Name] = typeSpec
	}

	return nil
}

func (pf *ParsedFile) parseFunctionsAndMethods() error {
	pf.methods = make(map[string][]*ast.FuncDecl)
	pf.constructors = make(map[string][]*ast.FuncDecl)
	for _, decl := range pf.node.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || !fd.Name.IsExported() {
			continue
		}

		var err error
		if fd.Recv == nil || fd.Recv.NumFields() == 0 {
			err = pf.parseConstructor(fd)
		} else {
			err = pf.parseMethod(fd)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (pf *ParsedFile) parseConstructor(fd *ast.FuncDecl) error {
	name := fd.Name.Name
	if !strings.HasPrefix(name, "New") {
		return nil // doesn't look like a constructor
	}

	res := fd.Type.Results

	if res.NumFields() != 2 {
		return errors.Errorf("Constructor %q should return exactly two values", name)
	}

	if pf.typeName(res.List[1].Type) != errorType {
		return errors.Errorf("Constructor %q should return 'error'", name)
	}

	typename := pf.typeName(res.List[0].Type)
	pf.constructors[typename] = append(pf.constructors[typename], fd)

	return nil
}

func (pf *ParsedFile) parseMethod(fd *ast.FuncDecl) error {
	name := fd.Name.Name

	res := fd.Type.Results
	if res.NumFields() < 1 {
		return errors.Errorf("Method %q should return at least one result (error)", name)
	}

	lastResType := pf.typeName(res.List[res.NumFields()-1].Type)
	if lastResType != errorType {
		return errors.Errorf(
			"Method %q should return 'error' as last value, but it's %q",
			name, lastResType,
		)
	}

	typename := pf.typeName(fd.Recv.List[0].Type)
	pf.methods[typename] = append(pf.methods[typename], fd)

	return nil
}

// ProxyPackageName guesses user friendly contract "name" from file name
// and/or package in the file
func (pf *ParsedFile) ProxyPackageName() (string, error) {
	match := regexp.MustCompile("([^/]+)/([^/]+).(go|insgoc)$").FindStringSubmatch(pf.name)
	if match == nil {
		return "", errors.New("couldn't match filename without extension and path")
	}

	packageName := pf.node.Name.Name

	proxyPackageName := packageName
	if proxyPackageName == mainPkg {
		proxyPackageName = match[2]
	}
	if proxyPackageName == mainPkg {
		proxyPackageName = match[1]
	}
	return proxyPackageName, nil
}

// ContractName returns name of the contract
func (pf *ParsedFile) ContractName() string {
	return pf.node.Name.Name
}

func checkMachineType(machineType insolar.MachineType) error {
	if machineType != insolar.MachineTypeGoPlugin &&
		machineType != insolar.MachineTypeBuiltin {
		return errors.New("Unsupported machine type")
	}
	return nil
}

func templatePathConstruct(tplType string) string {
	return path.Join(TemplateDirectory, tplType+".go.tpl")
}

func formatAndWrite(out io.Writer, templateName string, data map[string]interface{}) error {
	templatePath := templatePathConstruct(templateName)
	tmpl, err := openTemplate(templatePath)
	if err != nil {
		return errors.Wrap(err, "couldn't open template file for wrapper")
	}

	var buff bytes.Buffer

	err = tmpl.Execute(&buff, data)
	if err != nil {
		return errors.Wrap(err, "couldn't write code output handle")
	}

	fmtOut, err := format.Source(buff.Bytes())
	if err != nil {

		return errors.Wrap(err, "couldn't format code "+buff.String())
	}

	_, err = out.Write(fmtOut)
	if err != nil {
		return errors.Wrap(err, "couldn't write code to output")
	}

	return nil
}

// WriteWrapper generates and writes into `out` source code
// of wrapper for the contract
func (pf *ParsedFile) WriteWrapper(out io.Writer, packageName string) error {
	if err := checkMachineType(pf.machineType); err != nil {
		return err
	}

	functionsInfo := pf.functionInfoForWrapper(pf.constructors[pf.contract])
	for _, fi := range functionsInfo {
		if fi["SagaInfo"].(*SagaInfo).IsSaga {
			return fmt.Errorf("semantic error: '%s' can't be a saga because it's a constructor", fi["Name"].(string))
		}
	}

	methodsInfo := pf.functionInfoForWrapper(pf.methods[pf.contract])
	err := pf.checkSagaRollbackMethodsExistAndMatch(methodsInfo)
	if err != nil {
		return err
	}

	err = pf.checkSagaIsNotImmutable(methodsInfo)
	if err != nil {
		return err
	}

	err = pf.checkSagaMethodsReturnOnlySingleErrorValue(methodsInfo)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		"Package":            packageName,
		"ContractType":       pf.contract,
		"Methods":            methodsInfo,
		"Functions":          functionsInfo,
		"ParsedCode":         pf.code,
		"FoundationPath":     foundationPath,
		"Imports":            pf.generateImports(true),
		"GenerateInitialize": pf.machineType == insolar.MachineTypeBuiltin,
	}

	return formatAndWrite(out, "wrapper", data)
}

func (pf *ParsedFile) checkSagaIsNotImmutable(methodsInfo []map[string]interface{}) error {
	for _, mi := range methodsInfo {
		sagaInfo := mi["SagaInfo"].(*SagaInfo)
		if !sagaInfo.IsSaga {
			continue
		}

		if mi["Immutable"].(bool) {
			return fmt.Errorf("semantic error: '%s' can't be a saga because it's immutable", mi["Name"].(string))
		}

		for _, ri := range methodsInfo {
			if ri["Name"].(string) != sagaInfo.RollbackMethodName {
				continue
			}

			if ri["Immutable"].(bool) {
				return fmt.Errorf("semantic error: '%s' can't be saga's rollback method because it's immutable", ri["Name"].(string))
			}
		}
	}

	return nil
}

func (pf *ParsedFile) checkSagaMethodsReturnOnlySingleErrorValue(methodsInfo []map[string]interface{}) error {
	returnsOnlyError := func(info map[string]interface{}) bool {
		return info["Results"].(string) == "ret0" &&
			len(info["ErrorInterfaceInRes"].([]int)) == 1 &&
			info["ErrorInterfaceInRes"].([]int)[0] == 0
	}
	for _, mi := range methodsInfo {
		sagaInfo := mi["SagaInfo"].(*SagaInfo)
		if !sagaInfo.IsSaga {
			continue
		}

		if !returnsOnlyError(mi) {
			return fmt.Errorf("semantic error: '%s' is a saga accept method and thus should return only a single `error` value",
				mi["Name"].(string))
		}

		for _, ri := range methodsInfo {
			if ri["Name"].(string) != sagaInfo.RollbackMethodName {
				continue
			}

			if !returnsOnlyError(ri) {
				return fmt.Errorf("semantic error: '%s' is a saga rollback method and thus should return only a single `error` value",
					ri["Name"].(string))
			}
		}
	}

	return nil
}

func (pf *ParsedFile) checkSagaRollbackMethodsExistAndMatch(funcInfo []map[string]interface{}) error {
	type methodInfo struct {
		arguments string
	}
	methodNames := make(map[string]methodInfo)

	for _, info := range funcInfo {
		methodNames[info["Name"].(string)] = methodInfo{
			arguments: info["SagaInfo"].(*SagaInfo).Arguments,
		}
	}

	for _, info := range funcInfo {
		sagaInfo := info["SagaInfo"].(*SagaInfo)
		if !sagaInfo.IsSaga {
			continue
		}

		if sagaInfo.NumArguments != 1 {
			return fmt.Errorf(
				"Semantic error: '%v' is a saga with %v arguments. "+
					"Currently only one argument is allowed (hint: use a structure).",
				info["Name"].(string), sagaInfo.NumArguments)
		}

		// INS_FLAG_NO_ROLLBACK_METHOD allows to make saga calls between different
		// contract types despite of missing corresponding syntax support. Obviously
		// if validation fail there will be no rollback method to call. Please use
		// this flag with extra care!
		if sagaInfo.RollbackMethodName == "INS_FLAG_NO_ROLLBACK_METHOD" {
			// skip following semantic checks
			continue
		}

		rollbackInfo, exists := methodNames[sagaInfo.RollbackMethodName]
		if !exists {
			return fmt.Errorf(
				"Semantic error: '%v' is a saga with rollback method '%v', "+
					"but '%v' is not declared. Maybe a typo?",
				info["Name"].(string), sagaInfo.RollbackMethodName, sagaInfo.RollbackMethodName)
		}

		acceptArguments := info["SagaInfo"].(*SagaInfo).Arguments
		if acceptArguments != rollbackInfo.arguments {
			return fmt.Errorf(
				"Semantic error: '%v' is a saga with arguments '%v' and rollback method '%v', "+
					"but '%v' arguments '%v' dont't match. They should be exactly the same.",
				info["Name"].(string), acceptArguments, sagaInfo.RollbackMethodName,
				sagaInfo.RollbackMethodName, rollbackInfo.arguments)
		}
	}
	return nil
}

func (pf *ParsedFile) functionInfoForWrapper(list []*ast.FuncDecl) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, len(list))
	for _, fun := range list {
		info := map[string]interface{}{
			"Name":                fun.Name.Name,
			"ArgumentsZeroList":   generateZeroListOfTypes(pf, "args", fun.Type.Params),
			"Arguments":           numberedVars(fun.Type.Params, "args"),
			"Results":             numberedVars(fun.Type.Results, "ret"),
			"ErrorInterfaceInRes": typeIndexes(pf, fun.Type.Results, errorType),
			"Immutable":           isImmutable(fun),  // only for methods, not constructors
			"SagaInfo":            sagaInfo(pf, fun), // only for methods, not constructors
		}
		res = append(res, info)
	}
	return res
}

// WriteProxy generates and writes into `out` source code of contract's proxy
func (pf *ParsedFile) WriteProxy(classReference string, out io.Writer) error {
	proxyPackageName, err := pf.ProxyPackageName()
	if err != nil {
		return err
	}

	if classReference == "" {
		classReference = genesisrefs.GenerateProtoReferenceFromCode(0, pf.code).String()
	}

	_, err = insolar.NewReferenceFromBase58(classReference)
	if err != nil {
		return errors.Wrap(err, "can't write proxy: ")
	}

	if err := checkMachineType(pf.machineType); err != nil {
		return err
	}

	allMethodsProxies := pf.functionInfoForProxy(pf.methods[pf.contract])

	err = pf.checkSagaRollbackMethodsExistAndMatch(allMethodsProxies)
	if err != nil {
		return err
	}

	err = pf.checkSagaIsNotImmutable(allMethodsProxies)
	if err != nil {
		return err
	}

	err = pf.checkSagaMethodsReturnOnlySingleErrorValue(allMethodsProxies)
	if err != nil {
		return err
	}

	constructorProxies := pf.functionInfoForProxy(pf.constructors[pf.contract])
	for _, fi := range constructorProxies {
		if fi["SagaInfo"].(*SagaInfo).IsSaga {
			return fmt.Errorf("semantic error: '%s' can't be a saga because it's a constructor", fi["Name"].(string))
		}
	}

	sagaRollbackMethods := make(map[string]struct{})
	for _, methodInfo := range allMethodsProxies {
		sagaInfo := methodInfo["SagaInfo"].(*SagaInfo)
		if sagaInfo.IsSaga {
			sagaRollbackMethods[sagaInfo.RollbackMethodName] = struct{}{}
		}
	}

	// explicitly remove all saga Rollback methods from the proxy
	var filteredMethodsProxies []map[string]interface{} //nolint:prealloc
	for _, methodInfo := range allMethodsProxies {
		currentMethodName := methodInfo["Name"].(string)
		_, isRollback := sagaRollbackMethods[currentMethodName]
		if isRollback {
			continue
		}
		filteredMethodsProxies = append(filteredMethodsProxies, methodInfo)
	}

	// Need to guarantee order for generated files
	types := generateTypes(pf)
	sort.Strings(types)

	data := map[string]interface{}{
		"PackageName":         proxyPackageName,
		"Types":               types,
		"ContractType":        pf.contract,
		"MethodsProxies":      filteredMethodsProxies,
		"ConstructorsProxies": constructorProxies,
		"ClassReference":      classReference,
		"Imports":             pf.generateImports(false),
	}

	return formatAndWrite(out, "proxy", data)
}

func (pf *ParsedFile) functionInfoForProxy(list []*ast.FuncDecl) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, len(list))

	for _, fun := range list {
		info := map[string]interface{}{
			"Name":                fun.Name.Name,
			"Arguments":           genFieldList(pf, fun.Type.Params, true),
			"InitArgs":            generateInitArguments(fun.Type.Params),
			"ResultZeroList":      generateZeroListOfTypes(pf, "ret", fun.Type.Results),
			"Results":             numberedVars(fun.Type.Results, "ret"),
			"ErrorVar":            fmt.Sprintf("ret%d", fun.Type.Results.NumFields()-1),
			"ResultsWithErr":      commaAppend(numberedVarsI(fun.Type.Results.NumFields()-1, "ret"), "err"),
			"ResultsNilError":     commaAppend(numberedVarsI(fun.Type.Results.NumFields()-1, "ret"), "nil"),
			"ResultsTypes":        genFieldList(pf, fun.Type.Results, false),
			"ErrorInterfaceInRes": typeIndexes(pf, fun.Type.Results, errorType),
			"Immutable":           isImmutable(fun),
			"SagaInfo":            sagaInfo(pf, fun),
		}
		res = append(res, info)
	}
	return res
}

// ChangePackageToMain changes package of the parsed code to "main"
func (pf *ParsedFile) ChangePackageToMain() {
	pf.node.Name.Name = mainPkg
}

// Write prints `out` contract's code, it could be changed with a few methods
func (pf *ParsedFile) Write(out io.Writer) error {
	return printer.Fprint(out, pf.fileSet, pf.node)
}

// codeOfNode returns source code of an AST node
func (pf *ParsedFile) codeOfNode(n ast.Node) string {
	return string(pf.code[n.Pos()-1 : n.End()-1])
}

func (pf *ParsedFile) typeName(t ast.Expr) string {
	if tmp, ok := t.(*ast.StarExpr); ok { // *type
		t = tmp.X
	}
	return pf.codeOfNode(t)
}

func (pf *ParsedFile) generateImports(wrapper bool) map[string]bool {
	imports := make(map[string]bool)
	imports[fmt.Sprintf(`"%s"`, proxyctxPath)] = true
	imports[fmt.Sprintf(`"%s"`, foundationPath)] = true
	if !wrapper {
		imports[fmt.Sprintf(`"%s"`, corePath)] = true
	}
	for _, method := range pf.methods[pf.contract] {
		extendImportsMap(pf, method.Type.Params, imports)
		if !wrapper {
			extendImportsMap(pf, method.Type.Results, imports)
		}
	}
	for _, fun := range pf.constructors[pf.contract] {
		extendImportsMap(pf, fun.Type.Params, imports)
		if !wrapper {
			extendImportsMap(pf, fun.Type.Results, imports)
		}
	}

	return imports
}

func openTemplate(fileName string) (*template.Template, error) {
	functions := template.FuncMap{"Title": strings.Title}
	tmpl := template.New(path.Base(fileName)).Funcs(functions)

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, errors.Wrap(nil, "couldn't find info about current file")
	}
	templateDir := filepath.Join(filepath.Dir(currentFile), fileName)
	tmpl, err := tmpl.ParseFiles(templateDir)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't parse template for output")
	}

	return tmpl, nil
}

func numberedVars(list *ast.FieldList, name string) string {
	if list == nil || list.NumFields() == 0 {
		return ""
	}
	return numberedVarsI(list.NumFields(), name)
}

func commaAppend(l string, r string) string {
	if l == "" {
		return r
	}
	return l + ", " + r
}

func numberedVarsI(n int, name string) string {
	if n == 0 {
		return ""
	}

	res := ""
	for i := 0; i < n; i++ {
		res = commaAppend(res, name+strconv.Itoa(i))
	}
	return res
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

func isContractTypeSpec(typeNode *ast.TypeSpec) bool {
	baseContract := "foundation.BaseContract"
	st, ok := typeNode.Type.(*ast.StructType)
	if !ok {
		return false
	}
	if st.Fields == nil || st.Fields.NumFields() == 0 {
		return false
	}
	for _, fd := range st.Fields.List {
		if len(fd.Names) != 0 {
			continue // named struct field
		}
		selectField, ok := fd.Type.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		pack := selectField.X.(*ast.Ident).Name
		class := selectField.Sel.Name
		if baseContract == (pack + "." + class) {
			return true
		}
	}

	return false
}

func generateTypes(parsed *ParsedFile) []string {
	types := make([]string, 0, len(parsed.types))
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
		if parsed.codeOfNode(e.Type) == errorType {
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

func generateZeroListOfTypes(parsed *ParsedFile, name string, list *ast.FieldList) string {
	if list == nil || list.NumFields() == 0 {
		return fmt.Sprintf("%s := []interface{}{}\n", name)
	}

	text := fmt.Sprintf("%s := make([]interface{}, %d)\n", name, list.NumFields())

	for i, arg := range list.List {
		tname := parsed.codeOfNode(arg.Type)
		if tname == errorType {
			tname = "*foundation.Error"
		}

		text += fmt.Sprintf("\tvar %s%d %s\n", name, i, tname)
		text += fmt.Sprintf("\t%s[%d] = &%s%d\n", name, i, name, i)
	}

	return text
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

// GetRealApplicationDir return application dir path
func GetRealApplicationDir(dir string) (string, error) {
	gopath := build.Default.GOPATH
	if gopath == "" {
		return "", errors.Errorf("GOPATH is not set")
	}
	contractsPath := ""
	for _, p := range strings.Split(gopath, ":") {
		contractsPath = path.Join(p, "src/github.com/insolar/insolar/application/", dir)
		_, err := os.Stat(contractsPath)
		if err == nil {
			return contractsPath, nil
		}
	}
	return "", errors.Errorf("Not found github.com/insolar/insolar in GOPATH")
}

// GetRealContractsNames returns names of all real smart contracts
func GetRealContractsNames() ([]string, error) {
	pathWithContracts, err := GetRealApplicationDir("contract")
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

func isImmutable(decl *ast.FuncDecl) bool {
	var isImmutable = false
	if decl.Doc != nil && decl.Doc.List != nil {
		for _, comment := range decl.Doc.List {
			slice, err := skipCommentBeginning(comment.Text)
			if err != nil {
				// invalid comment beginning
				continue
			}
			if slice == immutableFlag {
				isImmutable = true
				break
			}
		}
	}
	return isImmutable
}

// skipCommentBegin converts '//comment' or '//[spaces]comment' to 'comment'
// The procedure returns an error if the string is not started with '//'
func skipCommentBeginning(comment string) (string, error) {
	slice := strings.TrimSpace(comment)
	sliceLen := len(slice)

	// skip '//'
	if !strings.HasPrefix(slice, "//") {
		return "", fmt.Errorf("invalid comment beginning")
	}
	slice = slice[2:sliceLen]
	sliceLen -= 2

	// skip all whitespaces after '//'
	for sliceLen > 0 && (slice[0] == ' ' || slice[0] == '\t') {
		slice = slice[1:sliceLen]
		sliceLen--
	}

	return slice, nil
}

func extractSagaInfoFromComment(comment string, info *SagaInfo) bool {
	slice, err := skipCommentBeginning(comment)
	if err != nil {
		return false
	}
	if strings.HasPrefix(slice, sagaFlagStart) &&
		strings.HasSuffix(slice, sagaFlagEnd) {
		rollbackName := slice[sagaFlagStartLength : len(slice)-len(sagaFlagEnd)]
		rollbackNameLen := len(rollbackName)
		if rollbackNameLen > 0 {
			sliceCopy := make([]byte, rollbackNameLen)
			copy(sliceCopy, rollbackName)
			info.IsSaga = true
			info.RollbackMethodName = string(sliceCopy)
			return true
		}
	}
	return false
}

func sagaInfo(pf *ParsedFile, decl *ast.FuncDecl) (info *SagaInfo) {
	info = &SagaInfo{
		Arguments:    genFieldList(pf, decl.Type.Params, true),
		NumArguments: len(decl.Type.Params.List),
	}
	if decl.Doc == nil || decl.Doc.List == nil {
		return // there are no comments
	}

	for _, comment := range decl.Doc.List {
		if extractSagaInfoFromComment(comment.Text, info) {
			return // info found
		}
	}

	return // no saga comment found
}

type ContractListEntry struct {
	Name       string
	Path       string
	Parsed     *ParsedFile
	ImportPath string
	Version    int
}

const (
	CodeType      = "code"
	PrototypeType = "prototype"
)

type ContractList []ContractListEntry

func generateContractList(contracts ContractList) interface{} {
	importList := make([]interface{}, 0)
	for _, contract := range contracts {
		data := map[string]interface{}{
			"Name":               contract.Name,
			"ImportName":         contract.Name,
			"ImportPath":         contract.ImportPath,
			"CodeReference":      genesisrefs.GenerateCodeReferenceFromContractID(CodeType, contract.Name, contract.Version).String(),
			"PrototypeReference": genesisrefs.GenerateProtoReferenceFromContractID(PrototypeType, contract.Name, contract.Version).String(),
		}
		importList = append(importList, data)
	}
	return importList
}

func GenerateInitializationList(out io.Writer, contracts ContractList) error {
	data := map[string]interface{}{
		"Contracts": generateContractList(contracts),
		"Package":   "builtin",
	}

	return formatAndWrite(out, "initialization", data)
}
