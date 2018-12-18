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
	"crypto"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
	"testing"
	"time"

	"github.com/insolar/insolar/ledger/pulsemanager"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/logicrunner/goplugin"

	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"

	"github.com/insolar/insolar/cryptography"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/messagebus"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils/network"
	"github.com/insolar/insolar/testutils/nodekeeper"

	"github.com/insolar/insolar/application/contract/member"
	"github.com/insolar/insolar/application/contract/member/signer"
	"github.com/insolar/insolar/application/contract/rootdomain"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/delegationtoken"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/core/utils"
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
var parallel = false

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

func PrepareLrAmCbPm(t *testing.T) (core.LogicRunner, core.ArtifactManager, *goplugintestutils.ContractsBuilder, core.PulseManager, func()) {
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

	mock := testutils.NewCryptographyServiceMock(t)
	mock.SignFunc = func(p []byte) (r *core.Signature, r1 error) {
		signature := core.SignatureFromBytes(nil)
		return &signature, nil
	}
	mock.GetPublicKeyFunc = func() (crypto.PublicKey, error) {
		return nil, nil
	}

	delegationTokenFactory := delegationtoken.NewDelegationTokenFactory()
	nk := nodekeeper.GetTestNodekeeper(mock)

	mb := testmessagebus.NewTestMessageBus(t)
	mb.PulseNumber = 0

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

	pulseStorage := l.PulseManager.(*pulsemanager.PulseManager).PulseStorage
	recentMock := recentstorage.NewProviderMock(t)

	parcelFactory := messagebus.NewParcelFactory()
	cm := &component.Manager{}
	cm.Register(platformpolicy.NewPlatformCryptographyScheme())
	am := l.GetArtifactManager()
	cm.Register(am, l.GetPulseManager(), l.GetJetCoordinator())

	cm.Inject(pulseStorage, nk, recentMock, l, lr, nw, mb, delegationTokenFactory, parcelFactory, mock)
	err = cm.Init(ctx)
	assert.NoError(t, err)
	err = cm.Start(ctx)
	assert.NoError(t, err)

	MessageBusTrivialBehavior(mb, lr)
	pm := l.GetPulseManager()

	currentPulse, _ := pulseStorage.Current(ctx)
	newPulseNumber := currentPulse.PulseNumber + 1
	err = lr.Ledger.GetPulseManager().Set(
		ctx,
		core.Pulse{PulseNumber: newPulseNumber, Entropy: core.Entropy{}},
		true,
	)
	require.NoError(t, err)

	mb.PulseNumber = newPulseNumber

	assert.NoError(t, err)
	if err != nil {
		t.Fatal("pulse set died, ", err)
	}
	cb := goplugintestutils.NewContractBuilder(am, icc)

	return lr, am, cb, pm, func() {
		cb.Clean()
		lr.Stop(ctx)
		cleaner()
		rundCleaner()
	}
}

func mockCryptographyService(t *testing.T) core.CryptographyService {
	mock := testutils.NewCryptographyServiceMock(t)
	mock.SignFunc = func(p []byte) (r *core.Signature, r1 error) {
		signature := core.SignatureFromBytes(nil)
		return &signature, nil
	}
	mock.VerifyFunc = func(p crypto.PublicKey, p1 core.Signature, p2 []byte) (r bool) {
		return true
	}
	return mock
}

func ValidateAllResults(t testing.TB, ctx context.Context, lr core.LogicRunner, mustfail ...core.RecordRef) {
	failmap := make(map[core.RecordRef]struct{})
	for _, r := range mustfail {
		failmap[r] = struct{}{}
	}

	rlr := lr.(*LogicRunner)

	for ref, state := range rlr.state {
		log.Debugf("TEST validating: %s", ref)

		msg := state.ExecutionState.Behaviour.(*ValidationSaver).caseBind.ToValidateMessage(
			ctx, ref, *rlr.pulse(ctx),
		)
		cb := NewCaseBindFromValidateMessage(ctx, rlr.MessageBus, msg)

		_, err := rlr.Validate(ctx, ref, *rlr.pulse(ctx), *cb)
		if _, ok := failmap[ref]; ok {
			assert.Error(t, err, "validation %s", ref)
		} else {
			assert.NoError(t, err, "validation %s", ref)
		}
	}
}

