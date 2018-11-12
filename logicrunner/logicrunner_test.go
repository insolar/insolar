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
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils/certificate"
	"github.com/insolar/insolar/testutils/network"
	"github.com/insolar/insolar/testutils/nodekeeper"

	"github.com/insolar/insolar/application/contract/member"
	"github.com/insolar/insolar/application/contract/member/signer"
	"github.com/insolar/insolar/application/contract/rootdomain"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/core/utils"
	cryptoHelper "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/ledger/ledgertestutils"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/testmessagebus"
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

func PrepareLrAmCbPm(t testing.TB) (core.LogicRunner, core.ArtifactManager, *goplugintestutils.ContractsBuilder, core.PulseManager, func()) {
	ctx := context.TODO()
	lrSock := os.TempDir() + "/" + testutils.RandomString() + ".sock"
	rundSock := os.TempDir() + "/" + testutils.RandomString() + ".sock"

	rundCleaner, err := goplugintestutils.StartInsgorund(runnerbin, "unix", rundSock, "unix", lrSock)
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

	ce := certificate.GetTestCertificate()
	nk := nodekeeper.GetTestNodekeeper(ce)

	var pulseNumber core.PulseNumber
	pulseNumber = 0

	mb := testmessagebus.NewTestMessageBus()
	mb.PulseNumber = pulseNumber

	nw := network.GetTestNetwork()
	// FIXME: TmpLedger is deprecated. Use mocks instead.
	l, cleaner := ledgertestutils.TmpLedger(
		t, "",
		core.Components{
			LogicRunner: lr,
			NodeNetwork: nk,
			MessageBus:  mb,
			Network:     nw,
		},
	)

	cm := &component.Manager{}
	cm.Register(nk, l, lr, nw, mb)
	err = cm.Start(ctx)
	assert.NoError(t, err)

	MessageBusTrivialBehavior(mb, lr)
	pm := l.GetPulseManager()

	newTestPulse(ctx, lr, mb)

	assert.NoError(t, err)
	if err != nil {
		t.Fatal("pulse set died, ", err)
	}
	am := l.GetArtifactManager()
	cb := goplugintestutils.NewContractBuilder(am, icc)

	return lr, am, cb, pm, func() {
		cb.Clean()
		lr.Stop(ctx)
		cleaner()
		rundCleaner()
	}
}

func newTestPulse(ctx context.Context, lr *LogicRunner, mb *testmessagebus.TestMessageBus) {
	currentPulse, _ := lr.Ledger.GetPulseManager().Current(ctx)
	newPulseNumber := currentPulse.PulseNumber + 1
	lr.Ledger.GetPulseManager().Set(
		ctx,
		core.Pulse{PulseNumber: newPulseNumber, Entropy: core.Entropy{}},
	)

	mb.PulseNumber = newPulseNumber
}

func ValidateAllResults(t testing.TB, ctx context.Context, lr core.LogicRunner, mustfail ...core.RecordRef) {
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
		log.Debugf("TEST validating: %s", ref)
		vstep, err := lr.Validate(ref, *rlr.pulse(ctx), cr)
		if _, ok := failmap[ref]; ok {
			assert.Error(t, err, "validation %s", ref)
			assert.True(t, len(cr) > vstep, "Validation failed before end %s", ref)
		} else {
			assert.NoError(t, err, "validation %s", ref)
			assert.Equal(t, len(cr), vstep, "Validation passed to the end %s", ref)
		}
	}
}

