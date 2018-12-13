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
	"path"
	"path/filepath"
	"testing"

	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
)

var randomTestCode = `
package main

import (
	"fmt"
	"errors"

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

func (hw *HelloWorlder) Hello() (string, error) {
	hw.Greeted++
	return "Hello world 2", nil
}

func (hw *HelloWorlder) Fail() (string, error) {
	hw.Greeted++
	return "", errors.New("We failed 2")
}

func (hw *HelloWorlder) Echo(s string) (string, error) {
	hw.Greeted++
	return s, nil
}

func (hw *HelloWorlder) HelloHuman(Name FullName) (PersonalGreeting, error) {
	hw.Greeted++
	return PersonalGreeting{
		Name:    Name,
		Message: fmt.Sprintf("Dear %s %s, we specially say hello to you", Name.First, Name.Last),
	}, nil
}

func (hw *HelloWorlder) HelloHumanPointer(Name FullName) (*PersonalGreeting, error) {
	hw.Greeted++
	return &PersonalGreeting{
		Name:    Name,
		Message: fmt.Sprintf("Dear %s %s, we specially say hello to you", Name.First, Name.Last),
	}, nil
}

func (hw *HelloWorlder) MultiArgs(Name FullName, s string, i int) (*PersonalGreeting, error) {
	hw.Greeted++
	return &PersonalGreeting{
		Name:    Name,
		Message: fmt.Sprintf("Dear %s %s, we specially say hello to you", Name.First, Name.Last),
	}, nil
}

func (hw HelloWorlder) ConstEcho(s string) (string, error) {
	return s, nil
}

func JustExportedStaticFunction(int, int) error { return nil }
`