func executeMethod(
	ctx context.Context, lr core.LogicRunner, pm core.PulseManager,
	objRef core.RecordRef, proxyPrototype core.RecordRef,
	nonce uint64,
	method string, arguments ...interface{},
) (
	core.Reply, error,
) {
	argsSerialized, err := core.Serialize(arguments)
	if err != nil {
		return nil, err
	}

	msg := &message.CallMethod{
		ObjectRef:      objRef,
		Method:         method,
		Arguments:      argsSerialized,
		ProxyPrototype: proxyPrototype,
	}
	msg.Caller = testutils.RandomRef()
	if nonce != 0 {
		msg.Nonce = nonce
	}

	pf := lr.(*LogicRunner).ParcelFactory
	parcel, _ := pf.Create(ctx, msg, testutils.RandomRef(), nil, *core.GenesisPulse)
	ctx = inslogger.ContextWithTrace(ctx, utils.RandTraceID())
	resp, err := lr.Execute(
		ctx,
		parcel,
	)

	return resp, err
}

func firstMethodRes(t *testing.T, resp core.Reply) interface{} {
	res := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	return res.([]interface{})[0]
}

func TestTypeCompatibility(t *testing.T) {
	var _ core.LogicRunner = (*LogicRunner)(nil)
}

func getRefFromID(id *core.RecordID) *core.RecordRef {
	ref := core.RecordRef{}
	ref.SetRecord(*id)
	return &ref
}

func TestSingleContract(t *testing.T) {
	if parallel {
		t.Parallel()
	}
	var contractOneCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type One struct {
	foundation.BaseContract
	Number int
}

func (c *One) Inc() (int, error) {
	c.Number++
	return c.Number, nil
}

func (c *One) Get() (int, error) {
	return c.Number, nil
}

func (c *One) Dec() (int, error) {
	c.Number--
	return c.Number, nil
}
`
	ctx := context.Background()

	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()

	err := cb.Build(map[string]string{"one": contractOneCode})
	assert.NoError(t, err)

	obj, prototype := getObjectInstance(t, ctx, am, cb, "one")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Get")
	assert.NoError(t, err, "contract call")
	assert.Equal(t, uint64(0), firstMethodRes(t, resp))

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Inc")
	assert.NoError(t, err, "contract call")
	assert.Equal(t, uint64(1), firstMethodRes(t, resp))

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Get")
	assert.NoError(t, err, "contract call")
	assert.Equal(t, uint64(1), firstMethodRes(t, resp))

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Dec")
	assert.NoError(t, err, "contract call")
	assert.Equal(t, uint64(0), firstMethodRes(t, resp))

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Get")
	assert.NoError(t, err, "contract call")
	assert.Equal(t, uint64(0), firstMethodRes(t, resp))

	ValidateAllResults(t, ctx, lr)
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

	obj, prototype := getObjectInstance(t, ctx, am, cb, "one")
	_, prototypeTwo := getObjectInstance(t, ctx, am, cb, "two")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Hello", "ins")
	assert.NoError(t, err, "contract call")
	assert.Equal(t, "Hi, ins! Two said: Hello you too, ins. 1 times!", firstMethodRes(t, resp))

	for i := 2; i <= 5; i++ {
		resp, err = executeMethod(ctx, lr, pm, *obj, *cb.Prototypes["one"], uint64(i), "Again", "ins")
		assert.NoError(t, err, "contract call")
		assert.Equal(
			t,
			fmt.Sprintf("Hi, ins! Two said: Hello you too, ins. %d times!", i),
			firstMethodRes(t, resp),
		)
	}

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "GetFriend")
	assert.NoError(t, err, "contract call")
	r0 := firstMethodRes(t, resp).([]uint8)
	var two core.RecordRef
	for i := 0; i < 64; i++ {
		two[i] = r0[i]
	}

	for i := 6; i <= 9; i++ {
		resp, err = executeMethod(ctx, lr, pm, two, *prototypeTwo, uint64(i), "Hello", "Insolar")
		assert.NoError(t, err, "contract call")
		assert.Equal(t, fmt.Sprintf("Hello you too, Insolar. %d times!", i), firstMethodRes(t, resp))
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

	err := cb.Build(map[string]string{"one": contractOneCode, "two": contractTwoCode})
	assert.NoError(t, err)

	obj, prototype := getObjectInstance(t, ctx, am, cb, "one")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Hello", "ins")
	assert.NoError(t, err)
	assert.Equal(t, "Hi, ins! Two said: Hello you too, ins. 644 times!", firstMethodRes(t, resp))

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "HelloFromDelegate", "ins")
	assert.NoError(t, err)
	assert.Equal(t, "Hello you too, ins. 1288 times!", firstMethodRes(t, resp))
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

func (r *One) Value() (int, error) {
	friend, err := two.GetImplementationFrom(r.GetReference())
	if err != nil {
		return 0, err
	}

	return friend.Value()
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

func (r *Two) Value() (int, error) {
	return r.X, nil
}
`
	ctx := context.TODO()
	// TODO: use am := testutil.NewTestArtifactManager() here
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()

	err := cb.Build(map[string]string{"one": contractOneCode, "two": contractTwoCode})
	assert.NoError(t, err)

	obj, prototype := getObjectInstance(t, ctx, am, cb, "one")

	_, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Hello")
	assert.NoError(t, err, "contract call")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Value")
	assert.NoError(t, err, "contract call")
	assert.Equal(t, uint64(644), firstMethodRes(t, resp))
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

	obj, prototype := getObjectInstance(t, ctx, am, cb, "one")

	res, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Hello")
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

	obj, prototype := getObjectInstance(t, ctx, am, cb, "one")

	_, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Kill")
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

	obj, prototype := getObjectInstance(t, ctx, am, cb, "one")

	_, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Panic")
	assert.Error(t, err)

	_, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "NotPanic")
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

