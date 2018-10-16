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
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/insolar/insolar/application/contract/member"
	"github.com/insolar/insolar/application/contract/member/signer"
	"github.com/insolar/insolar/application/contract/rootdomain"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	cryptoHelper "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/ledger/ledgertestutils"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
	"github.com/insolar/insolar/logicrunner/goplugin/preprocessor"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/testmessagebus"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/ugorji/go/codec"
)

var icc = ""
var runnerbin = ""
var parallel = true

func TestMain(m *testing.M) {
	var err error
	err = log.SetLevel("Debug")
	if err != nil {
		log.Errorln(err.Error())
	}
	if runnerbin, icc, err = goplugintestutils.Build(); err != nil {
		fmt.Println("Logic runner build failed, skip tests:", err.Error())
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func MessageBusTrivialBehavior(mb *testmessagebus.TestMessageBus, lr core.LogicRunner) {
	mb.ReRegister(core.TypeCallMethod, lr.Execute)
	mb.ReRegister(core.TypeCallConstructor, lr.Execute)
	mb.ReRegister(core.TypeValidateCaseBind, lr.ValidateCaseBind)
	mb.ReRegister(core.TypeValidationResults, lr.ProcessValidationResults)
	mb.ReRegister(core.TypeExecutorResults, lr.ExecutorResults)

}

func PrepareLrAmCb(t testing.TB) (core.LogicRunner, core.ArtifactManager, *goplugintestutils.ContractsBuilder, func()) {
	lrSock := os.TempDir() + "/" + testutils.RandomString() + ".sock"
	rundSock := os.TempDir() + "/" + testutils.RandomString() + ".sock"

	rundCleaner, err := testutils.StartInsgorund(runnerbin, "unix", rundSock, "unix", lrSock)
	assert.NoError(t, err)

	lr, err := NewLogicRunner(&configuration.LogicRunner{
		RPCListen:   lrSock,
		RPCProtocol: "unix",
		GoPlugin: &configuration.GoPlugin{
			RunnerListen:   rundSock,
			RunnerProtocol: "unix",
		},
	})
	assert.NoError(t, err, "Initialize runner")

	l, cleaner := ledgertestutils.TmpLedger(t, lr, "")
	mb := testmessagebus.NewTestMessageBus()
	assert.NoError(t, lr.Start(core.Components{
		Ledger:     l,
		MessageBus: mb,
	}), "starting logicrunner")
	MessageBusTrivialBehavior(mb, lr)
	err = l.GetPulseManager().Set(*pulsar.NewPulse(configuration.NewPulsar().NumberDelta, 0, &entropygenerator.StandardEntropyGenerator{}))
	if err != nil {
		t.Fatal("pulse set died, ", err)
	}
	am := l.GetArtifactManager()
	cb := goplugintestutils.NewContractBuilder(am, icc)

	return lr, am, cb, func() {
		cb.Clean()
		lr.Stop()
		cleaner()
		rundCleaner()
	}
}

func ValidateAllResults(t testing.TB, lr core.LogicRunner, mustfail ...core.RecordRef) {
	failmap := make(map[core.RecordRef]struct{})
	for _, r := range mustfail {
		failmap[r] = struct{}{}
	}
	rlr := lr.(*LogicRunner)
	rlr.caseBindMutex.Lock()
	rlrcbr := rlr.caseBind.Records
	rlr.caseBind.Records = make(map[core.RecordRef][]core.CaseRecord)
	rlr.caseBindMutex.Unlock()
	for ref, cr := range rlrcbr {
		assert.Equal(t, configuration.NewPulsar().NumberDelta, uint32(rlr.caseBind.Pulse.PulseNumber), "right pulsenumber")
		vstep, err := lr.Validate(ref, rlr.caseBind.Pulse, cr)
		if _, ok := failmap[ref]; ok {
			assert.Error(t, err, "validation")
			assert.True(t, len(cr) > vstep, "Validation failed before end")
		} else {
			assert.NoError(t, err, "validation")
			assert.Equal(t, len(cr), vstep, "Validation passed to the end")
		}
	}
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
	if parallel {
		t.Parallel()
	}
	lr, err := NewLogicRunner(&configuration.LogicRunner{})
	assert.NoError(t, err)
	lr.OnPulse(*pulsar.NewPulse(configuration.NewPulsar().NumberDelta, 0, &entropygenerator.StandardEntropyGenerator{}))

	mb := testmessagebus.NewTestMessageBus()
	comps := core.Components{
		Ledger:     &testLedger{am: goplugintestutils.NewTestArtifactManager()},
		MessageBus: mb,
	}
	assert.NoError(t, lr.Start(comps))
	assert.IsType(t, &LogicRunner{}, lr)
	MessageBusTrivialBehavior(mb, lr)

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

func (r *testLedger) HandleMessage(core.Message) (core.Reply, error) {
	panic("implement me")
}

func TestExecution(t *testing.T) {
	if parallel {
		t.Parallel()
	}
	am := goplugintestutils.NewTestArtifactManager()
	ld := &testLedger{am: am}
	lr, err := NewLogicRunner(&configuration.LogicRunner{})
	assert.NoError(t, err)
	mb := testmessagebus.NewTestMessageBus()
	MessageBusTrivialBehavior(mb, lr)
	lr.Start(core.Components{
		Ledger:     ld,
		MessageBus: mb,
	})
	lr.OnPulse(*pulsar.NewPulse(configuration.NewPulsar().NumberDelta, 0, &entropygenerator.StandardEntropyGenerator{}))

	codeRef := core.NewRefFromBase58("someCode")
	dataRef := core.NewRefFromBase58("someObject")
	classRef := core.NewRefFromBase58("someClass")
	am.Objects[dataRef] = &goplugintestutils.TestObjectDescriptor{
		AM:    am,
		Data:  []byte("origData"),
		Code:  &codeRef,
		Class: &classRef,
	}
	am.Classes[classRef] = &goplugintestutils.TestClassDescriptor{AM: am, ARef: &classRef, ACode: &codeRef}
	am.Codes[codeRef] = &goplugintestutils.TestCodeDescriptor{ARef: codeRef, AMachineType: core.MachineTypeGoPlugin}

	te := newTestExecutor()
	te.methodResponses = append(te.methodResponses, &testResp{data: []byte("data"), res: core.Arguments("res")})

	err = lr.RegisterExecutor(core.MachineTypeGoPlugin, te)
	assert.NoError(t, err)

	resp, err := lr.Execute(&message.CallMethod{ObjectRef: dataRef})
	assert.NoError(t, err)
	assert.Equal(t, []byte("data"), resp.(*reply.CallMethod).Data)
	assert.Equal(t, []byte("res"), resp.(*reply.CallMethod).Result)

	te.constructorResponses = append(te.constructorResponses, &testResp{data: []byte("data"), res: core.Arguments("res")})
	resp, err = lr.Execute(&message.CallConstructor{ClassRef: classRef})
	assert.NoError(t, err)
}

func TestContractCallingContract(t *testing.T) {
	if parallel {
		t.Parallel()
	}
	var contractOneCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import "github.com/insolar/insolar/application/proxy/two"
import "github.com/insolar/insolar/core"

type One struct {
	foundation.BaseContract
	Friend core.RecordRef
}

func (r *One) Hello(s string) string {
	holder := two.New()
	friend := holder.AsChild(r.GetReference())
	res := friend.Hello(s)
	r.Friend = friend.GetReference()
	return "Hi, " + s + "! Two said: " + res
}

func (r *One) Again(s string) string {
	res := two.GetObject(r.Friend).Hello(s)
	return "Hi, " + s + "! Two said: " + res
}

func (r *One)GetFriend() core.RecordRef {
	return r.Friend
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
	return &Two{X:0};
}

func (r *Two) Hello(s string) string {
	r.X ++
	return fmt.Sprintf("Hello you too, %s. %d times!", s, r.X)
}
`

	lr, am, cb, cleaner := PrepareLrAmCb(t)
	defer cleaner()

	err := cb.Build(map[string]string{"one": contractOneCode, "two": contractTwoCode})
	assert.NoError(t, err)

	obj, err := am.RegisterRequest(&message.CallConstructor{})
	_, err = am.ActivateObject(
		core.RecordRef{}, *obj,
		*cb.Classes["one"],
		*am.GenesisRef(),
		goplugintestutils.CBORMarshal(t, &struct{}{}),
	)
	assert.NoError(t, err)

	resp, err := lr.Execute(&message.CallMethod{
		ObjectRef: *obj,
		Method:    "Hello",
		Arguments: goplugintestutils.CBORMarshal(t, []interface{}{"ins"}),
	})
	assert.NoError(t, err, "contract call")
	r := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	f := r.([]interface{})[0]
	assert.Equal(t, "Hi, ins! Two said: Hello you too, ins. 1 times!", f)

	for i := 2; i <= 5; i++ {
		resp, err := lr.Execute(&message.CallMethod{
			ObjectRef: *obj,
			Method:    "Again",
			Arguments: goplugintestutils.CBORMarshal(t, []interface{}{"ins"}),
		})
		assert.NoError(t, err, "contract call")
		r := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
		f := r.([]interface{})[0]
		assert.Equal(t, fmt.Sprintf("Hi, ins! Two said: Hello you too, ins. %d times!", i), f)
	}

	resp, err = lr.Execute(&message.CallMethod{
		ObjectRef: *obj,
		Method:    "GetFriend",
		Arguments: goplugintestutils.CBORMarshal(t, []interface{}{}),
	})
	assert.NoError(t, err, "contract call")
	r = goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	r0 := r.([]interface{})[0].([]uint8)
	var two core.RecordRef
	for i := 0; i < 64; i++ {
		two[i] = r0[i]
	}

	for i := 6; i <= 9; i++ {
		resp, err := lr.Execute(&message.CallMethod{
			ObjectRef: two,
			Method:    "Hello",
			Arguments: goplugintestutils.CBORMarshal(t, []interface{}{"Insolar"}),
		})
		assert.NoError(t, err, "contract call")
		r := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
		f := r.([]interface{})[0]
		assert.Equal(t, fmt.Sprintf("Hello you too, Insolar. %d times!", i), f)
	}
	ValidateAllResults(t, lr)

}

func TestInjectingDelegate(t *testing.T) {
	if parallel {
		t.Parallel()
	}
	var contractOneCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import "github.com/insolar/insolar/application/proxy/two"

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
	lr, am, cb, cleaner := PrepareLrAmCb(t)
	gp := lr.(*LogicRunner).Executors[core.MachineTypeGoPlugin]
	defer cleaner()

	data := goplugintestutils.CBORMarshal(t, &struct{}{})
	argsSerialized := goplugintestutils.CBORMarshal(t, []interface{}{"ins"})

	err := cb.Build(map[string]string{"one": contractOneCode, "two": contractTwoCode})
	assert.NoError(t, err)

	obj, err := am.RegisterRequest(&message.CallConstructor{})
	_, err = am.ActivateObject(
		core.RecordRef{}, *obj,
		*cb.Classes["one"],
		*am.GenesisRef(),
		data,
	)
	assert.NoError(t, err)

	_, res, err := gp.CallMethod(
		&core.LogicCallContext{Class: cb.Classes["one"], Callee: obj}, *cb.Codes["one"],
		data, "Hello", argsSerialized,
	)
	assert.NoError(t, err)

	resParsed := goplugintestutils.CBORUnMarshalToSlice(t, res)
	assert.Equal(t, "Hi, ins! Two said: Hello you too, ins. 644 times!", resParsed[0])

	_, res, err = gp.CallMethod(
		&core.LogicCallContext{Class: cb.Classes["one"], Callee: obj}, *cb.Codes["one"],
		data, "HelloFromDelegate", argsSerialized,
	)
	assert.NoError(t, err)

	resParsed = goplugintestutils.CBORUnMarshalToSlice(t, res)

	assert.Equal(t, "Hello you too, ins. 1288 times!", resParsed[0])
}

func TestBasicNotificationCall(t *testing.T) {
	if parallel {
		t.Parallel()
	}
	var contractOneCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import "github.com/insolar/insolar/application/proxy/two"

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
	// TODO: use am := testutil.NewTestArtifactManager() here
	lr, am, cb, cleaner := PrepareLrAmCb(t)
	gp := lr.(*LogicRunner).Executors[core.MachineTypeGoPlugin]
	defer cleaner()

	data := goplugintestutils.CBORMarshal(t, &struct{}{})
	argsSerialized := goplugintestutils.CBORMarshal(t, []interface{}{})
	err := cb.Build(map[string]string{"one": contractOneCode, "two": contractTwoCode})
	assert.NoError(t, err)

	obj, err := am.RegisterRequest(&message.CallConstructor{})
	_, err = am.ActivateObject(
		core.RecordRef{},
		*obj,
		*cb.Classes["one"],
		*am.GenesisRef(),
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
	if parallel {
		t.Parallel()
	}
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
	lr, _, cb, cleaner := PrepareLrAmCb(t)
	gp := lr.(*LogicRunner).Executors[core.MachineTypeGoPlugin]
	defer cleaner()

	data := goplugintestutils.CBORMarshal(t, &struct{}{})
	argsSerialized := goplugintestutils.CBORMarshal(t, []struct{}{})

	err := cb.Build(map[string]string{"one": code})
	assert.NoError(t, err)

	_, res, err := gp.CallMethod(
		&core.LogicCallContext{Class: cb.Classes["one"]}, *cb.Codes["one"],
		data, "Hello", argsSerialized,
	)
	assert.NoError(t, err)

	resParsed := goplugintestutils.CBORUnMarshalToSlice(t, res)
	assert.Equal(t, cb.Classes["one"].String(), resParsed[0])
}

func TestDeactivation(t *testing.T) {
	if parallel {
		t.Parallel()
	}
	var code = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type One struct {
	foundation.BaseContract
}

func (r *One) Kill() {
	r.SelfDestruct()
}
`
	lr, am, cb, cleaner := PrepareLrAmCb(t)
	defer cleaner()

	err := cb.Build(map[string]string{"one": code})
	assert.NoError(t, err)

	obj, err := am.RegisterRequest(&message.CallConstructor{})
	_, err = am.ActivateObject(
		core.RecordRef{}, *obj,
		*cb.Classes["one"],
		*am.GenesisRef(),
		goplugintestutils.CBORMarshal(t, &struct{}{}),
	)
	assert.NoError(t, err)

	_, err = lr.Execute(&message.CallMethod{
		ObjectRef: *obj,
		Method:    "Kill",
		Arguments: goplugintestutils.CBORMarshal(t, []interface{}{}),
	})
	assert.NoError(t, err, "contract call")
}

func TestPanic(t *testing.T) {
	if parallel {
		t.Parallel()
	}
	var code = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type One struct {
	foundation.BaseContract
}

func (r *One) Panic() {
	panic("haha")
}
func (r *One) NotPanic() {
}
`
	lr, am, cb, cleaner := PrepareLrAmCb(t)
	defer cleaner()

	err := cb.Build(map[string]string{"one": code})
	assert.NoError(t, err)

	obj, err := am.RegisterRequest(&message.CallConstructor{})
	_, err = am.ActivateObject(
		core.RecordRef{}, *obj,
		*cb.Classes["one"],
		*am.GenesisRef(),
		goplugintestutils.CBORMarshal(t, &struct{}{}),
	)
	assert.NoError(t, err)

	_, err = lr.Execute(&message.CallMethod{
		ObjectRef: *obj,
		Method:    "Panic",
		Arguments: goplugintestutils.CBORMarshal(t, []interface{}{}),
	})
	assert.Error(t, err)

	_, err = lr.Execute(&message.CallMethod{
		ObjectRef: *obj,
		Method:    "NotPanic",
		Arguments: goplugintestutils.CBORMarshal(t, []interface{}{}),
	})
	assert.NoError(t, err)
}

func TestGetChildren(t *testing.T) {
	if parallel {
		t.Parallel()
	}
	goContract := `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/application/proxy/child"
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
	lr, am, cb, cleaner := PrepareLrAmCb(t)
	defer cleaner()

	err := cb.Build(map[string]string{"child": goChild})
	assert.NoError(t, err)
	err = cb.Build(map[string]string{"contract": goContract})
	assert.NoError(t, err)

	domain := core.NewRefFromBase58("c1")
	contract, err := am.RegisterRequest(&message.CallConstructor{ClassRef: core.NewRefFromBase58("dassads")})
	_, err = am.ActivateObject(domain, *contract, *cb.Classes["contract"], *am.GenesisRef(), goplugintestutils.CBORMarshal(t, nil))
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, contract, nil, "contract created")

	resp, err := lr.Execute(&message.CallMethod{
		ObjectRef: *contract,
		Method:    "NewChilds",
		Arguments: goplugintestutils.CBORMarshal(t, []interface{}{10}),
	})
	assert.NoError(t, err, "contract call")
	r := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	assert.Equal(t, []interface{}([]interface{}{uint64(45)}), r)

	resp, err = lr.Execute(&message.CallMethod{
		ObjectRef: *contract,
		Method:    "SumChilds",
		Arguments: goplugintestutils.CBORMarshal(t, []interface{}{}),
	})

	ValidateAllResults(t, lr)

	assert.NoError(t, err, "contract call")
	r = goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	assert.Equal(t, []interface{}([]interface{}{uint64(45)}), r)
}

func TestFailValidate(t *testing.T) {
	if parallel {
		t.Parallel()
	}
	goContract := `
package main

import (
	"math/rand"
	"time"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Contract struct {
	foundation.BaseContract
}

func (c *Contract) Rand() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(77)
}
`
	lr, am, cb, cleaner := PrepareLrAmCb(t)
	defer cleaner()

	err := cb.Build(map[string]string{"contract": goContract})
	assert.NoError(t, err)

	domain := core.NewRefFromBase58("c1")
	contract, err := am.RegisterRequest(&message.CallConstructor{ClassRef: core.NewRefFromBase58("dassads")})
	_, err = am.ActivateObject(domain, *contract, *cb.Classes["contract"], *am.GenesisRef(), goplugintestutils.CBORMarshal(t, nil))
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, contract, nil, "contract created")

	for i := 0; i < 5; i++ {
		_, err = lr.Execute(&message.CallMethod{
			ObjectRef: *contract,
			Method:    "Rand",
			Arguments: goplugintestutils.CBORMarshal(t, []interface{}{}),
		})
		assert.NoError(t, err, "contract call")
	}
	ValidateAllResults(t, lr, *contract)
}

func TestErrorInterface(t *testing.T) {
	if parallel {
		t.Parallel()
	}
	var contractOneCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/application/proxy/two"
)

type One struct {
	foundation.BaseContract
}

func (r *One) AnError() error {
	holder := two.New()
	friend := holder.AsChild(r.GetReference())

	return friend.AnError()
}

func (r *One) NoError() error {
	holder := two.New()
	friend := holder.AsChild(r.GetReference())

	return friend.NoError()
}
`

	var contractTwoCode = `
package main

import (
	"errors"

	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Two struct {
	foundation.BaseContract
}
func New() *Two {
	return &Two{}
}
func (r *Two) AnError() error {
	return errors.New("an error")
}
func (r *Two) NoError() error {
	return nil
}
`
	lr, am, cb, cleaner := PrepareLrAmCb(t)
	defer cleaner()

	err := cb.Build(map[string]string{
		"one": contractOneCode,
		"two": contractTwoCode,
	})
	assert.NoError(t, err)

	domain := core.NewRefFromBase58("c1")
	contract, err := am.RegisterRequest(&message.CallConstructor{})
	_, err = am.ActivateObject(domain, *contract, *cb.Classes["one"], *am.GenesisRef(), goplugintestutils.CBORMarshal(t, nil))
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, contract, nil, "contract created")

	resp, err := lr.Execute(&message.CallMethod{
		ObjectRef: *contract,
		Method:    "AnError",
		Arguments: goplugintestutils.CBORMarshal(t, []interface{}{}),
	})
	assert.NoError(t, err, "contract call")

	ValidateAllResults(t, lr)

	ch := new(codec.CborHandle)
	res := []interface{}{&foundation.Error{}}
	err = codec.NewDecoderBytes(resp.(*reply.CallMethod).Result, ch).Decode(&res)
	assert.NoError(t, err, "contract call")
	assert.Equal(t, &foundation.Error{S: "an error"}, res[0])

	resp, err = lr.Execute(&message.CallMethod{
		ObjectRef: *contract,
		Method:    "NoError",
		Arguments: goplugintestutils.CBORMarshal(t, []interface{}{}),
	})
	assert.NoError(t, err, "contract call")

	ValidateAllResults(t, lr)

	r := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	assert.Equal(t, []interface{}([]interface{}{nil}), r)
}

func TestNilResult(t *testing.T) {
	if parallel {
		t.Parallel()
	}
	var contractOneCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/application/proxy/two"
)

type One struct {
	foundation.BaseContract
}

func (r *One) Hello() *string {
	holder := two.New()
	friend := holder.AsChild(r.GetReference())

	return friend.Hello()
}
`

	var contractTwoCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Two struct {
	foundation.BaseContract
}
func New() *Two {
	return &Two{}
}
func (r *Two) Hello() *string {
	return nil
}
`
	lr, am, cb, cleaner := PrepareLrAmCb(t)
	defer cleaner()

	err := cb.Build(map[string]string{
		"one": contractOneCode,
		"two": contractTwoCode,
	})
	assert.NoError(t, err)

	domain := core.NewRefFromBase58("c1")
	contract, err := am.RegisterRequest(&message.CallConstructor{})
	_, err = am.ActivateObject(domain, *contract, *cb.Classes["one"], *am.GenesisRef(), goplugintestutils.CBORMarshal(t, nil))
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, contract, nil, "contract created")

	resp, err := lr.Execute(&message.CallMethod{
		ObjectRef: *contract,
		Method:    "Hello",
		Arguments: goplugintestutils.CBORMarshal(t, []interface{}{}),
	})
	assert.NoError(t, err, "contract call")

	ValidateAllResults(t, lr)

	r := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	assert.Equal(t, []interface{}([]interface{}{nil}), r)
}

type Caller struct {
	member string
	key    *ecdsa.PrivateKey
	lr     core.LogicRunner
	t      *testing.T
}

func (s *Caller) SignedCall(ref core.RecordRef, delegate core.RecordRef, method string, params []interface{}, resp []interface{}) {
	seed := make([]byte, 32)
	_, err := rand.Read(seed)
	assert.NoError(s.t, err)

	buf := goplugintestutils.CBORMarshal(s.t, params)

	serialized, err := signer.Serialize(ref[:], delegate[:], method, buf, seed)
	assert.NoError(s.t, err)

	sign, err := cryptoHelper.Sign(serialized, s.key)
	assert.NoError(s.t, err)
	res, err := s.lr.Execute(&message.CallMethod{
		ObjectRef: core.NewRefFromBase58(s.member),
		Method:    "AuthorizedCall",
		Arguments: goplugintestutils.CBORMarshal(s.t, []interface{}{ref, delegate, method, buf, seed, sign}),
	})
	assert.NoError(s.t, err, "contract call")
	result := goplugintestutils.CBORUnMarshal(s.t, res.(*reply.CallMethod).Result).([]interface{})
	assert.Nil(s.t, result[1])
	if result[0] != nil {
		ch := new(codec.CborHandle)
		err = codec.NewDecoderBytes(result[0].([]byte), ch).Decode(&resp)
	}
}

func TestRootDomainContract(t *testing.T) {
	if parallel {
		t.Parallel()
	}
	rootDomainCode, err := ioutil.ReadFile("../application/contract/rootdomain/rootdomain.go" +
		"")
	if err != nil {
		fmt.Print(err)
	}
	memberCode, err := ioutil.ReadFile("../application/contract/member/member.go")
	if err != nil {
		fmt.Print(err)
	}
	allowanceCode, err := ioutil.ReadFile("../application/contract/allowance/allowance.go")
	if err != nil {
		fmt.Print(err)
	}
	walletCode, err := ioutil.ReadFile("../application/contract/wallet/wallet.go")
	if err != nil {
		fmt.Print(err)
	}

	lr, am, cb, cleaner := PrepareLrAmCb(t)
	defer cleaner()
	err = cb.Build(map[string]string{"member": string(memberCode), "allowance": string(allowanceCode), "wallet": string(walletCode), "rootdomain": string(rootDomainCode)})
	assert.NoError(t, err)

	// Initializing Root Domain
	rootDomainRef, err := am.RegisterRequest(&message.BootstrapRequest{Name: "c1"})
	_, err = am.ActivateObject(core.RecordRef{}, *rootDomainRef, *cb.Classes["rootdomain"], *am.GenesisRef(), goplugintestutils.CBORMarshal(t, nil))
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, rootDomainRef, nil, "contract created")

	// Creating Root member
	rootKey, err := cryptoHelper.GeneratePrivateKey()
	assert.NoError(t, err)
	rootPubKey, err := cryptoHelper.ExportPublicKey(&rootKey.PublicKey)
	assert.NoError(t, err)

	rootMemberRef, err := am.RegisterRequest(&message.BootstrapRequest{Name: "c2"})
	assert.NoError(t, err)
	_, err = am.ActivateObject(core.RecordRef{}, *rootMemberRef, *cb.Classes["member"], *rootDomainRef, goplugintestutils.CBORMarshal(t, member.New("root", rootPubKey)))
	assert.NoError(t, err)

	// Updating root domain with root member
	_, err = am.UpdateObject(core.RecordRef{}, core.RecordRef{}, *rootDomainRef, goplugintestutils.CBORMarshal(t, rootdomain.RootDomain{RootMember: *rootMemberRef}))
	assert.NoError(t, err)

	root := Caller{rootMemberRef.String(), rootKey, lr, t}

	// Creating Member1
	member1Key, err := cryptoHelper.GeneratePrivateKey()
	assert.NoError(t, err)
	member1PubKey, err := cryptoHelper.ExportPublicKey(&member1Key.PublicKey)
	assert.NoError(t, err)

	res1 := []interface{}{""}
	root.SignedCall(*rootDomainRef, core.RecordRef{}, "CreateMember", []interface{}{"Member1", member1PubKey}, res1)
	member1Ref := res1[0].(string)
	assert.NotEqual(t, "", member1Ref)

	// Creating Member2
	member2Key, err := cryptoHelper.GeneratePrivateKey()
	assert.NoError(t, err)
	member2PubKey, err := cryptoHelper.ExportPublicKey(&member2Key.PublicKey)
	assert.NoError(t, err)

	res2 := []interface{}{""}
	root.SignedCall(*rootDomainRef, core.RecordRef{}, "CreateMember", []interface{}{"Member2", member2PubKey}, res2)
	member2Ref := res2[0].(string)
	assert.NotEqual(t, "", member2Ref)

	// Transfer 1 coin from Member1 to Member2
	member1 := Caller{member1Ref, member1Key, lr, t}
	z := core.NewRefFromBase58(member2Ref)
	member1.SignedCall(core.NewRefFromBase58(member1Ref), *cb.Classes["wallet"], "Transfer", []interface{}{1, &z}, nil)

	// Verify Member1 balance
	res3 := []interface{}{0}
	root.SignedCall(core.NewRefFromBase58(member1Ref), *cb.Classes["wallet"], "GetTotalBalance", []interface{}{}, res3)
	assert.Equal(t, 999, res3[0])

	// Verify Member2 balance
	res4 := []interface{}{0}
	root.SignedCall(core.NewRefFromBase58(member2Ref), *cb.Classes["wallet"], "GetTotalBalance", []interface{}{}, res4)
	assert.Equal(t, 1001, res4[0])
}

func BenchmarkContractCall(b *testing.B) {
	goParent := `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import "github.com/insolar/insolar/application/proxy/child"
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
	lr, am, cb, cleaner := PrepareLrAmCb(b)
	defer cleaner()
	err := cb.Build(map[string]string{"child": goChild, "parent": goParent})
	assert.NoError(b, err)

	domain := core.NewRefFromBase58("c1")
	parent, err := am.RegisterRequest(&message.CallConstructor{})
	_, err = am.ActivateObject(domain, *parent, *cb.Classes["parent"], *am.GenesisRef(), goplugintestutils.CBORMarshal(b, nil))
	assert.NoError(b, err, "create parent")
	assert.NotEqual(b, parent, nil, "parent created")
	child, err := am.RegisterRequest(&message.CallConstructor{ParentRef: *parent})
	_, err = am.ActivateObject(domain, *child, *cb.Classes["child"], *am.GenesisRef(), goplugintestutils.CBORMarshal(b, nil))
	assert.NoError(b, err, "create child")
	assert.NotEqual(b, child, nil, "child created")

	b.N = 1000
	for i := 0; i < b.N; i++ {
		resp, err := lr.Execute(&message.CallMethod{
			ObjectRef: *parent,
			Method:    "CCC",
			Arguments: goplugintestutils.CBORMarshal(b, []interface{}{child}),
		})
		assert.NoError(b, err, "parent call")
		r := goplugintestutils.CBORUnMarshal(b, resp.(*reply.CallMethod).Result)
		assert.Equal(b, []interface{}([]interface{}{uint64(5)}), r)
	}
}

func TestProxyGeneration(t *testing.T) {
	if parallel {
		t.Parallel()
	}
	contracts, err := preprocessor.GetRealContractsNames()
	assert.NoError(t, err)

	for _, contract := range contracts {
		t.Run(contract, func(t *testing.T) {
			parsed, err := preprocessor.ParseFile("../application/contract/" + contract + "/" + contract + ".go")
			assert.NoError(t, err)

			proxyPath, err := preprocessor.GetRealApplicationDir("proxy")
			assert.NoError(t, err)

			name, err := parsed.ProxyPackageName()
			assert.NoError(t, err)

			proxy := path.Join(proxyPath, name, name+".go")
			_, err = os.Stat(proxy)
			assert.NoError(t, err)

			buff := bytes.NewBufferString("")
			parsed.WriteProxy("", buff)

			cmd := exec.Command("diff", proxy, "-")
			cmd.Stdin = buff
			out, err := cmd.CombinedOutput()
			assert.NoError(t, err, string(out))
		})
	}
}