func executeMethod(ctx context.Context, lr core.LogicRunner, pm core.PulseManager, objRef core.RecordRef, nonce uint64, method string, arguments core.Arguments) (core.Reply, error) {
	msg := &message.CallMethod{
		ObjectRef: objRef,
		Method:    method,
		Arguments: arguments,
	}

	if nonce != 0 {
		msg.Nonce = nonce
	}

	pulse, _ := pm.Current(ctx)

	key, _ := cryptoHelper.GeneratePrivateKey()

	parcel, _ := message.NewParcel(ctx, msg, testutils.RandomRef(), key, pulse.PulseNumber, nil)
	ctx = inslogger.ContextWithTrace(ctx, utils.RandTraceID())
	resp, err := lr.Execute(
		ctx,
		parcel,
	)

	return resp, err
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

func getRefFromID(id *core.RecordID) *core.RecordRef {
	ref := core.RecordRef{}
	ref.SetRecord(*id)
	return &ref
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

func (r *One) Hello(s string) (string, error) {
	holder := two.New()
	friend, err := holder.AsChild(r.GetReference())
	if err != nil {
		return "1", err
	}

	res, err := friend.Hello(s)
	if err != nil {
		return "2", err
	}
	
	r.Friend = friend.GetReference()
	return "Hi, " + s + "! Two said: " + res, nil
}

func (r *One) Again(s string) (string, error) {
	res, err := two.GetObject(r.Friend).Hello(s)
	if err != nil {
		return "", err
	}
	
	return "Hi, " + s + "! Two said: " + res, nil
}

func (r *One)GetFriend() (core.RecordRef, error) {
	return r.Friend, nil
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

func New() (*Two, error) {
	return &Two{X:0}, nil;
}

func (r *Two) Hello(s string) (string, error) {
	r.X ++
	return fmt.Sprintf("Hello you too, %s. %d times!", s, r.X), nil
}
`
	ctx := context.Background()

	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()

	err := cb.Build(map[string]string{"one": contractOneCode, "two": contractTwoCode})
	assert.NoError(t, err)

	objID, err := am.RegisterRequest(ctx, &message.CallConstructor{})
	assert.NoError(t, err)
	obj := getRefFromID(objID)
	_, err = am.ActivateObject(
		ctx,
		core.RecordRef{}, *obj,
		*am.GenesisRef(),
		*cb.Prototypes["one"],
		false,
		goplugintestutils.CBORMarshal(t, &struct{}{}),
	)
	assert.NoError(t, err)

	resp, err := executeMethod(ctx, lr, pm, *obj, 0, "Hello", goplugintestutils.CBORMarshal(t, []interface{}{"ins"}))
	assert.NoError(t, err, "contract call")
	r := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	f := r.([]interface{})[0]
	assert.Equal(t, "Hi, ins! Two said: Hello you too, ins. 1 times!", f)

	for i := 2; i <= 5; i++ {
		resp, err = executeMethod(ctx, lr, pm, *obj, uint64(i), "Again", goplugintestutils.CBORMarshal(t, []interface{}{"ins"}))
		assert.NoError(t, err, "contract call")
		r := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
		f := r.([]interface{})[0]
		assert.Equal(t, fmt.Sprintf("Hi, ins! Two said: Hello you too, ins. %d times!", i), f)
	}

	resp, err = executeMethod(ctx, lr, pm, *obj, 0, "GetFriend", goplugintestutils.CBORMarshal(t, []interface{}{}))
	assert.NoError(t, err, "contract call")
	r = goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	r0 := r.([]interface{})[0].([]uint8)
	var two core.RecordRef
	for i := 0; i < 64; i++ {
		two[i] = r0[i]
	}

	for i := 6; i <= 9; i++ {
		resp, err = executeMethod(ctx, lr, pm, two, uint64(i), "Hello", goplugintestutils.CBORMarshal(t, []interface{}{"Insolar"}))
		assert.NoError(t, err, "contract call")
		r := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
		f := r.([]interface{})[0]
		assert.Equal(t, fmt.Sprintf("Hello you too, Insolar. %d times!", i), f)
	}
	ValidateAllResults(t, ctx, lr)

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

func (r *One) Hello(s string) (string, error) {
	holder := two.New()
	friend, err := holder.AsDelegate(r.GetReference())
	if err != nil {
		return "", err
	}

	res, err := friend.Hello(s)
	if err != nil {
		return "", err
	}

	return "Hi, " + s + "! Two said: " + res, nil
}

func (r *One) HelloFromDelegate(s string) (string, error) {
	friend, err := two.GetImplementationFrom(r.GetReference())
	if err != nil {
		return "", err
	}

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

func New() (*Two, error) {
	return &Two{X:322}, nil
}

func (r *Two) Hello(s string) (string, error) {
	r.X *= 2
	return fmt.Sprintf("Hello you too, %s. %d times!", s, r.X), nil
}
`
	ctx := context.TODO()
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()

	data := goplugintestutils.CBORMarshal(t, &struct{}{})

	err := cb.Build(map[string]string{"one": contractOneCode, "two": contractTwoCode})
	assert.NoError(t, err)

	objID, err := am.RegisterRequest(ctx, &message.CallConstructor{})
	assert.NoError(t, err)
	obj := getRefFromID(objID)
	_, err = am.ActivateObject(
		ctx,
		core.RecordRef{}, *obj,
		*am.GenesisRef(),
		*cb.Prototypes["one"],
		false,
		data,
	)
	assert.NoError(t, err)

	resp, err := executeMethod(ctx, lr, pm, *obj, 0, "Hello", goplugintestutils.CBORMarshal(t, []interface{}{"ins"}))
	assert.NoError(t, err)

	r := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	assert.Equal(t, []interface{}{"Hi, ins! Two said: Hello you too, ins. 644 times!", nil}, r)

	resp, err = executeMethod(ctx, lr, pm, *obj, 0, "HelloFromDelegate", goplugintestutils.CBORMarshal(t, []interface{}{"ins"}))
	assert.NoError(t, err)
	r = goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	assert.Equal(t, []interface{}{"Hello you too, ins. 1288 times!", nil}, r)

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

func (r *One) Hello() error {
	holder := two.New()

	friend, err := holder.AsDelegate(r.GetReference())
	if err != nil {
		return err
	}

	err = friend.HelloNoWait()
	if err != nil {
		return err
	}

	return nil
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

func New() (*Two, error) {
	return &Two{X:322}, nil
}

func (r *Two) Hello() (string, error) {
	r.X *= 2
	return fmt.Sprintf("Hello %d times!", r.X), nil
}
`
	ctx := context.TODO()
	// TODO: use am := testutil.NewTestArtifactManager() here
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()

	err := cb.Build(map[string]string{"one": contractOneCode, "two": contractTwoCode})
	assert.NoError(t, err)

	objID, err := am.RegisterRequest(ctx, &message.CallConstructor{})
	assert.NoError(t, err)
	obj := getRefFromID(objID)
	_, err = am.ActivateObject(
		ctx,
		core.RecordRef{},
		*obj,
		*am.GenesisRef(),
		*cb.Prototypes["one"],
		false,
		goplugintestutils.CBORMarshal(t, &struct{}{}),
	)
	assert.NoError(t, err)

	_, err = executeMethod(ctx, lr, pm, *obj, 0, "Hello", goplugintestutils.CBORMarshal(t, []interface{}{}))
	assert.NoError(t, err, "contract call")

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

func (r *One) Hello() (string, error) {
	return r.GetPrototype().String(), nil
}
`
	ctx := context.TODO()
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()

	err := cb.Build(map[string]string{"one": code})
	assert.NoError(t, err)

	objID, err := am.RegisterRequest(ctx, &message.CallConstructor{})
	assert.NoError(t, err)
	obj := getRefFromID(objID)
	_, err = am.ActivateObject(
		ctx,
		core.RecordRef{},
		*obj,
		*am.GenesisRef(),
		*cb.Prototypes["one"],
		false,
		goplugintestutils.CBORMarshal(t, &struct{}{}),
	)
	assert.NoError(t, err)

	res, err := executeMethod(ctx, lr, pm, *obj, 0, "Hello", goplugintestutils.CBORMarshal(t, []interface{}{}))
	assert.NoError(t, err)

	resParsed := goplugintestutils.CBORUnMarshalToSlice(t, res.(*reply.CallMethod).Result)
	assert.Equal(t, cb.Prototypes["one"].String(), resParsed[0])
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

func (r *One) Kill() error {
	r.SelfDestruct()
	return nil
}
`
	ctx := context.TODO()
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()

	err := cb.Build(map[string]string{"one": code})
	assert.NoError(t, err)

	objID, err := am.RegisterRequest(ctx, &message.CallConstructor{})
	assert.NoError(t, err)
	obj := getRefFromID(objID)
	_, err = am.ActivateObject(
		ctx,
		core.RecordRef{}, *obj,
		*am.GenesisRef(),
		*cb.Prototypes["one"],
		false,
		goplugintestutils.CBORMarshal(t, &struct{}{}),
	)
	assert.NoError(t, err)

	_, err = executeMethod(ctx, lr, pm, *obj, 0, "Kill", goplugintestutils.CBORMarshal(t, []interface{}{}))
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

func (r *One) Panic() error {
	panic("haha")
	return nil
}
func (r *One) NotPanic() error {
	return nil
}
`
	ctx := context.TODO()
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()

	err := cb.Build(map[string]string{"one": code})
	assert.NoError(t, err)

	objID, err := am.RegisterRequest(ctx, &message.CallConstructor{})
	assert.NoError(t, err)
	obj := getRefFromID(objID)
	_, err = am.ActivateObject(
		ctx,
		core.RecordRef{}, *obj,
		*am.GenesisRef(),
		*cb.Prototypes["one"],
		false,
		goplugintestutils.CBORMarshal(t, &struct{}{}),
	)
	assert.NoError(t, err)

	_, err = executeMethod(ctx, lr, pm, *obj, 0, "Panic", goplugintestutils.CBORMarshal(t, []interface{}{}))
	assert.Error(t, err)

	_, err = executeMethod(ctx, lr, pm, *obj, 0, "NotPanic", goplugintestutils.CBORMarshal(t, []interface{}{}))
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

func (c *Contract) NewChilds(cnt int) (int, error) {
	s := 0
	for i := 1; i < cnt; i++ {
        child.New(i).AsChild(c.GetReference())
		s += i
	} 
	return s, nil
}

func (c *Contract) SumChilds() (int, error) {
	s := 0
	childs, err := c.GetChildrenTyped(child.GetPrototype())
	if err != nil {
		return 0, err
	}
	for _, chref := range childs {
		o := child.GetObject(chref)
		n, err := o.GetNum()
		if err != nil {
			return 0, err
		}
		s += n
	}
	return s, nil
}

func (c *Contract) GetChildRefs() (ret []string, err error) {
	childs, err := c.GetChildrenTyped(child.GetPrototype())
	if err != nil {
		return nil, err
	}

	for _, chref := range childs {
		ret = append(ret, chref.String())
	}
	return ret, nil
}
`
	goChild := `
package main
import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type Child struct {
	foundation.BaseContract
	Num int
}

func (c *Child) GetNum() (int, error) {
	return c.Num, nil
}


func New(n int) (*Child, error) {
	return &Child{Num: n}, nil
}
`
	ctx := context.TODO()
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()

	err := cb.Build(map[string]string{"child": goChild})
	assert.NoError(t, err)
	err = cb.Build(map[string]string{"contract": goContract})
	assert.NoError(t, err)

	domain := core.NewRefFromBase58("c1")
	contractID, err := am.RegisterRequest(ctx, &message.CallConstructor{PrototypeRef: core.NewRefFromBase58("dassads")})
	assert.NoError(t, err)
	contract := getRefFromID(contractID)
	_, err = am.ActivateObject(
		ctx,
		domain,
		*contract,
		*am.GenesisRef(),
		*cb.Prototypes["contract"],
		false,
		goplugintestutils.CBORMarshal(t, nil),
	)
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, contract, nil, "contract created")

	resp, err := executeMethod(ctx, lr, pm, *contract, 0, "NewChilds", goplugintestutils.CBORMarshal(t, []interface{}{10}))
	assert.NoError(t, err, "contract call")
	r := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	assert.Equal(t, []interface{}{uint64(45), nil}, r)

	resp, err = executeMethod(ctx, lr, pm, *contract, 0, "SumChilds", goplugintestutils.CBORMarshal(t, []interface{}{}))
	assert.NoError(t, err, "contract call")

	ValidateAllResults(t, ctx, lr)

	assert.NoError(t, err, "contract call")
	r = goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	assert.Equal(t, []interface{}{uint64(45), nil}, r)
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

func (c *Contract) Rand() (int, error) {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(77), nil
}
`
	ctx := context.TODO()
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()

	err := cb.Build(map[string]string{"contract": goContract})
	assert.NoError(t, err)

	domain := core.NewRefFromBase58("c1")
	contractID, err := am.RegisterRequest(ctx, &message.CallConstructor{PrototypeRef: core.NewRefFromBase58("dassads")})
	assert.NoError(t, err)
	contract := getRefFromID(contractID)
	_, err = am.ActivateObject(
		ctx,
		domain,
		*contract,
		*am.GenesisRef(),
		*cb.Prototypes["contract"],
		false,
		goplugintestutils.CBORMarshal(t, nil),
	)
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, contract, nil, "contract created")

	for i := 0; i < 5; i++ {
		_, err = executeMethod(ctx, lr, pm, *contract, uint64(i), "Rand", goplugintestutils.CBORMarshal(t, []interface{}{}))
		assert.NoError(t, err, "contract call")
	}

	ValidateAllResults(t, ctx, lr, *contract)
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
	friend, err := holder.AsChild(r.GetReference())
	if err != nil {
		return err
	}

	return friend.AnError()
}

func (r *One) NoError() error {
	holder := two.New()
	friend, err := holder.AsChild(r.GetReference())
	if err != nil {
		return err
	}

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
func New() (*Two, error) {
	return &Two{}, nil
}
func (r *Two) AnError() error {
	return errors.New("an error")
}
func (r *Two) NoError() error {
	return nil
}
`
	ctx := context.TODO()
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()

	err := cb.Build(map[string]string{
		"one": contractOneCode,
		"two": contractTwoCode,
	})
	assert.NoError(t, err)

	domain := core.NewRefFromBase58("c1")
	contractID, err := am.RegisterRequest(ctx, &message.CallConstructor{})
	assert.NoError(t, err)
	contract := getRefFromID(contractID)
	_, err = am.ActivateObject(
		ctx,
		domain,
		*contract,
		*am.GenesisRef(),
		*cb.Prototypes["one"],
		false,
		goplugintestutils.CBORMarshal(t, nil),
	)
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, contract, nil, "contract created")

	resp, err := executeMethod(ctx, lr, pm, *contract, 0, "AnError", goplugintestutils.CBORMarshal(t, []interface{}{}))
	assert.NoError(t, err, "contract call")

	ch := new(codec.CborHandle)
	res := []interface{}{&foundation.Error{}}
	err = codec.NewDecoderBytes(resp.(*reply.CallMethod).Result, ch).Decode(&res)
	assert.NoError(t, err, "contract call")
	assert.Equal(t, &foundation.Error{S: "an error"}, res[0])

	resp, err = executeMethod(ctx, lr, pm, *contract, 0, "NoError", goplugintestutils.CBORMarshal(t, []interface{}{}))
	assert.NoError(t, err, "contract call")

	ValidateAllResults(t, ctx, lr)

	r := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	assert.Equal(t, []interface{}{nil}, r)
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

func (r *One) Hello() (*string, error) {
	holder := two.New()
	friend, err := holder.AsChild(r.GetReference())
	if err != nil {
		return nil, err
	}

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
func New() (*Two, error) {
	return &Two{}, nil
}
func (r *Two) Hello() (*string, error) {
	return nil, nil
}
`
	ctx := context.TODO()
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()

	err := cb.Build(map[string]string{
		"one": contractOneCode,
		"two": contractTwoCode,
	})
	assert.NoError(t, err)

	domain := core.NewRefFromBase58("c1")
	contractID, err := am.RegisterRequest(ctx, &message.CallConstructor{})
	assert.NoError(t, err)
	contract := getRefFromID(contractID)
	_, err = am.ActivateObject(
		ctx,
		domain,
		*contract,
		*am.GenesisRef(),
		*cb.Prototypes["one"],
		false,
		goplugintestutils.CBORMarshal(t, nil),
	)
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, contract, nil, "contract created")

	resp, err := executeMethod(ctx, lr, pm, *contract, 0, "Hello", goplugintestutils.CBORMarshal(t, []interface{}{}))
	assert.NoError(t, err, "contract call")

	ValidateAllResults(t, ctx, lr)

	r := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	assert.Equal(t, []interface{}{nil, nil}, r)
}

type Caller struct {
	member string
	key    *ecdsa.PrivateKey
	lr     core.LogicRunner
	t      *testing.T
}

func (s *Caller) SignedCall(ctx context.Context, pm core.PulseManager, rootDomain core.RecordRef, method string, params []interface{}) interface{} {
	seed := make([]byte, 32)
	_, err := rand.Read(seed)
	assert.NoError(s.t, err)

	buf := goplugintestutils.CBORMarshal(s.t, params)

	args, err := core.MarshalArgs(
		core.NewRefFromBase58(s.member),
		method,
		buf,
		seed)

	assert.NoError(s.t, err)

	sign, err := cryptoHelper.Sign(args, s.key)
	assert.NoError(s.t, err)

	res, err := executeMethod(ctx, s.lr, pm, core.NewRefFromBase58(s.member), 0, "Call", goplugintestutils.CBORMarshal(s.t, []interface{}{rootDomain, method, buf, seed, sign}))
	assert.NoError(s.t, err, "contract call")

	var result interface{}
	var contractErr error
	err = signer.UnmarshalParams(res.(*reply.CallMethod).Result, &result, &contractErr)
	assert.NoError(s.t, err, "unmarshal answer")
	assert.NoError(s.t, contractErr)

	return result
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

	ctx := context.TODO()
	// TODO need use pulseManager to sync all refs
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()
	err = cb.Build(map[string]string{"member": string(memberCode), "allowance": string(allowanceCode), "wallet": string(walletCode), "rootdomain": string(rootDomainCode)})
	assert.NoError(t, err)

	// Initializing Root Domain
	rootDomainID, err := am.RegisterRequest(ctx, &message.GenesisRequest{Name: "c1"})
	assert.NoError(t, err)
	rootDomainRef := getRefFromID(rootDomainID)
	rootDomainDesc, err := am.ActivateObject(
		ctx,
		core.RecordRef{},
		*rootDomainRef,
		*am.GenesisRef(),
		*cb.Prototypes["rootdomain"],
		false,
		goplugintestutils.CBORMarshal(t, nil),
	)
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, rootDomainRef, nil, "contract created")

	// Creating Root member
	rootKey, err := cryptoHelper.GeneratePrivateKey()
	assert.NoError(t, err)
	rootPubKey, err := cryptoHelper.ExportPublicKey(&rootKey.PublicKey)
	assert.NoError(t, err)

	rootMemberID, err := am.RegisterRequest(ctx, &message.GenesisRequest{Name: "c2"})
	assert.NoError(t, err)
	rootMemberRef := getRefFromID(rootMemberID)

	m, err := member.New("root", rootPubKey)
	assert.NoError(t, err)

	_, err = am.ActivateObject(
		ctx,
		core.RecordRef{},
		*rootMemberRef,
		*rootDomainRef,
		*cb.Prototypes["member"],
		false,
		goplugintestutils.CBORMarshal(t, m),
	)
	assert.NoError(t, err)

	// Updating root domain with root member
	_, err = am.UpdateObject(ctx, core.RecordRef{}, core.RecordRef{}, rootDomainDesc, goplugintestutils.CBORMarshal(t, rootdomain.RootDomain{RootMember: *rootMemberRef}))
	assert.NoError(t, err)

	root := Caller{rootMemberRef.String(), rootKey, lr, t}

	// Creating Member1
	member1Key, err := cryptoHelper.GeneratePrivateKey()
	assert.NoError(t, err)
	member1PubKey, err := cryptoHelper.ExportPublicKey(&member1Key.PublicKey)
	assert.NoError(t, err)

	res1 := root.SignedCall(ctx, pm, *rootDomainRef, "CreateMember", []interface{}{"Member1", member1PubKey})
	member1Ref := res1.(string)
	assert.NotEqual(t, "", member1Ref)

	// Creating Member2
	member2Key, err := cryptoHelper.GeneratePrivateKey()
	assert.NoError(t, err)
	member2PubKey, err := cryptoHelper.ExportPublicKey(&member2Key.PublicKey)
	assert.NoError(t, err)

	res2 := root.SignedCall(ctx, pm, *rootDomainRef, "CreateMember", []interface{}{"Member2", member2PubKey})
	member2Ref := res2.(string)
	assert.NotEqual(t, "", member2Ref)

	// Transfer 1 coin from Member1 to Member2
	member1 := Caller{member1Ref, member1Key, lr, t}
	member1.SignedCall(ctx, pm, *rootDomainRef, "Transfer", []interface{}{1, member2Ref})

	// Verify Member1 balance
	res3 := root.SignedCall(ctx, pm, *rootDomainRef, "GetBalance", []interface{}{member1Ref})
	assert.Equal(t, 999, int(res3.(uint64)))

	// Verify Member2 balance
	res4 := root.SignedCall(ctx, pm, *rootDomainRef, "GetBalance", []interface{}{member2Ref})
	assert.Equal(t, 1001, int(res4.(uint64)))
}