func (c *Contract) SumChildsByIterator() (int, error) {
	s := 0
	iterator, err := c.NewChildrenTypedIterator(child.GetPrototype())
	if err != nil {
		return 0, err
	}

	for iterator.HasNext() {
		chref, err := iterator.Next()
		if err != nil {
			return 0, err
		}

		o := child.GetObject(chref)
		n, err := o.GetNum()
		if err != nil {
			return 0, err
		}
		s += n
	}
	return s, nil
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

	obj, prototype := getObjectInstance(t, ctx, am, cb, "contract")

	// no childs, expect 0
	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "SumChildsByIterator")
	assert.NoError(t, err, "empty children")
	assert.Equal(t, uint64(0), firstMethodRes(t, resp))

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "NewChilds", 10)
	assert.NoError(t, err, "add children")
	assert.Equal(t, uint64(45), firstMethodRes(t, resp))

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "SumChildsByIterator")
	assert.NoError(t, err, "sum real children")
	assert.Equal(t, uint64(45), firstMethodRes(t, resp))

	ValidateAllResults(t, ctx, lr)
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

	obj, prototype := getObjectInstance(t, ctx, am, cb, "contract")

	for i := 0; i < 5; i++ {
		_, err = executeMethod(ctx, lr, pm, *obj, *prototype, uint64(i), "Rand")
		assert.NoError(t, err, "contract call")
	}

	ValidateAllResults(t, ctx, lr, *obj)
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

	obj, prototype := getObjectInstance(t, ctx, am, cb, "one")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "AnError")
	assert.NoError(t, err, "contract call")

	ch := new(codec.CborHandle)
	res := []interface{}{&foundation.Error{}}
	err = codec.NewDecoderBytes(resp.(*reply.CallMethod).Result, ch).Decode(&res)
	assert.Equal(t, &foundation.Error{S: "an error"}, res[0])

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "NoError")
	assert.NoError(t, err, "contract call")
	assert.Equal(t, nil, firstMethodRes(t, resp))

	ValidateAllResults(t, ctx, lr)
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

	obj, prototype := getObjectInstance(t, ctx, am, cb, "one")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Hello")
	assert.NoError(t, err, "contract call")
	assert.Equal(t, nil, firstMethodRes(t, resp))

	ValidateAllResults(t, ctx, lr)
}

type Caller struct {
	member string
	lr     core.LogicRunner
	t      *testing.T
	cs     core.CryptographyService
}

