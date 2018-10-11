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
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
)

var randomTestCode = `
package main

import (
	"fmt"

	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type HelloWorlder struct {
	foundation.BaseContract
	Greeted int
}

type FullName struct {
	First string
	Last  string
}

type PersonalGreeting struct {
	Name    FullName
	Message string
}

type Error struct {
	S string
}

func (e *Error) Error() string {
	return e.S
}

func (hw *HelloWorlder) Hello() (string, *Error) {
	hw.Greeted++
	return "Hello world 2", nil
}

func (hw *HelloWorlder) Fail() (string, *Error) {
	hw.Greeted++
	return "", &Error{"We failed 2"}
}

func (hw *HelloWorlder) Echo(s string) (string, *Error) {
	hw.Greeted++
	return s, nil
}

func (hw *HelloWorlder) HelloHuman(Name FullName) PersonalGreeting {
	hw.Greeted++
	return PersonalGreeting{
		Name:    Name,
		Message: fmt.Sprintf("Dear %s %s, we specially say hello to you", Name.First, Name.Last),
	}
}

func (hw *HelloWorlder) HelloHumanPointer(Name FullName) *PersonalGreeting {
	hw.Greeted++
	return &PersonalGreeting{
		Name:    Name,
		Message: fmt.Sprintf("Dear %s %s, we specially say hello to you", Name.First, Name.Last),
	}
}

func (hw *HelloWorlder) MultiArgs(Name FullName, s string, i int) *PersonalGreeting {
	hw.Greeted++
	return &PersonalGreeting{
		Name:    Name,
		Message: fmt.Sprintf("Dear %s %s, we specially say hello to you", Name.First, Name.Last),
	}
}

func (hw HelloWorlder) ConstEcho(s string) (string, *Error) {
	return s, nil
}

func JustExportedStaticFunction(int, int) {}
`

func TestBasicGeneration(t *testing.T) {
	t.Parallel()
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir) // nolint: errcheck

	err = goplugintestutils.WriteFile(tmpDir, "main.go", randomTestCode)
	assert.NoError(t, err)

	parsed, err := ParseFile(filepath.Join(tmpDir, "main.go"))
	assert.NoError(t, err)
	assert.NotNil(t, parsed)

	t.Run("wrapper", func(t *testing.T) {
		t.Parallel()

		buf := bytes.Buffer{}
		err = parsed.WriteWrapper(&buf)
		assert.NoError(t, err)

		code, err := ioutil.ReadAll(&buf)
		assert.NoError(t, err)
		assert.NotEmpty(t, code)
	})

	t.Run("proxy", func(t *testing.T) {
		t.Parallel()

		buf := bytes.Buffer{}
		err = parsed.WriteProxy("testRef", &buf)
		assert.NoError(t, err)

		code, err := ioutil.ReadAll(&buf)
		assert.NoError(t, err)
		assert.NotEmpty(t, code)
	})
}

func TestConstructorsParsing(t *testing.T) {
	t.Parallel()
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir) // nolint: errcheck

	code := `
package main

type One struct {
foundation.BaseContract
}

func New() *One {
	return &One{}
}

func NewFromString(s string) *One {
	return &One{}
}

func NewWrong() {
}
`

	err = goplugintestutils.WriteFile(tmpDir, "code1", code)
	assert.NoError(t, err)

	info, err := ParseFile(filepath.Join(tmpDir, "code1"))
	assert.NoError(t, err)

	assert.Equal(t, 1, len(info.constructors))
	assert.Equal(t, 2, len(info.constructors["One"]))
	assert.Equal(t, "New", info.constructors["One"][0].Name.Name)
	assert.Equal(t, "NewFromString", info.constructors["One"][1].Name.Name)
}

func TestCompileContractProxy(t *testing.T) {
	t.Parallel()

	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir) // nolint: errcheck

	err = os.MkdirAll(filepath.Join(tmpDir, "src/secondary"), 0777)
	assert.NoError(t, err)

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	// XXX: dirty hack to make `dep` installed packages available in generated code
	err = os.Symlink(filepath.Join(cwd, "../../../vendor"), filepath.Join(tmpDir, "src/secondary/vendor"))
	assert.NoError(t, err)

	proxyFh, err := os.OpenFile(filepath.Join(tmpDir, "/src/secondary/main.go"), os.O_WRONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	err = goplugintestutils.WriteFile(filepath.Join(tmpDir, "/contracts/secondary/"), "main.go", randomTestCode)
	assert.NoError(t, err)

	parsed, err := ParseFile(filepath.Join(tmpDir, "/contracts/secondary/main.go"))
	assert.NoError(t, err)

	err = parsed.WriteProxy("testRef", proxyFh)
	assert.NoError(t, err)

	err = proxyFh.Close()
	assert.NoError(t, err)

	err = goplugintestutils.WriteFile(tmpDir, "/test.go", `
package test

import (
	"github.com/insolar/insolar/core"
	"secondary"
)

func main() {
	_ = secondary.GetObject(core.NewRefFromBase58("some"))
}
	`)
	assert.NoError(t, err)

	cmd := exec.Command("go", "build", filepath.Join(tmpDir, "test.go"))
	cmd.Env = append(os.Environ(), "GOPATH="+goplugintestutils.PrependGoPath(tmpDir))
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err, string(out))
}

