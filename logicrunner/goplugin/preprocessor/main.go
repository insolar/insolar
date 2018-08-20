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

	types    map[string]string
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
		fs.Parse(os.Args[2:])

		for _, fn := range fs.Args() {
			err := generateContractWrapper(fn, output.writer)
			if err != nil {
				panic(err)
			}
		}
	case "proxy":
		fs := flag.NewFlagSet("proxy", flag.ExitOnError)
		fs.Parse(os.Args[2:])
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
	defer file.Close()

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
	out.Write([]byte("package " + packageName + "\n\n"))
	out.Write([]byte(generateWrappers(parsed) + "\n"))
	out.Write([]byte(generateExports(parsed) + "\n"))
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

	out.Write([]byte("package " + proxyPackageName + "\n\n"))

	out.Write([]byte(`import (
	"github.com/ugorji/go/codec"
)

`))

	out.Write([]byte(`// Contract proxy type
type ` + parsed.contract + ` struct {
	Reference string
	RPC struct{}
}

`))

	out.Write([]byte(`// GetObject
func GetObject(ref string) (r *` + parsed.contract + `) {
	return &` + parsed.contract + `{}
}
`))

	out.Write([]byte(generateMethodsProxies(parsed) + "\n"))
	return nil
}

func getMethods(parsed *parsedFile) {
	parsed.types = make(map[string]string)
	parsed.methods = make(map[string][]*ast.FuncDecl)
	for _, d := range parsed.node.Decls {
		switch td := d.(type) {
		case *ast.GenDecl:
			if td.Tok != token.TYPE {
				continue
			}

			typeNode := td.Specs[0].(*ast.TypeSpec)
			if strings.Contains(td.Doc.Text(), "@inscontract") {
				if parsed.contract != "" {
					panic("more than one contract in a file")
				}
				parsed.contract = typeNode.Name.Name
			} else {
				parsed.types[typeNode.Name.Name] = string(parsed.code[typeNode.Pos()-1 : typeNode.End()])
			}
		case *ast.FuncDecl:
			if td.Recv == nil || td.Recv.NumFields() == 0 { // not a method
				continue
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
		text += "type " + t + "\n"
	}

	return text
}

func generateWrappers(parsed *parsedFile) string {
	text := `import (
	"github.com/insolar/insolar/logicrunner/goplugin/testplugins/foundation"
	)` + "\n"

	for _, method := range parsed.methods[parsed.contract] {
		text += generateMethodWrapper(method, parsed.contract) + "\n"
	}
	return text
}

func generateMethodWrapper(method *ast.FuncDecl, class string) string {
	text := fmt.Sprintf("func (self *%s) INSWRAPER_%s(cbor foundation.CBORMarshaler, data []byte) ([]byte) {\n",
		class, method.Name.Name)
	text += fmt.Sprintf("\targs := [%d]interface{}{}\n", method.Type.Params.NumFields())

	args := []string{}
	for i, arg := range method.Type.Params.List {
		initializer := ""
		tname := fmt.Sprintf("%v", arg.Type)
		switch tname {
		case "uint", "int", "int8", "uint8", "int32", "uint32", "int64", "uint64":
			initializer = tname + "(0)"
		case "string":
			initializer = `""`
		default:
			initializer = tname + "{}"
		}
		text += fmt.Sprintf("\targs[%d] = %s\n", i, initializer)
		args = append(args, fmt.Sprintf("args[%d].(%s)", i, tname))
	}

	text += "\tcbor.Unmarshal(&args, data)\n"

	rets := []string{}
	for i := range method.Type.Results.List {
		rets = append(rets, fmt.Sprintf("ret%d", i))
	}
	ret := strings.Join(rets, ", ")
	text += fmt.Sprintf("\t%s := self.%s(%s)\n", ret, method.Name.Name, strings.Join(args, ", "))

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
	text += `
	ch := new(codec.CborHandle)
	var data []byte
	err := codec.NewEncoderBytes(&data, ch).Encode(*r)
	if err != nil {
		panic(err)
	}
`

	text += fmt.Sprintf("\tvar args [%d]interface{}\n", method.Type.Params.NumFields())
	for i, arg := range method.Type.Params.List {
		text += fmt.Sprintf("\targs[%d] = %s\n", i, arg.Names[0].Name)
	}

	text += `
	var argsSerialized []byte
	err = codec.NewEncoderBytes(&argsSerialized, ch).Encode(args)
	if err != nil {
		panic(err)
	}
`

	text += fmt.Sprintf(`\tdata, res, err := r.RPC.Exec(r.Reference, "%s", argsSerialized)`, method.Name.Name)

	text += `
	if err != nil {
		panic(err)
	}
`
	text += `
	err = codec.NewDecoderBytes(data, ch).Decode(r)
	if err != nil {
		panic(err)
	}
`
	text += "}\n"

	return text
}