func (s *Caller) SignedCall(ctx context.Context, pm core.PulseManager, rootDomain core.RecordRef, method string, proxyPrototype core.RecordRef, params []interface{}) interface{} {
	seed := make([]byte, 32)
	_, err := rand.Read(seed)
	assert.NoError(s.t, err)

	buf := goplugintestutils.CBORMarshal(s.t, params)

	memberRef, err := core.NewRefFromBase58(s.member)
	require.NoError(s.t, err)

	args, err := core.MarshalArgs(
		*memberRef,
		method,
		buf,
		seed)

	assert.NoError(s.t, err)

	signature, err := s.cs.Sign(args)
	assert.NoError(s.t, err)

	res, err := executeMethod(
		ctx, s.lr, pm, *memberRef, proxyPrototype, 0,
		"Call", rootDomain, method, buf, seed, signature.Bytes(),
	)
	assert.NoError(s.t, err, "contract call")

	var result interface{}
	var contractErr interface{}
	err = signer.UnmarshalParams(res.(*reply.CallMethod).Result, &result, &contractErr)
	assert.NoError(s.t, err, "unmarshal answer")
	require.Nilf(s.t, contractErr, "[ SignedCall ] Got error %v", contractErr)
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
	err = cb.Build(map[string]string{
		"member":     string(memberCode),
		"allowance":  string(allowanceCode),
		"wallet":     string(walletCode),
		"rootdomain": string(rootDomainCode),
	})
	assert.NoError(t, err)

	// Initializing Root Domain
	rootDomainID, err := am.RegisterRequest(
		ctx,
		&message.Parcel{
			Msg: &message.GenesisRequest{
				Name: "4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.7ZQboaH24PH42sqZKUvoa7UBrpuuubRtShp6CKNuWGZa",
			},
		},
	)
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

	kp := platformpolicy.NewKeyProcessor()

	// Creating Root member
	rootKey, err := kp.GeneratePrivateKey()
	assert.NoError(t, err)
	rootPubKey, err := kp.ExportPublicKey(kp.ExtractPublicKey(rootKey))
	assert.NoError(t, err)

	rootMemberID, err := am.RegisterRequest(
		ctx,
		&message.Parcel{
			Msg: &message.GenesisRequest{
				Name: "4FFB8zfQoGznSmzDxwv4njX1aR9ioL8GHSH17QXH2AFa.7ZQboaH24PH42sqZKUvoa7UBrpuuubRtShp6CKNuWGZa",
			},
		},
	)
	assert.NoError(t, err)
	rootMemberRef := getRefFromID(rootMemberID)

	m, err := member.New("root", string(rootPubKey))
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

	csRoot := cryptography.NewKeyBoundCryptographyService(rootKey)
	root := Caller{rootMemberRef.String(), lr, t, csRoot}

	// Creating Member1
	member1Key, err := kp.GeneratePrivateKey()
	assert.NoError(t, err)
	member1PubKey, err := kp.ExportPublicKey(kp.ExtractPublicKey(member1Key))
	assert.NoError(t, err)

	res1 := root.SignedCall(ctx, pm, *rootDomainRef, "CreateMember", *cb.Prototypes["member"], []interface{}{"Member1", member1PubKey})
	member1Ref := res1.(string)
	assert.NotEqual(t, "", member1Ref)

	// Creating Member2
	member2Key, err := kp.GeneratePrivateKey()
	assert.NoError(t, err)
	member2PubKey, err := kp.ExportPublicKey(kp.ExtractPublicKey(member2Key))
	assert.NoError(t, err)

	res2 := root.SignedCall(ctx, pm, *rootDomainRef, "CreateMember", *cb.Prototypes["member"], []interface{}{"Member2", member2PubKey})
	member2Ref := res2.(string)
	assert.NotEqual(t, "", member2Ref)

	// Transfer 1 coin from Member1 to Member2
	csMember1 := cryptography.NewKeyBoundCryptographyService(member1Key)
	member1 := Caller{member1Ref, lr, t, csMember1}
	resTransfer := member1.SignedCall(ctx, pm, *rootDomainRef, "Transfer", *cb.Prototypes["member"], []interface{}{1, member2Ref})
	assert.Equal(t, nil, resTransfer)

	// Verify Member1 balance
	res3 := root.SignedCall(ctx, pm, *rootDomainRef, "GetBalance", *cb.Prototypes["member"], []interface{}{member1Ref})
	assert.Equal(t, 999, int(res3.(uint64)))

	// Verify Member2 balance
	res4 := root.SignedCall(ctx, pm, *rootDomainRef, "GetBalance", *cb.Prototypes["member"], []interface{}{member2Ref})
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

	obj, prototype := getObjectInstance(t, ctx, am, cb, "contract")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "NewChilds", 1)
	assert.NoError(t, err, "contract call")
	assert.Equal(t, uint64(0), firstMethodRes(t, resp))

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
		true,
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
`
	ctx := context.TODO()
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()

	err := cb.Build(map[string]string{
		"one": contractOneCode,
		"two": contractTwoCode,
	})
	assert.NoError(t, err)

	obj, prototype := getObjectInstance(t, ctx, am, cb, "one")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Hello")
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

	obj, prototype := getObjectInstance(t, ctx, am, cb, "recursive")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Recursive")
	assert.NoError(t, err, "contract call")

	var contractErr *foundation.Error
	err = signer.UnmarshalParams(resp.(*reply.CallMethod).Result, &contractErr)
	assert.NoError(t, err, "unmarshal answer")
	assert.NotNil(t, contractErr)
	assert.Contains(t, contractErr.Error(), "loop detected")
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
	memberRef, refErr := core.NewRefFromBase58(member)
	if refErr != nil {
		return refErr
	}
	w, _ := wallet.GetImplementationFrom(*memberRef)
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
	err = cb.Build(map[string]string{
		"one":        contractOneCode,
		"member":     string(memberCode),
		"allowance":  string(allowanceCode),
		"wallet":     string(walletCode),
		"rootdomain": string(rootDomainCode),
	})
	assert.NoError(t, err)

	kp := platformpolicy.NewKeyProcessor()

	// Initializing Root Domain
	rootDomainID, err := am.RegisterRequest(ctx, &message.Parcel{Msg: &message.GenesisRequest{Name: "4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.7ZQboaH24PH42sqZKUvoa7UBrpuuubRtShp6CKNuWGZa"}})
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
	rootKey, err := kp.GeneratePrivateKey()
	assert.NoError(t, err)
	rootPubKey, err := kp.ExportPublicKey(kp.ExtractPublicKey(rootKey))
	assert.NoError(t, err)

	rootMemberID, err := am.RegisterRequest(
		ctx,
		&message.Parcel{
			Msg: &message.GenesisRequest{
				Name: "4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.4FFB8zfQoGznSmzDxwv4njX1aR9ioL8GHSH17QXH2AFa",
			},
		},
	)
	assert.NoError(t, err)
	rootMemberRef := getRefFromID(rootMemberID)

	m, err := member.New("root", string(rootPubKey))
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

	cs := cryptography.NewKeyBoundCryptographyService(rootKey)
	root := Caller{rootMemberRef.String(), lr, t, cs}

	// Creating Member
	memberKey, err := kp.GeneratePrivateKey()
	assert.NoError(t, err)
	memberPubKey, err := kp.ExportPublicKey(kp.ExtractPublicKey(memberKey))
	assert.NoError(t, err)

	res1 := root.SignedCall(ctx, pm, *rootDomainRef, "CreateMember", *cb.Prototypes["member"], []interface{}{"Member", string(memberPubKey)})
	memberRef := res1.(string)
	assert.NotEqual(t, "", memberRef)

	// Call CreateAllowance method in custom contract
	domain, err := core.NewRefFromBase58("7ZQboaH24PH42sqZKUvoa7UBrpuuubRtShp6CKNuWGZa.7ZQboaH24PH42sqZKUvoa7UBrpuuubRtShp6CKNuWGZa")
	require.NoError(t, err)
	contractID, err := am.RegisterRequest(ctx, &message.Parcel{Msg: &message.CallConstructor{}})
	assert.NoError(t, err)
	contract := getRefFromID(contractID)
	_, err = am.ActivateObject(
		ctx,
		*domain,
		*contract,
		*am.GenesisRef(),
		*cb.Prototypes["one"],
		false,
		goplugintestutils.CBORMarshal(t, nil),
	)
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, contract, nil, "contract created")

	resp, err := executeMethod(ctx, lr, pm, *contract, *cb.Prototypes["one"], 0, "CreateAllowance", memberRef)
	assert.NoError(t, err, "contract call")

	var contractErr *foundation.Error

	err = signer.UnmarshalParams(resp.(*reply.CallMethod).Result, &contractErr)
	assert.NoError(t, err, "unmarshal answer")
	assert.NotNil(t, contractErr)
	assert.Contains(t, contractErr.Error(), "[ New Allowance ] : Can't create allowance from not wallet contract")

	// Verify Member balance
	res3 := root.SignedCall(ctx, pm, *rootDomainRef, "GetBalance", *cb.Prototypes["member"], []interface{}{memberRef})
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

	obj, prototype := getObjectInstance(t, ctx, am, cb, "one")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "AddChildAndReturnMyselfAsParent")
	assert.Equal(t, *obj, Ref{}.FromSlice(firstMethodRes(t, resp).([]byte)))

	ValidateAllResults(t, ctx, lr)
}

func TestReleaseRequestsAfterPulse(t *testing.T) {
	t.Skip("Test for old architecture. Unskip when new queue mechanism will release.")
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

	obj, prototype := getObjectInstance(t, ctx, am, cb, "one")

	lr = getLogicRunnerWithoutValidation(lr)

	// hold executor
	go func() {
		log.Debugf("!!!!! Long start")
		executeMethod(ctx, lr, pm, *obj, *prototype, 0, "LongSleep")
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
			true,
		)
		log.Debugf("!!!!! Pulse end")
	}()

	// wait for holding and add to queue
	log.Debugf("!!!!! Short sleep")
	time.Sleep(time.Second)
	log.Debugf("!!!!! Short start")
	_, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "ShortSleep")
	log.Debugf("!!!!! Short end")
	assert.Error(t, err, "contract call")

	assert.Contains(t, err.Error(), "abort execution: new Pulse coming")
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

func TestGinsiderMustDieAfterInsolard(t *testing.T) {
	if parallel {
		t.Parallel()
	}

	var emptyMethodContract = `