func TestBasicGeneration(t *testing.T) {
	t.Parallel()
	tmpDir, err := ioutil.TempDir("", "test_")
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
		err := parsed.WriteWrapper(&buf)
		assert.NoError(t, err)

		code, err := ioutil.ReadAll(&buf)
		assert.NoError(t, err)
		assert.NotEmpty(t, code)
	})

	t.Run("proxy", func(t *testing.T) {
		t.Parallel()

		buf := bytes.Buffer{}
		err := parsed.WriteProxy(testutils.RandomRef().String(), &buf)
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

func New() (*One, error) {
	return &One{}, nil
}

func NewFromString(s string) (*One, error) {
	return &One{}, nil
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

	code = `
package main

type One struct {
	foundation.BaseContract
}

func New() {
	return
}
`

	err = goplugintestutils.WriteFile(tmpDir, "code1", code)
	assert.NoError(t, err)

	_, err = ParseFile(filepath.Join(tmpDir, "code1"))
	assert.Error(t, err)

	code = `
package main

type One struct {
	foundation.BaseContract
}

func New() *One {
	return &One{}
}
`

	err = goplugintestutils.WriteFile(tmpDir, "code1", code)
	assert.NoError(t, err)

	_, err = ParseFile(filepath.Join(tmpDir, "code1"))
	assert.Error(t, err)
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

	err = parsed.WriteProxy(testutils.RandomRef().String(), proxyFh)
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
	ref, _ := core.NewRefFromBase58("4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.7ZQboaH24PH42sqZKUvoa7UBrpuuubRtShp6CKNuWGZa")
	_ = secondary.GetObject(*ref)
}
	`)
	assert.NoError(t, err)

	cmd := exec.Command("go", "build", filepath.Join(tmpDir, "test.go"))
	cmd.Env = append(os.Environ(), "GOPATH="+goplugintestutils.PrependGoPath(tmpDir))
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err, string(out))
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

func ( a *A ) Get(
	a int, b bool, c string, d foundation.Reference,
) (
	int, bool, string, foundation.Reference, error,
) {
	return nil
}
`)

	assert.NoError(t, err)

	parsed, err := ParseFile(tmpDir + testContract)
	assert.NoError(t, err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy(testutils.RandomRef().String(), &bufProxy)
	assert.NoError(t, err)
	assert.Contains(t, bufProxy.String(), "var ret0 int")
	assert.Contains(t, bufProxy.String(), "ret[0] = &ret0")

	assert.Contains(t, bufProxy.String(), "var ret1 bool")
	assert.Contains(t, bufProxy.String(), "ret[1] = &ret1")

	assert.Contains(t, bufProxy.String(), "var ret2 string")
	assert.Contains(t, bufProxy.String(), "ret[2] = &ret2")

	assert.Contains(t, bufProxy.String(), "var ret3 foundation.Reference")
	assert.Contains(t, bufProxy.String(), "ret[3] = &ret3")
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

func (a *A) Get(
	a int, b bool, c string, d foundation.Reference,
) (
	int, bool, string, foundation.Reference, error,
) {
	return
}
`)
	assert.NoError(t, err)

	parsed, err := ParseFile(tmpDir + testContract)
	assert.NoError(t, err)

	var bufWrapper bytes.Buffer
	err = parsed.WriteWrapper(&bufWrapper)
	assert.NoError(t, err)
	assert.Contains(t, bufWrapper.String(), "var args0 int")
	assert.Contains(t, bufWrapper.String(), "args[0] = &args0")

	assert.Contains(t, bufWrapper.String(), "var args1 bool")
	assert.Contains(t, bufWrapper.String(), "args[1] = &args1")

	assert.Contains(t, bufWrapper.String(), "var args2 string")
	assert.Contains(t, bufWrapper.String(), "args[2] = &args2")

	assert.Contains(t, bufWrapper.String(), "var args3 foundation.Reference")
	assert.Contains(t, bufWrapper.String(), "args[3] = &args3")
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

func ( A ) Get(i path.SomeType) error {
	return nil
}

func ( A ) GetPointer(i *pointerPath.SomeType) error {
	return nil
}
`)
	assert.NoError(t, err)

	parsed, err := ParseFile(tmpDir + testContract)
	assert.NoError(t, err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy(testutils.RandomRef().String(), &bufProxy)
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

func ( A ) Get(i someAlias.SomeType) error {
	return nil
}
`)
	assert.NoError(t, err)

	parsed, err := ParseFile(tmpDir + testContract)
	assert.NoError(t, err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy(testutils.RandomRef().String(), &bufProxy)
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

func ( A ) Get() error {
	path.SomeMethod()
	return nil
}
`)
	assert.NoError(t, err)

	parsed, err := ParseFile(tmpDir + testContract)
	assert.NoError(t, err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy(testutils.RandomRef().String(), &bufProxy)
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

func ( A ) Get() (path.SomeValue, error) {
	f := path.SomeMethod()
	return f, nil
}
`)
	assert.NoError(t, err)

	parsed, err := ParseFile(tmpDir + testContract)
	assert.NoError(t, err)

	var bufProxy bytes.Buffer
	err = parsed.WriteProxy(testutils.RandomRef().String(), &bufProxy)
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
	err = parsed.WriteProxy(testutils.RandomRef().String(), &bufProxy)
	assert.EqualError(t, err, "couldn't match filename without extension and path")
}

func TestProxyGeneration(t *testing.T) {
	contracts, err := GetRealContractsNames()
	assert.NoError(t, err)

	for _, contract := range contracts {
		// Make a copy for proper work of closure inside gorutine
		contract := contract

		t.Run(contract, func(t *testing.T) {
			t.Parallel()
			parsed, err := ParseFile("../../../application/contract/" + contract + "/" + contract + ".go")
			assert.NotNil(t, parsed, "have parsed object")
			assert.NoError(t, err)

			proxyPath, err := GetRealApplicationDir("proxy")
			assert.NoError(t, err)

			name, err := parsed.ProxyPackageName()
			assert.NoError(t, err)

			proxy := path.Join(proxyPath, name, name+".go")
			_, err = os.Stat(proxy)
			assert.NoError(t, err)

			buff := bytes.NewBufferString("")
			err = parsed.WriteProxy("", buff)
			require.NoError(t, err)

			cmd := exec.Command("diff", "-u", proxy, "-")
			cmd.Stdin = buff
			out, err := cmd.CombinedOutput()
			assert.NoError(t, err, string(out))
		})
	}
}
