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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/insolar/insolar/logicrunner/goplugin/preprocessor"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/ugorji/go/codec"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/eventbus/event"
	"github.com/insolar/insolar/eventbus/reaction"
	"github.com/insolar/insolar/ledger/ledgertestutil"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/goplugin"
	"github.com/insolar/insolar/logicrunner/goplugin/testutil"
)

var icc = ""
var runnerbin = ""

func TestMain(m *testing.M) {
	var err error
	err = log.SetLevel("Debug")
	if err != nil {
		log.Errorln(err.Error())
	}
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
		"core.Ledger":   &testLedger{am: testutil.NewTestArtifactManager()},
		"core.EventBus": &testEventBus{},
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

type testEventBus struct {
	LogicRunner core.LogicRunner
}

func (*testEventBus) Start(components core.Components) error { return nil }
func (*testEventBus) Stop() error                            { return nil }
func (eb *testEventBus) Dispatch(event core.Event) (resp core.Reaction, err error) {
	return eb.LogicRunner.Execute(event)
}
func (*testEventBus) DispatchAsync(event core.Event) {}

func TestExecution(t *testing.T) {
	am := testutil.NewTestArtifactManager()
	ld := &testLedger{am: am}
	eb := &testEventBus{}
	lr, err := NewLogicRunner(configuration.LogicRunner{})
	assert.NoError(t, err)
	lr.Start(core.Components{
		"core.Ledger":   ld,
		"core.EventBus": eb,
	})
	eb.LogicRunner = lr

	codeRef := core.NewRefFromBase58("someCode")
	dataRef := core.NewRefFromBase58("someObject")
	classRef := core.NewRefFromBase58("someClass")
	am.Objects[dataRef] = &testutil.TestObjectDescriptor{
		AM:    am,
		Data:  []byte("origData"),
		Code:  &codeRef,
		Class: &classRef,
	}
	am.Classes[classRef] = &testutil.TestClassDescriptor{AM: am, ARef: &classRef, ACode: &codeRef}
	am.Codes[codeRef] = &testutil.TestCodeDescriptor{ARef: codeRef, AMachineType: core.MachineTypeGoPlugin}

	te := newTestExecutor()
	te.methodResponses = append(te.methodResponses, &testResp{data: []byte("data"), res: core.Arguments("res")})

	err = lr.RegisterExecutor(core.MachineTypeGoPlugin, te)
	assert.NoError(t, err)

	resp, err := lr.Execute(&event.CallMethodEvent{ObjectRef: dataRef})
	assert.NoError(t, err)
	assert.Equal(t, []byte("data"), resp.(*reaction.CommonReaction).Data)
	assert.Equal(t, []byte("res"), resp.(*reaction.CommonReaction).Result)

	te.constructorResponses = append(te.constructorResponses, &testResp{data: []byte("data"), res: core.Arguments("res")})
	resp, err = lr.Execute(&event.CallConstructorEvent{ClassRef: classRef})
	assert.NoError(t, err)
	assert.Equal(t, []byte("data"), resp.(*reaction.CommonReaction).Data)
	assert.Equal(t, []byte(nil), resp.(*reaction.CommonReaction).Result)
}

func TestContractCallingContract(t *testing.T) {
	var contractOneCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import "github.com/insolar/insolar/genesis/proxy/two"

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

	eb := &testEventBus{LogicRunner: lr}
	lr.EventBus = eb
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
		eb,
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

	cb := testutil.NewContractBuilder(am, icc)
	defer cb.Clean()
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
import "github.com/insolar/insolar/genesis/proxy/two"

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

	eb := &testEventBus{LogicRunner: lr}
	lr.EventBus = eb
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
		eb,
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

	cb := testutil.NewContractBuilder(am, icc)
	defer cb.Clean()
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

func TestBasicNotificationCall(t *testing.T) {
	var contractOneCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import "github.com/insolar/insolar/genesis/proxy/two"

type One struct {
	foundation.BaseContract
}

func (r *One) Hello() {
	holder := two.New()
	friend := holder.AsDelegate(r.GetReference())
	friend.HelloNoWait()
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

func (r *Two) Hello() string {
	r.X *= 2
	return fmt.Sprintf("Hello %d times!", r.X)
}
`

	lr, err := NewLogicRunner(configuration.LogicRunner{})
	assert.NoError(t, err)

	eb := &testEventBus{LogicRunner: lr}
	lr.EventBus = eb
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
		eb,
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
		[]interface{}{},
	)
	if err != nil {
		panic(err)
	}

	cb := testutil.NewContractBuilder(am, icc)
	defer cb.Clean()
	err = cb.Build(map[string]string{"one": contractOneCode, "two": contractTwoCode})
	assert.NoError(t, err)

	obj, err := am.ActivateObj(
		core.RecordRef{}, core.RecordRef{},
		*cb.Classes["one"],
		*am.RootRef(),
		data,
	)
	assert.NoError(t, err)

	_, _, err = gp.CallMethod(
		&core.LogicCallContext{Class: cb.Classes["one"], Callee: obj}, *cb.Codes["one"],
		data, "Hello", argsSerialized,
	)
	assert.NoError(t, err)
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

	cb := testutil.NewContractBuilder(am, icc)
	defer cb.Clean()
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
	"github.com/insolar/insolar/genesis/proxy/child"
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
	am.SetArchPref([]core.MachineType{core.MachineTypeGoPlugin})
	lr, err := NewLogicRunner(configuration.LogicRunner{
		GoPlugin: &configuration.GoPlugin{
			MainListen:     "127.0.0.1:7778",
			RunnerListen:   "127.0.0.1:7777",
			RunnerPath:     runnerbin,
			RunnerCodePath: insiderStorage,
		}})
	assert.NoError(t, err, "Initialize runner")

	assert.NoError(t, lr.Start(core.Components{
		"core.Ledger":   l,
		"core.EventBus": &testEventBus{LogicRunner: lr},
	}), "starting logicrunner")
	defer lr.Stop()

	cb := testutil.NewContractBuilder(am, icc)
	defer cb.Clean()
	err = cb.Build(map[string]string{"child": goChild})
	assert.NoError(t, err)
	err = cb.Build(map[string]string{"contract": goContract})
	assert.NoError(t, err)

	domain := core.NewRefFromBase58("c1")
	contract, err := am.ActivateObj(core.NewRefFromBase58("r1"), domain, *cb.Classes["contract"], *am.RootRef(), testutil.CBORMarshal(t, nil))
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, contract, nil, "contract created")

	resp, err := lr.Execute(&event.CallMethodEvent{
		Request:   core.NewRefFromBase58("r2"),
		ObjectRef: *contract,
		Method:    "NewChilds",
		Arguments: testutil.CBORMarshal(t, []interface{}{10}),
	})
	assert.NoError(t, err, "contract call")
	r := testutil.CBORUnMarshal(t, resp.(*reaction.CommonReaction).Result)
	assert.Equal(t, []interface{}([]interface{}{uint64(45)}), r)

	resp, err = lr.Execute(&event.CallMethodEvent{
		Request:   core.NewRefFromBase58("r3"),
		ObjectRef: *contract,
		Method:    "SumChilds",
		Arguments: testutil.CBORMarshal(t, []interface{}{}),
	})
	assert.NoError(t, err, "contract call")
	r = testutil.CBORUnMarshal(t, resp.(*reaction.CommonReaction).Result)
	assert.Equal(t, []interface{}([]interface{}{uint64(45)}), r)

}

func TestRootDomainContract(t *testing.T) {
	rootDomainCode, err := ioutil.ReadFile("../genesis/experiment/rootdomain/rootdomain.go" +
		"")
	if err != nil {
		fmt.Print(err)
	}
	memberCode, err := ioutil.ReadFile("../genesis/experiment/member/member.go")
	if err != nil {
		fmt.Print(err)
	}
	allowanceCode, err := ioutil.ReadFile("../genesis/experiment/allowance/allowance.go")
	if err != nil {
		fmt.Print(err)
	}
	walletCode, err := ioutil.ReadFile("../genesis/experiment/wallet/wallet.go")
	if err != nil {
		fmt.Print(err)
	}

	l, cleaner := ledgertestutil.TmpLedger(t, "")
	defer cleaner()

	insiderStorage, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(insiderStorage) // nolint: errcheck

	am := l.GetArtifactManager()
	am.SetArchPref([]core.MachineType{core.MachineTypeGoPlugin})
	fmt.Println("RUNNERPATH", runnerbin)
	lr, err := NewLogicRunner(configuration.LogicRunner{
		GoPlugin: &configuration.GoPlugin{
			MainListen:     "127.0.0.1:7778",
			RunnerListen:   "127.0.0.1:7777",
			RunnerPath:     runnerbin,
			RunnerCodePath: insiderStorage,
		}})
	assert.NoError(t, err, "Initialize runner")

	assert.NoError(t, lr.Start(core.Components{
		"core.Ledger":   l,
		"core.EventBus": &testEventBus{LogicRunner: lr},
	}), "starting logicrunner")
	defer lr.Stop()

	cb := testutil.NewContractBuilder(am, icc)
	defer cb.Clean()
	err = cb.Build(map[string]string{"member": string(memberCode), "allowance": string(allowanceCode), "wallet": string(walletCode), "rootDomain": string(rootDomainCode)})
	assert.NoError(t, err)

	domain := core.NewRefFromBase58("c1")
	request := core.NewRefFromBase58("c2")
	contract, err := am.ActivateObj(domain, request, *cb.Classes["rootDomain"], *am.RootRef(), testutil.CBORMarshal(t, nil))
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, contract, nil, "contract created")

	resp1, err := lr.Execute(&event.CallMethodEvent{
		Request:   request,
		ObjectRef: *contract,
		Method:    "CreateMember",
		Arguments: testutil.CBORMarshal(t, []interface{}{"member1"}),
	})
	assert.NoError(t, err, "contract call")
	r1 := testutil.CBORUnMarshal(t, resp1.(*reaction.CommonReaction).Result)
	member1Ref := r1.([]interface{})[0].(string)

	resp2, err := lr.Execute(&event.CallMethodEvent{
		Request:   request,
		ObjectRef: *contract,
		Method:    "CreateMember",
		Arguments: testutil.CBORMarshal(t, []interface{}{"member2"}),
	})
	assert.NoError(t, err, "contract call")
	r2 := testutil.CBORUnMarshal(t, resp2.(*reaction.CommonReaction).Result)
	member2Ref := r2.([]interface{})[0].(string)

	_, err = lr.Execute(&event.CallMethodEvent{
		Request:   request,
		ObjectRef: *contract,
		Method:    "SendMoney",
		Arguments: testutil.CBORMarshal(t, []interface{}{member1Ref, member2Ref, 1}),
	})
	assert.NoError(t, err, "contract call")

	resp4, err := lr.Execute(&event.CallMethodEvent{
		Request:   request,
		ObjectRef: *contract,
		Method:    "DumpAllUsers",
		Arguments: testutil.CBORMarshal(t, []interface{}{}),
	})
	assert.NoError(t, err, "contract call")
	r := testutil.CBORUnMarshal(t, resp4.(*reaction.CommonReaction).Result)

	var res []map[string]interface{}
	var expected = map[interface{}]float64{"member1": 999, "member2": 1001}

	err = json.Unmarshal(r.([]interface{})[0].([]byte), &res)
	assert.NoError(t, err)
	for _, member := range res {
		assert.Equal(t, expected[member["member"]], member["wallet"])
	}
}

func BenchmarkContractCall(b *testing.B) {
	goParent := `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import "github.com/insolar/insolar/genesis/proxy/child"
import "github.com/insolar/insolar/core"

type Parent struct {
	foundation.BaseContract
}

func (c *Parent) CCC(ref *core.RecordRef) int {	
	o := child.GetObject(*ref)	
	return o.GetNum()
}
`
	goChild := `
package main
import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type Child struct {
	foundation.BaseContract
}

func (c *Child) GetNum() int {
	return 5
}
`
	l, cleaner := ledgertestutil.TmpLedger(b, "")
	defer cleaner()

	insiderStorage, err := ioutil.TempDir("", "test-")
	assert.NoError(b, err)
	defer os.RemoveAll(insiderStorage) // nolint: errcheck

	am := l.GetArtifactManager()
	lr, err := NewLogicRunner(configuration.LogicRunner{
		GoPlugin: &configuration.GoPlugin{
			MainListen:     "127.0.0.1:7778",
			RunnerListen:   "127.0.0.1:7777",
			RunnerPath:     runnerbin,
			RunnerCodePath: insiderStorage,
		}})
	assert.NoError(b, err, "Initialize runner")

	assert.NoError(b, lr.Start(core.Components{
		"core.Ledger":   l,
		"core.EventBus": &testEventBus{LogicRunner: lr},
	}), "starting logicrunner")
	defer lr.Stop()

	cb := testutil.NewContractBuilder(am, icc)
	defer cb.Clean()
	err = cb.Build(map[string]string{"child": goChild, "parent": goParent})
	assert.NoError(b, err)

	domain := core.NewRefFromBase58("c1")
	parent, err := am.ActivateObj(core.NewRefFromBase58("r1"), domain, *cb.Classes["parent"], *am.RootRef(), testutil.CBORMarshal(b, nil))
	assert.NoError(b, err, "create parent")
	assert.NotEqual(b, parent, nil, "parent created")
	child, err := am.ActivateObj(core.NewRefFromBase58("r2"), domain, *cb.Classes["child"], *am.RootRef(), testutil.CBORMarshal(b, nil))
	assert.NoError(b, err, "create child")
	assert.NotEqual(b, child, nil, "child created")

	b.N = 1000
	for i := 0; i < b.N; i++ {
		resp, err := lr.Execute(&event.CallMethodEvent{
			Request:   core.NewRefFromBase58("rr"),
			ObjectRef: *parent,
			Method:    "CCC",
			Arguments: testutil.CBORMarshal(b, []interface{}{child}),
		})
		assert.NoError(b, err, "parent call")
		r := testutil.CBORUnMarshal(b, resp.(*reaction.CommonReaction).Result)
		assert.Equal(b, []interface{}([]interface{}{uint64(5)}), r)
	}
}

func TestProxyGeneration(t *testing.T) {
	contracts, err := preprocessor.GetRealContractsNames()
	assert.NoError(t, err)

	for _, contract := range contracts {
		t.Run(contract, func(t *testing.T) {
			parsed, err := preprocessor.ParseFile("../genesis/experiment/" + contract + "/" + contract + ".go")
			assert.NoError(t, err)

			proxyPath, err := preprocessor.GetRealGenesisDir("proxy")
			assert.NoError(t, err)

			name, err := preprocessor.ProxyPackageName(parsed)
			assert.NoError(t, err)

			proxy := path.Join(proxyPath, name, name+".go")
			_, err = os.Stat(proxy)
			assert.NoError(t, err)

			buff := bytes.NewBufferString("")
			preprocessor.GenerateContractProxy(parsed, "", buff)

			cmd := exec.Command("diff", proxy, "-")
			cmd.Stdin = buff
			out, err := cmd.CombinedOutput()
			assert.NoError(t, err, string(out))
		})
	}
}