package main

import (
	"time"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	
)
type One struct {
   foundation.BaseContract
}

func New() (*One, error){
   return nil, nil
}

func (r *One) EmptyMethod() (error) {
	time.Sleep(200 * time.Millisecond)
	return nil
}

`
	ctx := inslogger.ContextWithTrace(context.Background(), utils.RandTraceID())
	lr, am, cb, _, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()

	err := cb.Build(map[string]string{
		"one": emptyMethodContract,
	})
	assert.NoError(t, err)

	_, prototype := getObjectInstance(t, ctx, am, cb, "one")

	proto, err := am.GetObject(ctx, *prototype, nil, false)
	codeRef, err := proto.Code()

	assert.NoError(t, err, "get contract code")

	rlr := lr.(*LogicRunner)
	gp, err := goplugin.NewGoPlugin(rlr.Cfg, rlr.MessageBus, rlr.ArtifactManager)

	callContext := &core.LogicCallContext{
		Caller:          nil,
		Callee:          nil,
		Request:         nil,
		Time:            time.Now(),
		Pulse:           *rlr.pulse(ctx),
		TraceID:         inslogger.TraceID(ctx),
		CallerPrototype: nil,
	}
	res := rpctypes.DownCallMethodResp{}
	req := rpctypes.DownCallMethodReq{
		Context:   callContext,
		Code:      *codeRef,
		Method:    "EmptyMethod",
		Arguments: goplugintestutils.CBORMarshal(t, []interface{}{}),
	}

	client, err := gp.Downstream(ctx)

	// call method without waiting of it execution
	client.Go("RPC.CallMethod", req, res, nil)

	// emulate death
	rlr.sock.Close()

	// wait for gorund try to send answer back, it will see closing connection, after that it needs to die
	time.Sleep(300 * time.Millisecond)

	// ping to goPlugin, it has to be dead
	_, err = rpc.Dial(gp.Cfg.GoPlugin.RunnerProtocol, gp.Cfg.GoPlugin.RunnerListen)
	assert.Error(t, err, "rpc Dial")
	assert.Contains(t, err.Error(), "connect: connection refused")
}

func TestGetRemoteData(t *testing.T) {
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
 func (r *One) GetChildCode() (core.RecordRef, error) {
	holder := two.New()
	child, err := holder.AsChild(r.GetReference())
	if err != nil {
		return core.RecordRef{}, err
	}

 	return child.GetCode()
 }

 func (r *One) GetChildPrototype() (core.RecordRef, error) {
	holder := two.New()
	child, err := holder.AsChild(r.GetReference())
	if err != nil {
		return core.RecordRef{}, err
	}

 	return child.GetPrototype()
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
 `
	ctx := context.Background()
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()
	err := cb.Build(map[string]string{"one": contractOneCode, "two": contractTwoCode})
	assert.NoError(t, err)

	obj, prototype := getObjectInstance(t, ctx, am, cb, "one")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "GetChildCode")
	assert.NoError(t, err, "contract call")
	assert.Equal(t, *cb.Codes["two"], Ref{}.FromSlice(firstMethodRes(t, resp).([]byte)), "Compare Code Refs")

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "GetChildPrototype")
	assert.NoError(t, err, "contract call")
	assert.Equal(t, *cb.Prototypes["two"], Ref{}.FromSlice(firstMethodRes(t, resp).([]byte)), "Compare Code Prototypes")
}