func TestFullValidationCycle(t *testing.T) {
	t.Skip("test is terribly wrong")
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

func (c *Contract) NewChilds(cnt int) (int, error) {
	s := 0
	for i := 1; i < cnt; i++ {
        child.New(i).AsChild(c.GetReference())
		s += i
	} 
	return s, nil
}

func (c *Contract) SumChilds() (int, error) {
	s := 0
	childs, err := c.GetChildrenTyped(child.GetImage())
	if err != nil {
		return 0, err
	}
	for _, chref := range childs {
		o := child.GetObject(chref)
		n, err := o.GetNum()
		if err != nil {
			return 0, err
		}
		s += n
	}
	return s, nil
}

func (c *Contract) GetChildRefs() (ret []string, err error) {
	childs, err := c.GetChildrenTyped(child.GetImage())
	if err != nil {
		return nil, err
	}

	for _, chref := range childs {
		ret = append(ret, chref.String())
	}
	return ret, nil
}
`
	goChild := `
package main
import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type Child struct {
	foundation.BaseContract
	Num int
}

func (c *Child) GetNum() (int, error) {
	return c.Num, nil
}


func New(n int) (*Child, error) {
	return &Child{Num: n}, nil
}
`
	ctx := context.TODO()
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()

	err := cb.Build(map[string]string{"child": goChild, "contract": goContract})
	assert.NoError(t, err)

	domain := core.NewRefFromBase58("c1")
	contractID, err := am.RegisterRequest(ctx, &message.CallConstructor{PrototypeRef: core.NewRefFromBase58("dassads")})
	assert.NoError(t, err)
	contract := getRefFromID(contractID)
	_, err = am.ActivateObject(
		ctx,
		domain,
		*contract,
		*am.GenesisRef(),
		*cb.Prototypes["contract"],
		false,
		goplugintestutils.CBORMarshal(t, nil),
	)
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, contract, nil, "contract created")

	resp, err := executeMethod(ctx, lr, pm, *contract, 0, "NewChilds", goplugintestutils.CBORMarshal(t, []interface{}{1}))
	assert.NoError(t, err, "contract call")
	r := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	assert.Equal(t, []interface{}{uint64(0), nil}, r)

	mb := lr.(*LogicRunner).MessageBus.(*testmessagebus.TestMessageBus)
	toValidate := make([]core.Parcel, 0)
	mb.ReRegister(core.TypeValidateCaseBind, func(ctx context.Context, m core.Parcel) (core.Reply, error) {
		toValidate = append(toValidate, m)
		return nil, nil
	})
	toExecute := make([]core.Parcel, 0)
	mb.ReRegister(core.TypeExecutorResults, func(ctx context.Context, m core.Parcel) (core.Reply, error) {
		toExecute = append(toExecute, m)
		return nil, nil
	})
	toCheckValidate := make([]core.Parcel, 0)
	mb.ReRegister(core.TypeValidationResults, func(ctx context.Context, m core.Parcel) (core.Reply, error) {
		toCheckValidate = append(toCheckValidate, m)
		return nil, nil
	})

	err = lr.(*LogicRunner).Ledger.GetPulseManager().Set(
		ctx,
		core.Pulse{PulseNumber: 1231234, Entropy: core.Entropy{}},
	)
	assert.NoError(t, err)

	for _, m := range toValidate {
		lr.ValidateCaseBind(ctx, m)
	}

	for _, m := range toExecute {
		lr.ExecutorResults(ctx, m)
	}

	for _, m := range toCheckValidate {
		lr.ProcessValidationResults(ctx, m)
	}
}

func TestConstructorReturnNil(t *testing.T) {
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

func (r *One) Hello() (*string, error) {
	holder := two.New()
	_, err := holder.AsChild(r.GetReference())
	if err != nil {
		return nil, err
	}
	ok := "all was well"
	return &ok, nil
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
func New() (*Two, error) {
	return nil, nil
}
// Contract without methods can't build because of import error in proxy
// TODO: INS-737
func (r *Two) Hello() (*string, error) {
	return nil, nil
}
`
	ctx := context.TODO()
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()

	err := cb.Build(map[string]string{
		"one": contractOneCode,
		"two": contractTwoCode,
	})
	assert.NoError(t, err)

	domain := core.NewRefFromBase58("c1")
	contractID, err := am.RegisterRequest(ctx, &message.CallConstructor{})
	assert.NoError(t, err)
	contract := getRefFromID(contractID)
	_, err = am.ActivateObject(
		ctx,
		domain,
		*contract,
		*am.GenesisRef(),
		*cb.Prototypes["one"],
		false,
		goplugintestutils.CBORMarshal(t, nil),
	)
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, contract, nil, "contract created")

	resp, err := executeMethod(ctx, lr, pm, *contract, 0, "Hello", goplugintestutils.CBORMarshal(t, []interface{}{}))
	assert.NoError(t, err, "contract call")

	var result interface{}
	var contractErr *foundation.Error

	err = signer.UnmarshalParams(resp.(*reply.CallMethod).Result, &result, &contractErr)
	assert.NoError(t, err, "unmarshal answer")
	assert.NotNil(t, contractErr)
	assert.Contains(t, contractErr.Error(), "[ FakeNew ] ( INSCONSTRUCTOR_* ) ( Generated Method ) Constructor returns nil")
}

