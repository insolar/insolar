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
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

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

func generateContractWrapper(fileName string, out io.Writer) error {
	parsed, err := parseFile(fileName)
	if err != nil {
		return errors.Wrap(err, "couldn't parse")
	}

	packageName := parsed.node.Name.Name
	if packageName != "main" {
		panic("Contract must be in main package")
	}

	code := "package " + packageName + "\n\n"
	code += generateWrappers(parsed) + "\n"
	code += generateExports(parsed) + "\n"

	_, err = out.Write([]byte(code))
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

	code := "package " + proxyPackageName + "\n\n"

	code += `import (
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

`

	code += generateTypes(parsed) + "\n"

	code += `// Contract proxy type
type ` + parsed.contract + ` struct {
	Reference string
}

`

	code += `// GetObject
func GetObject(ref string) (r *` + parsed.contract + `) {
	return &` + parsed.contract + `{Reference: ref}
}
`

	code += generateMethodsProxies(parsed) + "\n"

	_, err = out.Write([]byte(code))
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
func generateTypes(parsed *parsedFile) string {
	text := ""
	for _, t := range parsed.types {
		text += "type " + string(parsed.code[t.Pos()-1:t.End()-1]) + "\n\n"
	}

	return text
}

func generateWrappers(parsed *parsedFile) string {
	text := `import (
	"github.com/insolar/insolar/logicrunner/goplugin/testplugins/foundation"
	)` + "\n"

	for _, method := range parsed.methods[parsed.contract] {
		text += generateMethodWrapper(parsed, method) + "\n"
	}
	return text
}

func generateZeroListOfTypes(parsed *parsedFile, name string, list *ast.FieldList) (string, string) {
	text := fmt.Sprintf("\t%s := [%d]interface{}{}\n", name, list.NumFields())

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

func generateMethodWrapper(parsed *parsedFile, method *ast.FuncDecl) string {
	text := fmt.Sprintf("func (self *%s) INSWRAPER_%s(cbor foundation.CBORMarshaler, data []byte) ([]byte) {\n",
		parsed.contract, method.Name.Name)

	argsInit, argsList := generateZeroListOfTypes(parsed, "args", method.Type.Params)
	text += argsInit

	text += "\tcbor.Unmarshal(&args, data)\n"

	rets := []string{}
	for i := range method.Type.Results.List {
		rets = append(rets, fmt.Sprintf("ret%d", i))
	}
	ret := strings.Join(rets, ", ")
	text += fmt.Sprintf("\t%s := self.%s(%s)\n", ret, method.Name.Name, argsList)

	text += fmt.Sprintf("\treturn cbor.Marshal([]interface{}{%s})\n", strings.Join(rets, ", "))
	text += "}\n"
	return text
}

/* generated snipped must be something like this

func (hw *HelloWorlder) INSWRAPER_Echo(cbor cborer, data []byte) ([]byte, error) {
	args := [1]interface{}{}
	args[0] = ""
	cbor.Unmarshal(&args, data)
	ret1, ret2 := hw.Echo(args[0].(string))
	return cbor.Marshal([]interface{}{ret1, ret2}), nil
}
*/

func generateExports(parsed *parsedFile) string {
	text := "var INSEXPORT " + parsed.contract + "\n"
	return text
}

func generateMethodsProxies(parsed *parsedFile) string {
	text := ""

	for _, method := range parsed.methods[parsed.contract] {
		text += generateMethodProxy(parsed, method) + "\n"
	}
	return text
}

func generateMethodProxy(parsed *parsedFile, method *ast.FuncDecl) string {
	text := fmt.Sprintf("func (r *%s) %s(", parsed.contract, method.Name.Name)
	for i, arg := range method.Type.Params.List {
		if i > 0 {
			text += ", "
		}
		text += arg.Names[0].Name
		text += " " + string(parsed.code[arg.Type.Pos()-1:arg.Type.End()-1])
	}
	text += ") ("

	for i, arg := range method.Type.Results.List {
		if i > 0 {
			text += ", "
		}
		text += string(parsed.code[arg.Type.Pos()-1 : arg.Type.End()-1])
	}

	text += ") {\n"

	text += fmt.Sprintf("\tvar args [%d]interface{}\n", method.Type.Params.NumFields())
	for i, arg := range method.Type.Params.List {
		text += fmt.Sprintf("\targs[%d] = %s\n", i, arg.Names[0].Name)
	}

	text += `
	var argsSerialized []byte
	err := proxyctx.Current.Serialize(args, &argsSerialized)
	if err != nil {
		panic(err)
	}
`

	text += fmt.Sprintf("\t"+`res, err := proxyctx.Current.RouteCall(r.Reference, "%s", argsSerialized)`, method.Name.Name)

	text += `
	if err != nil {
		panic(err)
	}
`
	resInit, resList := generateZeroListOfTypes(parsed, "resList", method.Type.Results)
	text += resInit

	text += `
	err = proxyctx.Current.Deserialize(res, &resList)
	if err != nil {
		panic(err)
	}
`

	if method.Type.Results.NumFields() > 0 {
		text += "\treturn " + resList
	}

	text += "}\n"

	return text
}
