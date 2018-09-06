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
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/logicrunner/goplugin/testutil"
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
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir) // nolint: errcheck

	err = testutil.WriteFile(tmpDir, "main.go", randomTestCode)
	assert.NoError(t, err)

	t.Run("wrapper", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := generateContractWrapper(tmpDir+"/main.go", &buf)
		assert.NoError(t, err)

		code, err := ioutil.ReadAll(&buf)
		assert.NoError(t, err)
		assert.NotEmpty(t, code)
	})

	t.Run("proxy", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := generateContractProxy(tmpDir+"/main.go", "testRef", &buf)
		assert.NoError(t, err)

		code, err := ioutil.ReadAll(&buf)
		assert.NoError(t, err)
		assert.NotEmpty(t, code)
	})
}

func TestConstructorsParsing(t *testing.T) {
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

	err = testutil.WriteFile(tmpDir, "code1", code)
	assert.NoError(t, err)

	info, err := parseFile(tmpDir + "/code1")
	assert.NoError(t, err)

	assert.Equal(t, 1, len(info.constructors))
	assert.Equal(t, 2, len(info.constructors["One"]))
	assert.Equal(t, "New", info.constructors["One"][0].Name.Name)
	assert.Equal(t, "NewFromString", info.constructors["One"][1].Name.Name)
}

func TestCompileContractProxy(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(cwd) // nolint: errcheck

	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir) // nolint: errcheck

	err = os.MkdirAll(tmpDir+"/src/secondary/", 0777)
	assert.NoError(t, err)

	// XXX: dirty hack to make `dep` installed packages available in generated code
	err = os.Symlink(cwd+"/../../../vendor/", tmpDir+"/src/secondary/vendor")
	assert.NoError(t, err)

	proxyFh, err := os.OpenFile(tmpDir+"/src/secondary/main.go", os.O_WRONLY|os.O_CREATE, 0644)
	assert.NoError(t, err)

	err = testutil.WriteFile(tmpDir+"/contracts/secondary/", "main.go", randomTestCode)
	assert.NoError(t, err)

	err = generateContractProxy(tmpDir+"/contracts/secondary/main.go", "testRef", proxyFh)
	assert.NoError(t, err)

	err = proxyFh.Close()
	assert.NoError(t, err)

	err = testutil.WriteFile(tmpDir, "/test.go", `
package test

import "secondary"

func main() {
	_ = secondary.GetObject("some")
}
	`)
	assert.NoError(t, err)

	err = os.Chdir(tmpDir)
	assert.NoError(t, err)

	origGoPath, err := testutil.ChangeGoPath(tmpDir)
	assert.NoError(t, err)
	defer os.Setenv("GOPATH", origGoPath) // nolint: errcheck

	out, err := exec.Command("go", "build", "test.go").CombinedOutput()
	assert.NoError(t, err, string(out))
}

func TestGenerateProxyAndWrapperForBoolParams(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir) //nolint: errcheck

	testContract := "/test.go"
	err = testutil.WriteFile(tmpDir, testContract, `
package main
type A struct{
foundation.BaseContract
}

func ( A ) Get( b bool ) bool{
	return true
}
`)
	assert.NoError(t, err)

	var bufProxy bytes.Buffer
	err = generateContractProxy(tmpDir+testContract, "testRef", &bufProxy)
	assert.NoError(t, err)
	assert.Contains(t, bufProxy.String(), "resList[0] = bool(false)")

	var bufWrapper bytes.Buffer
	err = generateContractWrapper(tmpDir+testContract, &bufWrapper)
	assert.NoError(t, err)
	assert.Contains(t, bufWrapper.String(), "args[0] = bool(false)")

}