func TestRecursiveCall(t *testing.T) {
	if parallel {
		t.Parallel()
	}

	var recursiveContractCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/application/proxy/recursive"
)
type One struct {
	foundation.BaseContract
}

func New() (*One, error) {
	return &One{}, nil
}

func (r *One) Recursive() (error) {
	remoteSelf := recursive.GetObject(r.GetReference())
	err := remoteSelf.Recursive()
	return err
}

`

	ctx := inslogger.ContextWithTrace(context.Background(), utils.RandTraceID())
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()

	err := cb.Build(map[string]string{
		"recursive": recursiveContractCode,
	})
	assert.NoError(t, err)

	domain := core.NewRefFromBase58("c1")
	contractID, err := am.RegisterRequest(ctx, &message.CallConstructor{PrototypeRef: core.NewRefFromBase58("recursive")})
	assert.NoError(t, err)
	contract := getRefFromID(contractID)
	_, err = am.ActivateObject(
		ctx, domain, *contract, *am.GenesisRef(), *cb.Prototypes["recursive"], false,
		goplugintestutils.CBORMarshal(t, nil),
	)
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, contract, nil, "contract created")

	resp, err := executeMethod(ctx, lr, pm, *contract, 0, "Recursive", goplugintestutils.CBORMarshal(t, []interface{}{}))
	assert.NoError(t, err, "contract call")
	r := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	assert.Equal(t, []interface{}{map[interface{}]interface{}{"S": "on calling main API: couldn't dispatch event: loop detected"}}, r)
}

func TestNewAllowanceNotFromWallet(t *testing.T) {
	if parallel {
		t.Parallel()
	}
	var contractOneCode = `
