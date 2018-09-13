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
	"github.com/insolar/insolar/ledger/ledgertestutil"
	"github.com/insolar/insolar/logicrunner/goplugin"
	"github.com/insolar/insolar/logicrunner/goplugin/testutil"
	"github.com/insolar/insolar/messagerouter/message"
)

var icc = ""
var runnerbin = ""

func TestMain(m *testing.M) {
	var err error
	log.SetLevel(log.DebugLevel)
	if runnerbin, icc, err = testutil.Build(); err != nil {
		fmt.Println("Logic runner build failed, skip tests:", err.Error())
		os.Exit(1)
	}
	os.Exit(m.Run())
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

func (r *testExecutor) CallMethod(ctx *core.LogicCallContext, code core.RecordRef, data []byte, method string, args core.Arguments) ([]byte, core.Arguments, error) {
	if len(r.methodResponses) < 1 {
		panic(errors.New("no expected 'CallMethod' calls"))
	}

	res := r.methodResponses[0]
	r.methodResponses = r.methodResponses[1:]
	return res.data, res.res, res.err
}

func (r *testExecutor) CallConstructor(ctx *core.LogicCallContext, code core.RecordRef, name string, args core.Arguments) ([]byte, error) {
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

func (r *testLedger) GetPulseManager() core.PulseManager {
	panic("implement me")
}

func (r *testLedger) GetJetCoordinator() core.JetCoordinator {
	panic("implement me")
}

func (r *testLedger) Start(components core.Components) error   { return nil }
func (r *testLedger) Stop() error                              { return nil }
func (r *testLedger) GetArtifactManager() core.ArtifactManager { return r.am }

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
		AM:    am,
		Data:  []byte("origData"),
		Code:  &codeRef,
		Class: &classRef,
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

func TestContractCallingContract(t *testing.T) {
	var contractOneCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import "contract-proxy/two"

type One struct {
	foundation.BaseContract
}

func (r *One) Hello(s string) string {
	holder := two.New()
	friend := holder.AsChild(r.GetReference())

	res := friend.Hello(s)

	return "Hi, " + s + "! Two said: " + res
}
`

	var contractTwoCode = `
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

	lr, err := NewLogicRunner(configuration.LogicRunner{})
	assert.NoError(t, err)

	mr := &testMessageRouter{LogicRunner: lr}
	am := testutil.NewTestArtifactManager()
	lr.ArtifactManager = am

	insiderStorage, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(insiderStorage) // nolint: errcheck

	gp, err := goplugin.NewGoPlugin(
		&configuration.GoPlugin{
			MainListen:     "127.0.0.1:7778",
			RunnerListen:   "127.0.0.1:7777",
			RunnerPath:     runnerbin,
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
	assert.NoError(t, err)

	var argsSerialized []byte
	err = codec.NewEncoderBytes(&argsSerialized, ch).Encode(
		[]interface{}{"ins"},
	)
	assert.NoError(t, err)

	cb, cleaner := testutil.NewContractBuilder(am, icc)
	defer cleaner()
	err = cb.Build(map[string]string{"one": contractOneCode, "two": contractTwoCode})
	assert.NoError(t, err)

	obj, err := am.ActivateObj(
		core.RecordRef{}, core.RecordRef{},
		*cb.Classes["one"],
		*am.RootRef(),
		data,
	)
	assert.NoError(t, err)

	_, res, err := gp.CallMethod(
		&core.LogicCallContext{Class: cb.Classes["one"], Callee: obj}, *cb.Codes["one"],
		data, "Hello", argsSerialized,
	)
	assert.NoError(t, err)

	var resParsed []interface{}
	err = codec.NewDecoderBytes(res, ch).Decode(&resParsed)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Hi, ins! Two said: Hello you too, ins. 644 times!", resParsed[0])
}

func TestInjectingDelegate(t *testing.T) {
	var contractOneCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import "contract-proxy/two"

type One struct {
	foundation.BaseContract
}

func (r *One) Hello(s string) string {
	holder := two.New()
	friend := holder.AsDelegate(r.GetReference())

	res := friend.Hello(s)

	return "Hi, " + s + "! Two said: " + res
}

func (r *One) HelloFromDelegate(s string) string {
	friend := two.GetImplementationFrom(r.GetReference())
	return friend.Hello(s)
}
`

	var contractTwoCode = `
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

	lr, err := NewLogicRunner(configuration.LogicRunner{})
	assert.NoError(t, err)

	mr := &testMessageRouter{LogicRunner: lr}
	am := testutil.NewTestArtifactManager()
	lr.ArtifactManager = am

	insiderStorage, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(insiderStorage) // nolint: errcheck

	gp, err := goplugin.NewGoPlugin(
		&configuration.GoPlugin{
			MainListen:     "127.0.0.1:7778",
			RunnerListen:   "127.0.0.1:7777",
			RunnerPath:     runnerbin,
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

	cb, cleaner := testutil.NewContractBuilder(am, icc)
	defer cleaner()
	err = cb.Build(map[string]string{"one": contractOneCode, "two": contractTwoCode})
	assert.NoError(t, err)

	obj, err := am.ActivateObj(
		core.RecordRef{}, core.RecordRef{},
		*cb.Classes["one"],
		*am.RootRef(),
		data,
	)
	assert.NoError(t, err)

	_, res, err := gp.CallMethod(
		&core.LogicCallContext{Class: cb.Classes["one"], Callee: obj}, *cb.Codes["one"],
		data, "Hello", argsSerialized,
	)
	assert.NoError(t, err)

	var resParsed []interface{}
	err = codec.NewDecoderBytes(res, ch).Decode(&resParsed)
	assert.NoError(t, err)
	assert.Equal(t, "Hi, ins! Two said: Hello you too, ins. 644 times!", resParsed[0])

	_, res, err = gp.CallMethod(
		&core.LogicCallContext{Class: cb.Classes["one"], Callee: obj}, *cb.Codes["one"],
		data, "HelloFromDelegate", argsSerialized,
	)
	assert.NoError(t, err)

	err = codec.NewDecoderBytes(res, ch).Decode(&resParsed)
	assert.NoError(t, err)
	assert.Equal(t, "Hello you too, ins. 1288 times!", resParsed[0])
}

func TestContextPassing(t *testing.T) {
	var code = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type One struct {
	foundation.BaseContract
}

func (r *One) Hello() string {
	return r.GetClass().String()
}
`

	am := testutil.NewTestArtifactManager()

	insiderStorage, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(insiderStorage) // nolint: errcheck

	gp, err := goplugin.NewGoPlugin(
		&configuration.GoPlugin{
			MainListen:     "127.0.0.1:7778",
			RunnerListen:   "127.0.0.1:7777",
			RunnerPath:     runnerbin,
			RunnerCodePath: insiderStorage,
		},
		nil,
		am,
	)
	assert.NoError(t, err)
	defer gp.Stop()

	ch := new(codec.CborHandle)
	var data []byte
	err = codec.NewEncoderBytes(&data, ch).Encode(&struct{}{})
	assert.NoError(t, err)

	var argsSerialized []byte
	err = codec.NewEncoderBytes(&argsSerialized, ch).Encode([]interface{}{})
	assert.NoError(t, err)

	cb, cleaner := testutil.NewContractBuilder(am, icc)
	defer cleaner()
	err = cb.Build(map[string]string{"one": code})
	assert.NoError(t, err)

	_, res, err := gp.CallMethod(
		&core.LogicCallContext{Class: cb.Classes["one"]}, *cb.Codes["one"],
		data, "Hello", argsSerialized,
	)
	assert.NoError(t, err)

	resParsed := []interface{}{""}
	err = codec.NewDecoderBytes(res, ch).Decode(&resParsed)
	assert.NoError(t, err)
	assert.Equal(t, cb.Classes["one"].String(), resParsed[0])
}

func TestGetChildren(t *testing.T) {
	goContract := `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
//	"github.com/insolar/insolar/core"
	"contract-proxy/child"
)

type Contract struct {
	foundation.BaseContract
}

func (c *Contract) NewChilds(cnt int) int {
	s := 0
	for i := 1; i < cnt; i++ {
        child.New(i).AsChild(c.GetReference())
		s += i
	} 
	return s
}

func (c *Contract) SumChilds() int {
	s := 0
	childs, err := c.GetChildrenTyped(child.GetClass())
	if err != nil {
		panic(err)
	}
	for _, chref := range childs {
		o := child.GetObject(chref)
		s += o.GetNum()
	}
	return s
}

func (c *Contract) GetChildRefs() (ret []string) {
	childs, err := c.GetChildrenTyped(child.GetClass())
	if err != nil {
		panic(err)
	}

	for _, chref := range childs {
		ret = append(ret, chref.String())
	}
	return ret
}
`
	goChild := `
package main
import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type Child struct {
	foundation.BaseContract
	Num int
}

func (c *Child) GetNum() int {
	return c.Num
}


func New(n int) *Child {
	return &Child{Num: n};
}
`
	l, cleaner := ledgertestutil.TmpLedger(t, "")
	defer cleaner()

	insiderStorage, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(insiderStorage) // nolint: errcheck

	am := l.GetArtifactManager()
	lr, err := NewLogicRunner(configuration.LogicRunner{
		GoPlugin: &configuration.GoPlugin{
			MainListen:     "127.0.0.1:7778",
			RunnerListen:   "127.0.0.1:7777",
			RunnerPath:     runnerbin,
			RunnerCodePath: insiderStorage,
		}})
	assert.NoError(t, err, "Initialize runner")

	assert.NoError(t, lr.Start(core.Components{
		"core.Ledger":        l,
		"core.MessageRouter": &testMessageRouter{LogicRunner: lr},
	}), "starting logicrunner")
	defer lr.Stop()

	cb, cleaner := testutil.NewContractBuilder(am, icc)
	//	defer cleaner()
	err = cb.Build(map[string]string{"child": goChild})
	assert.NoError(t, err)
	err = cb.Build(map[string]string{"contract": goContract})
	assert.NoError(t, err)

	domain := core.String2Ref("c1")
	contract, err := am.ActivateObj(core.String2Ref("r1"), domain, *cb.Classes["contract"], *am.RootRef(), testutil.CBORMarshal(t, nil))
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, contract, nil, "contract created")

	resp := lr.Execute(&message.CallMethodMessage{
		Request:   core.String2Ref("r2"),
		ObjectRef: *contract,
		Method:    "NewChilds",
		Arguments: testutil.CBORMarshal(t, []interface{}{100}),
	})
	assert.NoError(t, resp.Error, "contract call")
	r := testutil.CBORUnMarshal(t, resp.Result)
	assert.Equal(t, []interface{}([]interface{}{uint64(4950)}), r)

	resp = lr.Execute(&message.CallMethodMessage{
		Request:   core.String2Ref("r3"),
		ObjectRef: *contract,
		Method:    "SumChilds",
		Arguments: testutil.CBORMarshal(t, []interface{}{}),
	})
	assert.NoError(t, resp.Error, "contract call")
	r = testutil.CBORUnMarshal(t, resp.Result)
	assert.Equal(t, []interface{}([]interface{}{uint64(4950)}), r)

}