// TODO - unskip when we decide how to work with NotificationCalls (NoWaitMethods)
func TestNoLoopsWhileNotificationCall(t *testing.T) {
	t.Skip()
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
 func (r *One) GetChildCode() (int, error) {
	holder := two.New()
	child, err := holder.AsChild(r.GetReference())
	if err != nil {
		return 0, err
	}

	for i := 0; i < 100; i++ {
		child.IncreaseNoWait()
	}

 	return child.GetCounter()
 }
`
	var contractTwoCode = `
 package main
 import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
 )
 type Two struct {
	foundation.BaseContract
	Counter int
 }
 func New() (*Two, error) {
	return &Two{}, nil
 }

 func (r *Two) Increase() error {
 	r.Counter++
	return nil
 }

 func (r *Two) GetCounter() (int, error) {
	return r.Counter, nil
 }

`

	ctx := context.Background()
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()
	err := cb.Build(map[string]string{"one": contractOneCode, "two": contractTwoCode})
	assert.NoError(t, err)

	obj, prototype := getObjectInstance(t, ctx, am, cb, "one")

	for i := 0; i < 100; i++ {

	}

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "GetChildCode", goplugintestutils.CBORMarshal(t, []interface{}{}))
	assert.NoError(t, err, "contract call")
	r := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	assert.Equal(t, []interface{}{uint64(100), nil}, r)
}

func TestPrototypeMismatch(t *testing.T) {
	if parallel {
		t.Parallel()
	}
	testContract := `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/application/proxy/first"
	"github.com/insolar/insolar/core"
)