package main
import (
	"fmt"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/application/proxy/allowance"
	"github.com/insolar/insolar/application/proxy/wallet"
	"github.com/insolar/insolar/core"
)
type One struct {
	foundation.BaseContract
}
func (r *One) CreateAllowance(member string) (error) {
	w, _ := wallet.GetImplementationFrom(core.NewRefFromBase58(member))
	walletRef := w.GetReference()
	ah := allowance.New(&walletRef, 111, r.GetContext().Time.Unix()+10)
	_, err := ah.AsChild(walletRef)
	if err != nil {
		return fmt.Errorf("Error:", err.Error())
	}
	return nil
}
`
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

	ctx := context.TODO()
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()
	err = cb.Build(map[string]string{"one": contractOneCode, "member": string(memberCode), "allowance": string(allowanceCode), "wallet": string(walletCode), "rootdomain": string(rootDomainCode)})
	assert.NoError(t, err)

	// Initializing Root Domain
	rootDomainID, err := am.RegisterRequest(ctx, &message.GenesisRequest{Name: "c1"})
	assert.NoError(t, err)
	rootDomainRef := getRefFromID(rootDomainID)
	rootDomainDesc, err := am.ActivateObject(
		ctx,
		core.RecordRef{},
		*rootDomainRef,
		*am.GenesisRef(),
		*cb.Prototypes["rootdomain"],
		false,
		goplugintestutils.CBORMarshal(t, nil),
	)
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, rootDomainRef, nil, "contract created")

	// Creating Root member
	rootKey, err := cryptoHelper.GeneratePrivateKey()
	assert.NoError(t, err)
	rootPubKey, err := cryptoHelper.ExportPublicKey(&rootKey.PublicKey)
	assert.NoError(t, err)

	rootMemberID, err := am.RegisterRequest(ctx, &message.GenesisRequest{Name: "c2"})
	assert.NoError(t, err)
	rootMemberRef := getRefFromID(rootMemberID)

	m, err := member.New("root", rootPubKey)
	assert.NoError(t, err)

	_, err = am.ActivateObject(
		ctx,
		core.RecordRef{},
		*rootMemberRef,
		*rootDomainRef,
		*cb.Prototypes["member"],
		false,
		goplugintestutils.CBORMarshal(t, m),
	)
	assert.NoError(t, err)

	// Updating root domain with root member
	_, err = am.UpdateObject(ctx, core.RecordRef{}, core.RecordRef{}, rootDomainDesc, goplugintestutils.CBORMarshal(t, rootdomain.RootDomain{RootMember: *rootMemberRef}))
	assert.NoError(t, err)

	root := Caller{rootMemberRef.String(), rootKey, lr, t}

	// Creating Member
	memberKey, err := cryptoHelper.GeneratePrivateKey()
	assert.NoError(t, err)
	memberPubKey, err := cryptoHelper.ExportPublicKey(&memberKey.PublicKey)
	assert.NoError(t, err)

	res1 := root.SignedCall(ctx, pm, *rootDomainRef, "CreateMember", []interface{}{"Member", memberPubKey})
	memberRef := res1.(string)
	assert.NotEqual(t, "", memberRef)

	// Call CreateAllowance method in custom contract
	domain := core.NewRefFromBase58("c1")
	contractID, err := am.RegisterRequest(ctx, &message.CallConstructor{})
	assert.NoError(t, err)
	contract := getRefFromID(contractID)
	_, err = am.ActivateObject(
		ctx,
		domain,
		*contract,
		*am.GenesisRef(),
		*cb.Prototypes["one"],
		false,
		goplugintestutils.CBORMarshal(t, nil),
	)
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, contract, nil, "contract created")

	resp, err := executeMethod(ctx, lr, pm, *contract, 0, "CreateAllowance", goplugintestutils.CBORMarshal(t, []interface{}{memberRef}))
	assert.NoError(t, err, "contract call")

	var contractErr *foundation.Error

	err = signer.UnmarshalParams(resp.(*reply.CallMethod).Result, &contractErr)
	assert.NoError(t, err, "unmarshal answer")
	assert.NotNil(t, contractErr)
	assert.Contains(t, contractErr.Error(), "[ New Allowance ] : Can't create allowance from not wallet contract")

	// Verify Member balance
	res3 := root.SignedCall(ctx, pm, *rootDomainRef, "GetBalance", []interface{}{memberRef})
	assert.Equal(t, 1000, int(res3.(uint64)))
}

func TestGetParent(t *testing.T) {
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
}
 func (r *One) AddChildAndReturnMyselfAsParent() (core.RecordRef, error) {
	holder := two.New()
	friend, err := holder.AsChild(r.GetReference())
	if err != nil {
		return core.RecordRef{}, err
	}

 	return friend.GetParent()
}
`
	var contractTwoCode = `
package main
 import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)
 type Two struct {
	foundation.BaseContract
}
 func New() (*Two, error) {
	return &Two{}, nil
}
 func (r *Two) GetParent() (core.RecordRef, error) {
	return *r.GetContext().Parent, nil
}
 `
	ctx := context.Background()
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()
	err := cb.Build(map[string]string{"one": contractOneCode, "two": contractTwoCode})
	assert.NoError(t, err)
	objID, err := am.RegisterRequest(ctx, &message.CallConstructor{})
	assert.NoError(t, err)
	obj := getRefFromID(objID)
	_, err = am.ActivateObject(
		ctx,
		core.RecordRef{},
		*obj,
		*am.GenesisRef(),
		*cb.Prototypes["one"],
		false,
		goplugintestutils.CBORMarshal(t, &struct{}{}),
	)
	assert.NoError(t, err)
	resp, err := executeMethod(ctx, lr, pm, *obj, 0, "AddChildAndReturnMyselfAsParent", goplugintestutils.CBORMarshal(t, []interface{}{}))
	r := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)

	refFromMethod := r.([]interface{})[0].([]byte)
	assert.Equal(t, *obj, Ref{}.FromSlice(refFromMethod))
	ValidateAllResults(t, ctx, lr)
}