func TestGenerateProxyAndWrapperWithoutReturnValue(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = testutil.WriteFile(tmpDir, testContract, `
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

	var bufProxy bytes.Buffer
	err = generateContractProxy(tmpDir+testContract, "testRef", &bufProxy)
	assert.NoError(t, err)
	code, err := ioutil.ReadAll(&bufProxy)
	assert.NoError(t, err)
	assert.NotEqual(t, len(code), 0)

	var bufWrapper bytes.Buffer
	err = generateContractWrapper(tmpDir+testContract, &bufWrapper)
	assert.NoError(t, err)
	assert.Contains(t, bufWrapper.String(), "    self.Get(  )")
}

func TestFailIfThereAreNoContract(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir) //nolint: errcheck

	testContract := "/test.go"
	err = testutil.WriteFile(tmpDir, testContract, `
package main
type A struct{
	ttt ppp.TTT
}
`)
	assert.NoError(t, err)

	var bufProxy bytes.Buffer
	err = generateContractProxy(tmpDir+testContract, "testRef", &bufProxy)
	assert.EqualError(t, err, "couldn't parse: Only one smart contract must exist")

	var bufWrapper bytes.Buffer
	err = generateContractWrapper(tmpDir+testContract, &bufWrapper)
	assert.EqualError(t, err, "couldn't parse: Only one smart contract must exist")
}

func TestContractOnlyIfEmbedBaseContract(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir) //nolint: errcheck

	testContract := "/test.go"

	err = testutil.WriteFile(tmpDir, testContract, `
package main
// A contains object of type foundation.BaseContract, but it must embed it
type A struct{
	tt foundation.BaseContract
}
`)
	assert.NoError(t, err)

	var bufProxy bytes.Buffer
	err = generateContractProxy(tmpDir+testContract, "testRef", &bufProxy)
	assert.EqualError(t, err, "couldn't parse: Only one smart contract must exist")

	var bufWrapper bytes.Buffer
	err = generateContractWrapper(tmpDir+testContract, &bufWrapper)
	assert.EqualError(t, err, "couldn't parse: Only one smart contract must exist")

}

func TestOnlyOneSmartContractMustExist(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir) //nolint: errcheck

	testContract := "/test.go"

	err = testutil.WriteFile(tmpDir, testContract, `
package main

type A struct{
	foundation.BaseContract
}

type B struct{
	foundation.BaseContract
}
`)
	assert.NoError(t, err)

	var bufProxy bytes.Buffer
	err = generateContractProxy(tmpDir+testContract, "testRef", &bufProxy)
	assert.EqualError(t, err, "couldn't parse: : more than one contract in a file")

	var bufWrapper bytes.Buffer
	err = generateContractWrapper(tmpDir+testContract, &bufWrapper)
	assert.EqualError(t, err, "couldn't parse: : more than one contract in a file")
}

func TestImportsFromContractInWrapper(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = testutil.WriteFile(tmpDir, testContract, `
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

	var bufProxy bytes.Buffer
	err = generateContractProxy(tmpDir+testContract, "testRef", &bufProxy)
	assert.NoError(t, err)
	code, err := ioutil.ReadAll(&bufProxy)
	assert.NoError(t, err)
	assert.NotEqual(t, len(code), 0)

	var bufWrapper bytes.Buffer
	err = generateContractWrapper(tmpDir+testContract, &bufWrapper)
	assert.NoError(t, err)
	assert.Contains(t, bufWrapper.String(), `"some/test/import/path"`)
	assert.Contains(t, bufWrapper.String(), `"some/test/import/pointerPath"`)
	assert.Contains(t, bufWrapper.String(), `"github.com/insolar/insolar/logicrunner/goplugin/foundation"`)
}

func TestAliasImportsFromContractInWrapper(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = testutil.WriteFile(tmpDir, testContract, `
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

	var bufProxy bytes.Buffer
	err = generateContractProxy(tmpDir+testContract, "testRef", &bufProxy)
	assert.NoError(t, err)
	code, err := ioutil.ReadAll(&bufProxy)
	assert.NoError(t, err)
	assert.NotEqual(t, len(code), 0)

	var bufWrapper bytes.Buffer
	err = generateContractWrapper(tmpDir+testContract, &bufWrapper)
	assert.NoError(t, err)
	assert.Contains(t, bufWrapper.String(), `someAlias "some/test/import/path"`)
	assert.Contains(t, bufWrapper.String(), `"github.com/insolar/insolar/logicrunner/goplugin/foundation"`)
}

func TestImportsFromContractNotInWrapperIfNotInputValue(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test.go"
	err = testutil.WriteFile(tmpDir, testContract, `
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

	var bufProxy bytes.Buffer
	err = generateContractProxy(tmpDir+testContract, "testRef", &bufProxy)
	assert.NoError(t, err)
	code, err := ioutil.ReadAll(&bufProxy)
	assert.NoError(t, err)
	assert.NotEqual(t, len(code), 0)

	var bufWrapper bytes.Buffer
	err = generateContractWrapper(tmpDir+testContract, &bufWrapper)
	assert.NoError(t, err)
	assert.NotContains(t, bufWrapper.String(), `someAlias "some/test/import/path"`)
	assert.Contains(t, bufWrapper.String(), `"github.com/insolar/insolar/logicrunner/goplugin/foundation"`)
}

func TestNotMatchFileNameForProxy(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	testContract := "/test_not_go_file.test"
	err = testutil.WriteFile(tmpDir, testContract, `
package main

type A struct{
	foundation.BaseContract
}
`)

	var bufProxy bytes.Buffer
	err = generateContractProxy(tmpDir+testContract, "testRef", &bufProxy)
	assert.EqualError(t, err, "couldn't match filename without extension and path")
}
