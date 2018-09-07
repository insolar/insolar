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

package logicrunner

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/ugorji/go/codec"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin"
	"github.com/insolar/insolar/logicrunner/goplugin/testutil"
	"github.com/insolar/insolar/messagerouter/message"
)

var icc = "../cmd/icc/icc"

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestTypeCompatibility(t *testing.T) {
	var _ core.LogicRunner = (*LogicRunner)(nil)
}

type testExecutor struct {
	constructorResponses []*testResp
	methodResponses      []*testResp
}

func (r *testExecutor) Stop() error {
	return nil
}

type testResp struct {
	data []byte
	res  core.Arguments
	err  error
}

func newTestExecutor() *testExecutor {
	return &testExecutor{
		constructorResponses: make([]*testResp, 0),
		methodResponses:      make([]*testResp, 0),
	}
}

func (r *testExecutor) CallMethod(ref core.RecordRef, data []byte, method string, args core.Arguments) ([]byte, core.Arguments, error) {
	if len(r.methodResponses) < 1 {
		panic(errors.New("no expected 'CallMethod' calls"))
	}

	res := r.methodResponses[0]
	r.methodResponses = r.methodResponses[1:]
	return res.data, res.res, res.err
}

func (r *testExecutor) CallConstructor(ref core.RecordRef, name string, args core.Arguments) ([]byte, error) {
	if len(r.constructorResponses) < 1 {
		panic(errors.New("no expected 'CallConstructor' calls"))
	}

	res := r.constructorResponses[0]
	r.constructorResponses = r.constructorResponses[1:]
	return res.data, res.err
}

func TestBasics(t *testing.T) {
	lr, err := NewLogicRunner(configuration.LogicRunner{})
	assert.NoError(t, err)

	comps := core.Components{
		"core.Ledger":        &testLedger{am: testutil.NewTestArtifactManager()},
		"core.MessageRouter": &testMessageRouter{},
	}
	assert.NoError(t, lr.Start(comps))
	assert.IsType(t, &LogicRunner{}, lr)

	_, err = lr.GetExecutor(core.MachineTypeGoPlugin)
	assert.Error(t, err)

	te := newTestExecutor()

	err = lr.RegisterExecutor(core.MachineTypeGoPlugin, te)
	assert.NoError(t, err)

	te2, err := lr.GetExecutor(core.MachineTypeGoPlugin)
	assert.NoError(t, err)
	assert.Equal(t, te, te2)
}

type testLedger struct {
	am core.ArtifactManager
}

func (r *testLedger) Start(components core.Components) error { return nil }
func (r *testLedger) Stop() error                            { return nil }
func (r *testLedger) GetManager() core.ArtifactManager       { return r.am }

type testMessageRouter struct {
	LogicRunner core.LogicRunner
}

func (*testMessageRouter) Start(components core.Components) error { return nil }
func (*testMessageRouter) Stop() error                            { return nil }
func (r *testMessageRouter) Route(msg core.Message) (resp core.Response, err error) {
	res := r.LogicRunner.Execute(msg)
	return *res, nil
}

func TestExecution(t *testing.T) {
	am := testutil.NewTestArtifactManager()
	ld := &testLedger{am: am}
	mr := &testMessageRouter{}
	lr, err := NewLogicRunner(configuration.LogicRunner{})
	assert.NoError(t, err)
	lr.Start(core.Components{
		"core.Ledger":        ld,
		"core.MessageRouter": mr,
	})

	codeRef := core.String2Ref("someCode")
	dataRef := core.String2Ref("someObject")
	classRef := core.String2Ref("someClass")
	am.Objects[dataRef] = &testutil.TestObjectDescriptor{
		AM:   am,
		Data: []byte("origData"),
		Code: &codeRef,
	}
	am.Classes[classRef] = &testutil.TestClassDescriptor{AM: am, ARef: &classRef, ACode: &codeRef}
	am.Codes[codeRef] = &testutil.TestCodeDescriptor{ARef: &codeRef, AMachineType: core.MachineTypeGoPlugin}

	te := newTestExecutor()
	te.methodResponses = append(te.methodResponses, &testResp{data: []byte("data"), res: core.Arguments("res")})

	err = lr.RegisterExecutor(core.MachineTypeGoPlugin, te)
	assert.NoError(t, err)

	resp := lr.Execute(&message.CallMethodMessage{ObjectRef: dataRef})
	assert.NoError(t, resp.Error)
	assert.Equal(t, []byte("data"), resp.Data)
	assert.Equal(t, []byte("res"), resp.Result)

	te.constructorResponses = append(te.constructorResponses, &testResp{data: []byte("data"), res: core.Arguments("res")})
	resp = lr.Execute(&message.CallConstructorMessage{ClassRef: classRef})
	assert.NoError(t, resp.Error)
	assert.Equal(t, []byte("data"), resp.Data)
	assert.Equal(t, []byte(nil), resp.Result)
}

func buildCLI(name string) error {
	out, err := exec.Command("go", "build", "-o", "./goplugin/"+name+"/"+name, "./goplugin/"+name+"/").CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "can't build %s: %s", name, string(out))
	}
	return nil
}

func buildInciderCLI() error {
	return buildCLI("ginsider-cli")
}

func buildPreprocessor() error {
	out, err := exec.Command("go", "build", "-o", icc, "../cmd/icc/").CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "can't build %s: %s", icc, string(out))
	}
	return nil

}

const contractOneCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import "contract-proxy/two"

type One struct {
	foundation.BaseContract
}