func TestGenerateProxyAndWrapperWithoutReturnValue(t *testing.T) {
	t.Parallel()
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = goplugintestutils.WriteFile(tmpDir, testContract, `
package main
type A struct{
	int C
	foundation.BaseContract
	int M
}

func ( A ) Get(){
	return
}
`)
	assert.NoError(t, err)

	parsed, err := ParseFile(tmpDir + testContract)
	assert.NoError(t, err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy("testRef", &bufProxy)
	assert.NoError(t, err)
	code, err := ioutil.ReadAll(&bufProxy)
	assert.NoError(t, err)
	assert.NotEqual(t, len(code), 0)

	var bufWrapper bytes.Buffer
	err = parsed.WriteWrapper(&bufWrapper)
	assert.NoError(t, err)
	assert.Contains(t, bufWrapper.String(), "    self.Get(  )")
}

func TestFailIfThereAreNoContract(t *testing.T) {
	t.Parallel()
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir) //nolint: errcheck

	testContract := "/test.go"
	err = goplugintestutils.WriteFile(tmpDir, testContract, `
package main
type A struct{
	ttt ppp.TTT
}
`)
	assert.NoError(t, err)

	_, err = ParseFile(tmpDir + testContract)
	assert.EqualError(t, err, "Only one smart contract must exist")
}

func TestInitializationFunctionParamsProxy(t *testing.T) {
	t.Parallel()
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir) // nolint: errcheck

	testContract := "/test.go"

	err = goplugintestutils.WriteFile(tmpDir, testContract, `
package main

type A struct{
	foundation.BaseContract
}

func ( a *A )Get( a int, b bool, c string, d foundation.Reference ) ( int, bool, string, foundation.Reference ){
	return
}
`)

	assert.NoError(t, err)

	parsed, err := ParseFile(tmpDir + testContract)
	assert.NoError(t, err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy("testRef", &bufProxy)
	assert.NoError(t, err)
	assert.Contains(t, bufProxy.String(), "var a0 int")
	assert.Contains(t, bufProxy.String(), "resList[0] = a0")

	assert.Contains(t, bufProxy.String(), "var a1 bool")
	assert.Contains(t, bufProxy.String(), "resList[1] = a1")

	assert.Contains(t, bufProxy.String(), "var a2 string")
	assert.Contains(t, bufProxy.String(), "resList[2] = a2")

	assert.Contains(t, bufProxy.String(), "var a3 foundation.Reference")
	assert.Contains(t, bufProxy.String(), "resList[3] = a3")
}

func TestInitializationFunctionParamsWrapper(t *testing.T) {
	t.Parallel()
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir) //nolint: errcheck

	testContract := "/test.go"

	err = goplugintestutils.WriteFile(tmpDir, testContract, `
package main

type A struct{
	foundation.BaseContract
}

func ( a *A )Get( a int, b bool, c string, d foundation.Reference ) ( int, bool, string, foundation.Reference ){
	return
}
`)
	assert.NoError(t, err)

	parsed, err := ParseFile(tmpDir + testContract)
	assert.NoError(t, err)

	var bufWrapper bytes.Buffer
	err = parsed.WriteWrapper(&bufWrapper)
	assert.NoError(t, err)
	assert.Contains(t, bufWrapper.String(), "var a0 int")
	assert.Contains(t, bufWrapper.String(), "args[0] = a0")

	assert.Contains(t, bufWrapper.String(), "var a1 bool")
	assert.Contains(t, bufWrapper.String(), "args[1] = a1")

	assert.Contains(t, bufWrapper.String(), "var a2 string")
	assert.Contains(t, bufWrapper.String(), "args[2] = a2")

	assert.Contains(t, bufWrapper.String(), "var a3 foundation.Reference")
	assert.Contains(t, bufWrapper.String(), "args[3] = a3")
}

func TestContractOnlyIfEmbedBaseContract(t *testing.T) {
	t.Parallel()
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir) //nolint: errcheck

	testContract := "/test.go"

	err = goplugintestutils.WriteFile(tmpDir, testContract, `
package main
// A contains object of type foundation.BaseContract, but it must embed it
type A struct{
	tt foundation.BaseContract
}
`)
	assert.NoError(t, err)

	_, err = ParseFile(tmpDir + testContract)
	assert.EqualError(t, err, "Only one smart contract must exist")
}

func TestOnlyOneSmartContractMustExist(t *testing.T) {
	t.Parallel()
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir) //nolint: errcheck

	testContract := "/test.go"

	err = goplugintestutils.WriteFile(tmpDir, testContract, `
package main

type A struct{
	foundation.BaseContract
}

type B struct{
	foundation.BaseContract
}
`)
	assert.NoError(t, err)

	_, err = ParseFile(tmpDir + testContract)
	assert.EqualError(t, err, ": more than one contract in a file")
}