func TestReleaseRequestsAfterPulse(t *testing.T) {
	if parallel {
		t.Parallel()
	}

	var sleepContract = `
package main

import (
   "github.com/insolar/insolar/logicrunner/goplugin/foundation"
   "time"
)
type One struct {
   foundation.BaseContract
   N int
}

func New() (*One, error){
   return nil, nil
}

func (r *One) LongSleep() (error) {
   time.Sleep(7 * time.Second)
   r.N++
   return nil
}

func (r *One) ShortSleep() (error) {
   time.Sleep(1 * time.Microsecond)
   r.N++
   return nil
}

`
	ctx := inslogger.ContextWithTrace(context.Background(), utils.RandTraceID())
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()

	err := cb.Build(map[string]string{
		"one": sleepContract,
	})
	assert.NoError(t, err)

	domain := core.NewRefFromBase58("c1")
	contractID, err := am.RegisterRequest(ctx, &message.CallConstructor{PrototypeRef: core.NewRefFromBase58("one")})
	assert.NoError(t, err)
	contract := getRefFromID(contractID)
	_, err = am.ActivateObject(
		ctx, domain, *contract, *am.GenesisRef(), *cb.Prototypes["one"], false,
		goplugintestutils.CBORMarshal(t, nil),
	)
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, contract, nil, "contract created")

	lr = getLogicRunnerWithoutValidation(lr)

	// hold executor
	go func() {
		log.Debugf("!!!!! Long start")
		executeMethod(ctx, lr, pm, *contract, 0, "LongSleep", goplugintestutils.CBORMarshal(t, []interface{}{}))
		log.Debugf("!!!!! Long end")
	}()

	// wait both method calls, send new pulse
	go func() {
		log.Debugf("!!!!! Pulse sleep")
		time.Sleep(3 * time.Second)
		log.Debugf("!!!!! Pulse start")
		err = pm.Set(
			ctx,
			core.Pulse{PulseNumber: 1, Entropy: core.Entropy{}},
		)
		log.Debugf("!!!!! Pulse end")
	}()

	// wait for holding and add to queue
	log.Debugf("!!!!! Short sleep")
	time.Sleep(time.Second)
	log.Debugf("!!!!! Short start")
	_, err = executeMethod(ctx, lr, pm, *contract, 0, "ShortSleep", goplugintestutils.CBORMarshal(t, []interface{}{}))
	log.Debugf("!!!!! Short end")

	// TODO check for 302
	assert.Error(t, err, "contract call")
	assert.Equal(t, "Abort execution: New Pulse coming", err.Error())
}

func getLogicRunnerWithoutValidation(lr core.LogicRunner) *LogicRunner {
	rlr := lr.(*LogicRunner)
	newmb := rlr.MessageBus.(*testmessagebus.TestMessageBus)

	emptyFunc := func(context.Context, core.Parcel) (res core.Reply, err error) {
		return nil, nil
	}

	newmb.ReRegister(core.TypeValidationResults, emptyFunc)
	newmb.ReRegister(core.TypeExecutorResults, emptyFunc)

	rlr.MessageBus = newmb

	return rlr
}