type Contract struct {
	foundation.BaseContract
}

func (c *Contract) Test(firstRef *core.RecordRef) (string, error) {
	return first.GetObject(*firstRef).GetName()
}
`

	// right contract
	firstContract := `
package main
import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type First struct {
	foundation.BaseContract
}

func (c *First) GetName() (string, error) {
	return "first", nil
}
`

	// malicious contract with same method signature and another behaviour
	secondContract := `
package main
import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type First struct {
	foundation.BaseContract
}

func (c *First) GetName() (string, error) {
	return "YOU ARE ROBBED!", nil
}
`
	ctx := context.TODO()
	lr, am, cb, pm, cleaner := PrepareLrAmCbPm(t)
	defer cleaner()

	err := cb.Build(map[string]string{"test": testContract, "first": firstContract, "second": secondContract})
	assert.NoError(t, err)

	testObj, testPrototype := getObjectInstance(t, ctx, am, cb, "test")
	secondObj, _ := getObjectInstance(t, ctx, am, cb, "second")

	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, secondObj, nil, "contract created")

	resp, err := executeMethod(ctx, lr, pm, *testObj, *testPrototype, 0, "Test", *secondObj)
	assert.NoError(t, err, "contract call")

	ch := new(codec.CborHandle)
	res := []interface{}{&foundation.Error{}}
	err = codec.NewDecoderBytes(resp.(*reply.CallMethod).Result, ch).Decode(&res)
	assert.Equal(t, map[interface{}]interface{}(map[interface{}]interface{}{"S": "[ RouteCall ] on calling main API: couldn't dispatch event: proxy call error: try to call method of prototype as method of another prototype"}), res[1])
}

func getObjectInstance(t *testing.T, ctx context.Context, am core.ArtifactManager, cb *goplugintestutils.ContractsBuilder, contractName string) (*core.RecordRef, *core.RecordRef) {
	domain, err := core.NewRefFromBase58("4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.7ZQboaH24PH42sqZKUvoa7UBrpuuubRtShp6CKNuWGZa")
	require.NoError(t, err)
	contractID, err := am.RegisterRequest(
		ctx,
		&message.Parcel{Msg: &message.CallConstructor{PrototypeRef: testutils.RandomRef()}},
	)
	assert.NoError(t, err)
	objectRef := getRefFromID(contractID)

	_, err = am.ActivateObject(
		ctx,
		*domain,
		*objectRef,
		*am.GenesisRef(),
		*cb.Prototypes[contractName],
		false,
		goplugintestutils.CBORMarshal(t, nil),
	)
	assert.NoError(t, err, "create contract")
	assert.NotEqual(t, objectRef, nil, "contract created")

	return objectRef, cb.Prototypes[contractName]
}