func TestImportsFromContract(t *testing.T) {
	t.Parallel()
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = goplugintestutils.WriteFile(tmpDir, testContract, `
package main
import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"some/test/import/path"
	"some/test/import/pointerPath"
)

type A struct{
	foundation.BaseContract
}

func ( A ) Get(i path.SomeType){
	return
}

func ( A ) GetPointer(i *pointerPath.SomeType){
	return
}
`)
	assert.NoError(t, err)

	parsed, err := ParseFile(tmpDir + testContract)
	assert.NoError(t, err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy("testRef", &bufProxy)
	assert.NoError(t, err)
	assert.Contains(t, bufProxy.String(), `"some/test/import/path"`)
	assert.Contains(t, bufProxy.String(), `"some/test/import/pointerPath"`)
	assert.Contains(t, bufProxy.String(), `"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"`)
	code, err := ioutil.ReadAll(&bufProxy)
	assert.NoError(t, err)
	assert.NotEqual(t, len(code), 0)

	var bufWrapper bytes.Buffer
	err = parsed.WriteWrapper(&bufWrapper)
	assert.NoError(t, err)
	assert.Contains(t, bufWrapper.String(), `"some/test/import/path"`)
	assert.Contains(t, bufWrapper.String(), `"some/test/import/pointerPath"`)
}

func TestAliasImportsFromContract(t *testing.T) {
	t.Parallel()
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = goplugintestutils.WriteFile(tmpDir, testContract, `
package main
import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	someAlias "some/test/import/path"
)

type A struct{
	foundation.BaseContract
}

func ( A ) Get(i someAlias.SomeType){
	return
}
`)
	assert.NoError(t, err)

	parsed, err := ParseFile(tmpDir + testContract)
	assert.NoError(t, err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy("testRef", &bufProxy)
	assert.NoError(t, err)
	assert.Contains(t, bufProxy.String(), `someAlias "some/test/import/path"`)
	assert.Contains(t, bufProxy.String(), `"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"`)
	code, err := ioutil.ReadAll(&bufProxy)
	assert.NoError(t, err)
	assert.NotEqual(t, len(code), 0)

	var bufWrapper bytes.Buffer
	err = parsed.WriteWrapper(&bufWrapper)
	assert.NoError(t, err)
	assert.Contains(t, bufWrapper.String(), `someAlias "some/test/import/path"`)
	assert.NotContains(t, bufProxy.String(), `"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"`)
}

func TestImportsFromContractUseInsideFunc(t *testing.T) {
	t.Parallel()
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = goplugintestutils.WriteFile(tmpDir, testContract, `
package main
import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"some/test/import/path"
)

type A struct{
	foundation.BaseContract
}

func ( A ) Get() {
	path.SomeMethod()
	return
}
`)
	assert.NoError(t, err)

	parsed, err := ParseFile(tmpDir + testContract)
	assert.NoError(t, err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy("testRef", &bufProxy)
	assert.NoError(t, err)
	assert.NotContains(t, bufProxy.String(), `"some/test/import/path"`)
	code, err := ioutil.ReadAll(&bufProxy)
	assert.NoError(t, err)
	assert.NotEqual(t, len(code), 0)

	var bufWrapper bytes.Buffer
	err = parsed.WriteWrapper(&bufWrapper)
	assert.NoError(t, err)
	assert.NotContains(t, bufWrapper.String(), `"some/test/import/path"`)
}

func TestImportsFromContractUseForReturnValue(t *testing.T) {
	t.Parallel()
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = goplugintestutils.WriteFile(tmpDir, testContract, `
package main
import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"some/test/import/path"
)

type A struct{
	foundation.BaseContract
}

func ( A ) Get() path.SomeValue {
	f := path.SomeMethod()
	return f
}
`)
	assert.NoError(t, err)

	parsed, err := ParseFile(tmpDir + testContract)
	assert.NoError(t, err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy("testRef", &bufProxy)
	assert.NoError(t, err)
	assert.Contains(t, bufProxy.String(), `"some/test/import/path"`)
	code, err := ioutil.ReadAll(&bufProxy)
	assert.NoError(t, err)
	assert.NotEqual(t, len(code), 0)

	var bufWrapper bytes.Buffer
	err = parsed.WriteWrapper(&bufWrapper)
	assert.NoError(t, err)
	assert.NotContains(t, bufWrapper.String(), `"some/test/import/path"`)
}

func TestNotMatchFileNameForProxy(t *testing.T) {
	t.Parallel()
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test_not_go_file.test"
	err = goplugintestutils.WriteFile(tmpDir, testContract, `
package main

type A struct{
	foundation.BaseContract
}
`)
	assert.NoError(t, err)

	parsed, err := ParseFile(tmpDir + testContract)
	assert.NoError(t, err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy("testRef", &bufProxy)
	assert.EqualError(t, err, "couldn't match filename without extension and path")
}