func (r *One) Hello(s string) string {
	holder := two.New()
	friend := holder.AsChild("")

	res := friend.Hello(s)

	return "Hi, " + s + "! Two said: " + res
}
`

const contractTwoCode = `
package main

import (
	"fmt"

	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Two struct {
	foundation.BaseContract
	X int
}

func New() *Two {
	return &Two{X:322};
}

func (r *Two) Hello(s string) string {
	r.X *= 2
	return fmt.Sprintf("Hello you too, %s. %d times!", s, r.X)
}
`

func generateContractProxy(root string, name string) error {
	dstDir := root + "/src/contract-proxy/" + name

	err := os.MkdirAll(dstDir, 0777)
	if err != nil {
		return err
	}

	contractPath := root + "/src/contract/" + name + "/main.go"

	out, err := exec.Command(icc, "proxy", "-o", dstDir+"/main.go", "--code-reference", "Class"+name, contractPath).CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "can't generate proxy: "+string(out))
	}
	return nil
}

func buildContractPlugin(root string, name string) error {
	dstDir := root + "/plugins/"

	err := os.MkdirAll(dstDir, 0777)
	if err != nil {
		return err
	}

	origGoPath, err := testutil.ChangeGoPath(root)
	if err != nil {
		return err
	}
	defer os.Setenv("GOPATH", origGoPath) // nolint: errcheck

	//contractPath := root + "/src/contract/" + name + "/main.go"

	out, err := exec.Command("go", "build", "-buildmode=plugin", "-o", dstDir+"/"+name+".so", "contract/"+name).CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "can't build contract: "+string(out))
	}
	return nil
}

func generateContractWrapper(root string, name string) error {
	contractPath := root + "/src/contract/" + name + "/main.go"
	wrapperPath := root + "/src/contract/" + name + "/main_wrapper.go"

	out, err := exec.Command(icc, "wrapper", "-o", wrapperPath, contractPath).CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "can't generate wrapper for contract '"+name+"': "+string(out))
	}
	return nil
}

func buildContracts(root string, names ...string) error {
	for _, name := range names {
		err := generateContractProxy(root, name)
		if err != nil {
			return err
		}
		err = generateContractWrapper(root, name)
		if err != nil {
			return err
		}
	}

	for _, name := range names {
		err := buildContractPlugin(root, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func suckInContracts(am *testutil.TestArtifactManager, root string, names ...string) {
	for _, name := range names {
		pluginBinary, err := ioutil.ReadFile(root + "/plugins/" + name + ".so")
		if err != nil {
			panic(err)
		}

		ref := core.String2Ref(name)
		am.Codes[ref] = &testutil.TestCodeDescriptor{
			ARef:         &ref,
			ACode:        pluginBinary,
			AMachineType: core.MachineTypeGoPlugin,
		}
	}
}

func TestContractCallingContract(t *testing.T) {
	err := buildInciderCLI()
	if err != nil {
		t.Fatal(err)
	}

	err = buildPreprocessor()
	if err != nil {
		t.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd) // nolint: errcheck

	tmpDir, err := ioutil.TempDir("", "test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // nolint: errcheck

	err = testutil.WriteFile(tmpDir+"/src/contract/one/", "main.go", contractOneCode)
	if err != nil {
		t.Fatal(err)
	}
	err = testutil.WriteFile(tmpDir+"/src/contract/two/", "main.go", contractTwoCode)
	if err != nil {
		t.Fatal(err)
	}

	err = buildContracts(tmpDir, "one", "two")
	if err != nil {
		t.Fatal(err)
	}

	insiderStorage := tmpDir + "/insider-storage/"

	err = os.MkdirAll(insiderStorage, 0777)
	if err != nil {
		t.Fatal(err)
	}

	lr, err := NewLogicRunner(configuration.LogicRunner{})
	assert.NoError(t, err)

	mr := &testMessageRouter{LogicRunner: lr}
	am := testutil.NewTestArtifactManager()
	lr.ArtifactManager = am

	gp, err := goplugin.NewGoPlugin(
		&configuration.GoPlugin{
			MainListen:     "127.0.0.1:7778",
			RunnerListen:   "127.0.0.1:7777",
			RunnerPath:     "./goplugin/ginsider-cli/ginsider-cli",
			RunnerCodePath: insiderStorage,
		},
		mr,
		am,
	)
	assert.NoError(t, err)
	defer gp.Stop()

	err = lr.RegisterExecutor(core.MachineTypeGoPlugin, gp)
	assert.NoError(t, err)

	ch := new(codec.CborHandle)
	var data []byte
	err = codec.NewEncoderBytes(&data, ch).Encode(
		&struct{}{},
	)
	if err != nil {
		t.Fatal(err)
	}

	var argsSerialized []byte
	err = codec.NewEncoderBytes(&argsSerialized, ch).Encode(
		[]interface{}{"ins"},
	)
	if err != nil {
		panic(err)
	}

	suckInContracts(am, tmpDir, "one", "two")

	codeRef := core.String2Ref("two")
	am.Classes[core.String2Ref("Classtwo")] = &testutil.TestClassDescriptor{
		AM:    am,
		ACode: &codeRef,
	}

	_, res, err := gp.CallMethod(core.String2Ref("one"), data, "Hello", argsSerialized)
	if err != nil {
		panic(err)
	}

	var resParsed []interface{}
	err = codec.NewDecoderBytes(res, ch).Decode(&resParsed)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Hi, ins! Two said: Hello you too, ins. 644 times!", resParsed[0])
}
