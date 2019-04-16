//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package logicrunner

import (
	"context"
	"crypto"
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/ugorji/go/codec"

	"github.com/insolar/insolar/application/contract/member"
	"github.com/insolar/insolar/application/contract/member/signer"
	"github.com/insolar/insolar/application/contract/rootdomain"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/contractrequester"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/delegationtoken"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage/drop"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/goplugin"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/insolar/insolar/logicrunner/pulsemanager"
	"github.com/insolar/insolar/messagebus"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/nodekeeper"
	"github.com/insolar/insolar/testutils/testmessagebus"
)

var parallel = false

type LogicRunnerFuncSuite struct {
	suite.Suite

	runnerBin    string
	icc          string
	contractsDir string
}

func FindContractsDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("failed to find folder")
	}

	// we're located in /logicrunner/logicrunner_test.go file, so we must take two basenames
	projectRoot := path.Dir(path.Dir(file))
	contractsDir := path.Join(projectRoot, "application", "contract")
	return contractsDir, nil
}

func (s *LogicRunnerFuncSuite) SetupSuite() {
	var err error
	if s.runnerBin, s.icc, err = goplugintestutils.Build(); err != nil {
		s.Fail("Logic runner build failed, skip tests: ", err.Error())
	}

	if s.contractsDir, err = FindContractsDir(); err != nil {
		s.contractsDir = ""
		log.Error("Failed to find contracts dir: ", err.Error())
	}
}

func MessageBusTrivialBehavior(mb *testmessagebus.TestMessageBus, lr insolar.LogicRunner) {
	mb.ReRegister(insolar.TypeCallMethod, lr.Execute)
	mb.ReRegister(insolar.TypeCallConstructor, lr.Execute)
	mb.ReRegister(insolar.TypeValidateCaseBind, lr.HandleValidateCaseBindMessage)
	mb.ReRegister(insolar.TypeValidationResults, lr.HandleValidationResultsMessage)
	mb.ReRegister(insolar.TypeExecutorResults, lr.HandleExecutorResultsMessage)
}

func (s *LogicRunnerFuncSuite) PrepareLrAmCbPm() (insolar.LogicRunner, artifacts.Client, *goplugintestutils.ContractsBuilder, insolar.PulseManager, func()) {
	ctx := context.TODO()
	lrSock := os.TempDir() + "/" + testutils.RandomString() + ".sock"
	rundSock := os.TempDir() + "/" + testutils.RandomString() + ".sock"

	rundCleaner, err := goplugintestutils.StartInsgorund(s.runnerBin, "unix", rundSock, "unix", lrSock)
	s.NoError(err)

	lr, err := NewLogicRunner(&configuration.LogicRunner{
		RPCListen:   lrSock,
		RPCProtocol: "unix",
		GoPlugin: &configuration.GoPlugin{
			RunnerListen:   rundSock,
			RunnerProtocol: "unix",
		},
	})
	s.NoError(err, "Initialize runner")

	mock := testutils.NewCryptographyServiceMock(s.T())
	mock.SignFunc = func(p []byte) (r *insolar.Signature, r1 error) {
		signature := insolar.SignatureFromBytes(nil)
		return &signature, nil
	}
	mock.GetPublicKeyFunc = func() (crypto.PublicKey, error) {
		return nil, nil
	}

	delegationTokenFactory := delegationtoken.NewDelegationTokenFactory()
	nk := nodekeeper.GetTestNodekeeper(mock)

	mb := testmessagebus.NewTestMessageBus(s.T())

	nw := testutils.GetTestNetwork(s.T())
	// FIXME: TmpLedger is deprecated. Use mocks instead.
	l, db, cleaner := artifacts.TmpLedger(
		s.T(),
		"",
		insolar.Components{
			LogicRunner: lr,
			NodeNetwork: nk,
			MessageBus:  mb,
			Network:     nw,
		},
	)

	indexMock := recentstorage.NewRecentIndexStorageMock(s.T())
	indexMock.AddObjectMock.Return()

	providerMock := recentstorage.NewProviderMock(s.T())
	providerMock.GetIndexStorageMock.Return(indexMock)
	providerMock.DecreaseIndexesTTLMock.Return(nil)

	parcelFactory := messagebus.NewParcelFactory()
	cm := &component.Manager{}
	cm.Register(platformpolicy.NewPlatformCryptographyScheme())
	am := l.GetArtifactManager()
	cm.Register(am, l.GetPulseManager(), l.GetJetCoordinator())
	cr, err := contractrequester.New()
	pulseAccessor := l.PulseManager.(*pulsemanager.PulseManager).PulseAccessor
	nth := testutils.NewTerminationHandlerMock(s.T())

	cm.Inject(db, pulseAccessor, nk, providerMock, l, lr, nw, mb, cr, delegationTokenFactory, parcelFactory, nth, mock)
	err = cm.Init(ctx)
	s.NoError(err)
	err = cm.Start(ctx)
	s.NoError(err)

	MessageBusTrivialBehavior(mb, lr)
	pm := l.GetPulseManager()

	s.incrementPulseHelper(ctx, lr, pm)

	cb := goplugintestutils.NewContractBuilder(am, s.icc)

	return lr, am, cb, pm, func() {
		cb.Clean()
		lr.Stop(ctx)
		cleaner()
		rundCleaner()
	}
}

func (s *LogicRunnerFuncSuite) incrementPulseHelper(ctx context.Context, lr insolar.LogicRunner, pm insolar.PulseManager) {
	pulseStorage := pm.(*pulsemanager.PulseManager).PulseAccessor
	currentPulse, _ := pulseStorage.Latest(ctx)

	newPulseNumber := currentPulse.PulseNumber + 1
	err := pm.Set(
		ctx,
		insolar.Pulse{PulseNumber: newPulseNumber, Entropy: insolar.Entropy{}},
		true,
	)
	s.Require().NoError(err)

	rootJetId := *insolar.NewJetID(0, nil)
	_, err = lr.(*LogicRunner).MessageBus.Send(
		ctx,
		&message.HotData{
			Jet:             *insolar.NewReference(insolar.DomainID, insolar.ID(rootJetId)),
			Drop:            drop.Drop{Pulse: 1, JetID: rootJetId},
			RecentObjects:   nil,
			PendingRequests: nil,
			PulseNumber:     newPulseNumber,
		}, nil,
	)
	s.Require().NoError(err)
}

func mockCryptographyService(t *testing.T) insolar.CryptographyService {
	mock := testutils.NewCryptographyServiceMock(t)
	mock.SignFunc = func(p []byte) (r *insolar.Signature, r1 error) {
		signature := insolar.SignatureFromBytes(nil)
		return &signature, nil
	}
	mock.VerifyFunc = func(p crypto.PublicKey, p1 insolar.Signature, p2 []byte) (r bool) {
		return true
	}
	return mock
}

func ValidateAllResults(t testing.TB, ctx context.Context, lr insolar.LogicRunner, mustfail ...insolar.Reference) {
	return // TODO REMOVE
	failmap := make(map[insolar.Reference]struct{})
	for _, r := range mustfail {
		failmap[r] = struct{}{}
	}

	rlr := lr.(*LogicRunner)

	a := assert.New(t)

	for ref, state := range rlr.state {
		log.Debugf("TEST validating: %s", ref)

		msg := state.ExecutionState.Behaviour.(*ValidationSaver).caseBind.ToValidateMessage(
			ctx, ref, *rlr.pulse(ctx),
		)
		cb := NewCaseBindFromValidateMessage(ctx, rlr.MessageBus, msg)

		_, err := rlr.Validate(ctx, ref, *rlr.pulse(ctx), *cb)
		if _, ok := failmap[ref]; ok {
			a.Error(err, "validation %s", ref)
		} else {
			a.NoError(err, "validation %s", ref)
		}
	}
}

func executeMethod(
	ctx context.Context, lr insolar.LogicRunner, pm insolar.PulseManager,
	objRef insolar.Reference, proxyPrototype insolar.Reference,
	nonce uint64,
	method string, arguments ...interface{},
) (
	insolar.Reply, error,
) {
	ctx = inslogger.ContextWithTrace(ctx, utils.RandTraceID())

	argsSerialized, err := insolar.Serialize(arguments)
	if err != nil {
		return nil, err
	}

	rlr := lr.(*LogicRunner)

	bm := message.BaseLogicMessage{
		Caller: testutils.RandomRef(),
		Nonce:  nonce,
	}

	rep, err := rlr.ContractRequester.CallMethod(ctx, &bm, false, false, &objRef, method, argsSerialized, &proxyPrototype)
	return rep, err
}

func firstMethodRes(t *testing.T, resp insolar.Reply) interface{} {
	res := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	return res.([]interface{})[0]
}

func (s *LogicRunnerFuncSuite) TestTypeCompatibilityError() {
	var _ insolar.LogicRunner = (*LogicRunner)(nil)
}

func getRefFromID(id *insolar.ID) *insolar.Reference {
	ref := insolar.Reference{}
	ref.SetRecord(*id)
	return &ref
}

func (s *LogicRunnerFuncSuite) TestSingleContractError() {
	if parallel {
		s.T().Parallel()
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

	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()

	err := cb.Build(map[string]string{"one": contractOneCode})
	s.NoError(err)

	obj, prototype := s.getObjectInstance(ctx, am, cb, "one")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Get")
	s.NoError(err, "contract call")
	s.Equal(uint64(0), firstMethodRes(s.T(), resp))

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Inc")
	s.NoError(err, "contract call")
	s.Equal(uint64(1), firstMethodRes(s.T(), resp))

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Get")
	s.NoError(err, "contract call")
	s.Equal(uint64(1), firstMethodRes(s.T(), resp))

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Dec")
	s.NoError(err, "contract call")
	s.Equal(uint64(0), firstMethodRes(s.T(), resp))

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Get")
	s.NoError(err, "contract call")
	s.Equal(uint64(0), firstMethodRes(s.T(), resp))

	ValidateAllResults(s.T(), ctx, lr)
}

func (s *LogicRunnerFuncSuite) TestContractCallingContractError() {
	if parallel {
		s.T().Parallel()
	}
	var contractOneCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import "github.com/insolar/insolar/application/proxy/two"
import "github.com/insolar/insolar/insolar"
import "errors"

type One struct {
	foundation.BaseContract
	Friend insolar.Reference
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

func (r *One)GetFriend() (insolar.Reference, error) {
	return r.Friend, nil
}

func (r *One)TestPayload() (two.Payload, error) {
	f := two.GetObject(r.Friend)
	err := f.SetPayload(two.Payload{Int: 10, Str: "HiHere"})
	if err != nil { return two.Payload{}, err }

	p, err := f.GetPayload()
	if err != nil { return two.Payload{}, err }

	str, err := f.GetPayloadString()	
	if err != nil { return two.Payload{}, err }

	if p.Str != str { return two.Payload{}, errors.New("Oops") }

	return p, nil

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
	P Payload
}

type Payload struct {
	Int int
	Str string
}

func New() (*Two, error) {
	return &Two{X:0}, nil;
}

func (r *Two) Hello(s string) (string, error) {
	r.X ++
	return fmt.Sprintf("Hello you too, %s. %d times!", s, r.X), nil
}

func (r *Two) GetPayload() (Payload, error) {
	return r.P, nil
}

func (r *Two) SetPayload(P Payload) (error) {
	r.P = P
	return nil
}

func (r *Two) GetPayloadString() (string, error) {
	return r.P.Str, nil
}
`
	ctx := context.Background()

	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()

	err := cb.Build(map[string]string{"one": contractOneCode, "two": contractTwoCode})
	s.NoError(err)

	obj, prototype := s.getObjectInstance(ctx, am, cb, "one")
	_, prototypeTwo := s.getObjectInstance(ctx, am, cb, "two")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Hello", "ins")
	s.NoError(err, "contract call")
	s.Equal("Hi, ins! Two said: Hello you too, ins. 1 times!", firstMethodRes(s.T(), resp))

	for i := 2; i <= 5; i++ {
		resp, err = executeMethod(ctx, lr, pm, *obj, *cb.Prototypes["one"], uint64(i), "Again", "ins")
		s.NoError(err, "contract call")
		s.Equal(
			fmt.Sprintf("Hi, ins! Two said: Hello you too, ins. %d times!", i),
			firstMethodRes(s.T(), resp),
		)
	}

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "GetFriend")
	s.NoError(err, "contract call")
	r0 := firstMethodRes(s.T(), resp).([]uint8)
	var two insolar.Reference
	for i := 0; i < 64; i++ {
		two[i] = r0[i]
	}

	for i := 6; i <= 9; i++ {
		resp, err = executeMethod(ctx, lr, pm, two, *prototypeTwo, uint64(i), "Hello", "Insolar")
		s.NoError(err, "contract call")
		s.Equal(fmt.Sprintf("Hello you too, Insolar. %d times!", i), firstMethodRes(s.T(), resp))
	}

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 7, "TestPayload")
	s.NoError(err, "contract call")
	res := firstMethodRes(s.T(), resp).(map[interface{}]interface{})["Str"]
	s.Equal("HiHere", res)

	ValidateAllResults(s.T(), ctx, lr)
}

func (s *LogicRunnerFuncSuite) TestInjectingDelegateError() {
	if parallel {
		s.T().Parallel()
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
	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()

	err := cb.Build(map[string]string{"one": contractOneCode, "two": contractTwoCode})
	s.NoError(err)

	obj, prototype := s.getObjectInstance(ctx, am, cb, "one")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Hello", "ins")
	s.NoError(err)
	s.Equal("Hi, ins! Two said: Hello you too, ins. 644 times!", firstMethodRes(s.T(), resp))

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "HelloFromDelegate", "ins")
	s.NoError(err)
	s.Equal("Hello you too, ins. 1288 times!", firstMethodRes(s.T(), resp))
}

func (s *LogicRunnerFuncSuite) TestBasicNotificationCallError() {
	if parallel {
		s.T().Parallel()
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
	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()

	err := cb.Build(map[string]string{"one": contractOneCode, "two": contractTwoCode})
	s.NoError(err)

	obj, prototype := s.getObjectInstance(ctx, am, cb, "one")

	_, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Hello")
	s.NoError(err, "contract call")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Value")
	s.NoError(err, "contract call")
	s.Equal(uint64(644), firstMethodRes(s.T(), resp))
}

func (s *LogicRunnerFuncSuite) TestContextPassingError() {
	if parallel {
		s.T().Parallel()
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
	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()

	err := cb.Build(map[string]string{"one": code})
	s.NoError(err)

	obj, prototype := s.getObjectInstance(ctx, am, cb, "one")

	res, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Hello")
	s.NoError(err)

	resParsed := goplugintestutils.CBORUnMarshalToSlice(s.T(), res.(*reply.CallMethod).Result)
	s.Equal(cb.Prototypes["one"].String(), resParsed[0])
}

func (s *LogicRunnerFuncSuite) TestDeactivationError() {
	if parallel {
		s.T().Parallel()
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
	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()

	err := cb.Build(map[string]string{"one": code})
	s.NoError(err)

	obj, prototype := s.getObjectInstance(ctx, am, cb, "one")

	_, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Kill")
	s.NoError(err, "contract call")
}

func (s *LogicRunnerFuncSuite) TestPanicError() {
	if parallel {
		s.T().Parallel()
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
	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()

	err := cb.Build(map[string]string{"one": code})
	s.NoError(err)

	obj, prototype := s.getObjectInstance(ctx, am, cb, "one")

	_, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Panic")
	s.Error(err)

	_, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "NotPanic")
	s.NoError(err)
}

func (s *LogicRunnerFuncSuite) TestGetChildrenError() {
	if parallel {
		s.T().Parallel()
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
	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()

	err := cb.Build(map[string]string{"child": goChild})
	s.NoError(err)
	err = cb.Build(map[string]string{"contract": goContract})
	s.NoError(err)

	obj, prototype := s.getObjectInstance(ctx, am, cb, "contract")

	// no childs, expect 0
	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "SumChildsByIterator")
	s.NoError(err, "empty children")
	s.Equal(uint64(0), firstMethodRes(s.T(), resp))

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "NewChilds", 10)
	s.NoError(err, "add children")
	s.Equal(uint64(45), firstMethodRes(s.T(), resp))

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "SumChildsByIterator")
	s.NoError(err, "sum real children")
	s.Equal(uint64(45), firstMethodRes(s.T(), resp))

	ValidateAllResults(s.T(), ctx, lr)
}

func (s *LogicRunnerFuncSuite) TestFailValidateError() {
	if parallel {
		s.T().Parallel()
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
	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()

	err := cb.Build(map[string]string{"contract": goContract})
	s.NoError(err)

	obj, prototype := s.getObjectInstance(ctx, am, cb, "contract")

	for i := 0; i < 5; i++ {
		_, err = executeMethod(ctx, lr, pm, *obj, *prototype, uint64(i), "Rand")
		s.NoError(err, "contract call")
	}

	ValidateAllResults(s.T(), ctx, lr, *obj)
}

func (s *LogicRunnerFuncSuite) TestErrorInterfaceError() {
	if parallel {
		s.T().Parallel()
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
	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()

	err := cb.Build(map[string]string{
		"one": contractOneCode,
		"two": contractTwoCode,
	})
	s.NoError(err)

	obj, prototype := s.getObjectInstance(ctx, am, cb, "one")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "AnError")
	s.NoError(err, "contract call")

	ch := new(codec.CborHandle)
	res := []interface{}{&foundation.Error{}}
	err = codec.NewDecoderBytes(resp.(*reply.CallMethod).Result, ch).Decode(&res)
	s.Equal(&foundation.Error{S: "an error"}, res[0])

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "NoError")
	s.NoError(err, "contract call")
	s.Nil(firstMethodRes(s.T(), resp))

	ValidateAllResults(s.T(), ctx, lr)
}

func (s *LogicRunnerFuncSuite) TestNilResultError() {
	if parallel {
		s.T().Parallel()
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
	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()

	err := cb.Build(map[string]string{
		"one": contractOneCode,
		"two": contractTwoCode,
	})
	s.NoError(err)

	obj, prototype := s.getObjectInstance(ctx, am, cb, "one")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Hello")
	s.NoError(err, "contract call")
	s.Nil(firstMethodRes(s.T(), resp))

	ValidateAllResults(s.T(), ctx, lr)
}

type Caller struct {
	member string
	lr     insolar.LogicRunner
	cs     insolar.CryptographyService
	suite  *LogicRunnerFuncSuite
}

func (s *Caller) SignedCall(ctx context.Context, pm insolar.PulseManager, rootDomain insolar.Reference, method string, proxyPrototype insolar.Reference, params []interface{}) interface{} {
	seed := make([]byte, 32)
	_, err := rand.Read(seed)
	s.suite.NoError(err)

	buf := goplugintestutils.CBORMarshal(s.suite.T(), params)

	memberRef, err := insolar.NewReferenceFromBase58(s.member)
	s.suite.Require().NoError(err)

	args, err := insolar.MarshalArgs(
		*memberRef,
		method,
		buf,
		seed)

	s.suite.NoError(err)

	signature, err := s.cs.Sign(args)
	s.suite.NoError(err)

	res, err := executeMethod(
		ctx, s.lr, pm, *memberRef, proxyPrototype, 0,
		"Call", rootDomain, method, buf, seed, signature.Bytes(),
	)
	s.suite.NoError(err, "contract call")

	var result interface{}
	var contractErr interface{}
	err = signer.UnmarshalParams(res.(*reply.CallMethod).Result, &result, &contractErr)
	s.suite.NoError(err, "unmarshal answer")
	s.suite.Require().Nilf(contractErr, "[ SignedCall ] Got error %v", contractErr)
	return result
}

func (s *LogicRunnerFuncSuite) LoadBasicContracts(contracts []string) map[string]string {
	contractCode := make(map[string]string)
	for _, contract := range contracts {
		code, err := ioutil.ReadFile(path.Join(s.contractsDir, contract, contract+".go"))
		if err != nil {
			s.Failf("Failed to load contract %s: %s", contract, err.Error())
		}
		contractCode[contract] = string(code)
	}
	return contractCode
}

func (s *LogicRunnerFuncSuite) TestRootDomainContractError() {
	if parallel {
		s.T().Parallel()
	}

	contracts := []string{"member", "allowance", "wallet", "rootdomain"}
	contractCode := s.LoadBasicContracts(contracts)
	ctx := context.TODO()
	// TODO need use pulseManager to sync all refs
	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()
	err := cb.Build(contractCode)
	s.NoError(err)

	// Initializing Root Domain
	rootDomainID, err := am.RegisterRequest(
		ctx,
		insolar.GenesisRecord.Ref(),
		&message.Parcel{
			Msg: &message.GenesisRequest{
				Name: "4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.7ZQboaH24PH42sqZKUvoa7UBrpuuubRtShp6CKNuWGZa",
			},
		},
	)
	s.NoError(err)
	rootDomainRef := getRefFromID(rootDomainID)
	rootDomainDesc, err := am.ActivateObject(
		ctx,
		insolar.Reference{},
		*rootDomainRef,
		insolar.GenesisRecord.Ref(),
		*cb.Prototypes["rootdomain"],
		false,
		goplugintestutils.CBORMarshal(s.T(), nil),
	)
	s.NoError(err, "create contract")
	s.NotEqual(rootDomainRef, nil, "contract created")

	kp := platformpolicy.NewKeyProcessor()

	// Creating Root member
	rootKey, err := kp.GeneratePrivateKey()
	s.NoError(err)
	rootPubKey, err := kp.ExportPublicKeyPEM(kp.ExtractPublicKey(rootKey))
	s.NoError(err)

	rootMemberID, err := am.RegisterRequest(
		ctx,
		insolar.GenesisRecord.Ref(),
		&message.Parcel{
			Msg: &message.GenesisRequest{
				Name: "4FFB8zfQoGznSmzDxwv4njX1aR9ioL8GHSH17QXH2AFa.7ZQboaH24PH42sqZKUvoa7UBrpuuubRtShp6CKNuWGZa",
			},
		},
	)
	s.NoError(err)
	rootMemberRef := getRefFromID(rootMemberID)

	m, err := member.New("root", string(rootPubKey))
	s.NoError(err)

	_, err = am.ActivateObject(
		ctx,
		insolar.Reference{},
		*rootMemberRef,
		*rootDomainRef,
		*cb.Prototypes["member"],
		false,
		goplugintestutils.CBORMarshal(s.T(), m),
	)
	s.NoError(err)

	// Updating root domain with root member
	_, err = am.UpdateObject(ctx, insolar.Reference{}, insolar.Reference{}, rootDomainDesc, goplugintestutils.CBORMarshal(s.T(), rootdomain.RootDomain{RootMember: *rootMemberRef}))
	s.NoError(err)

	csRoot := cryptography.NewKeyBoundCryptographyService(rootKey)
	root := Caller{rootMemberRef.String(), lr, csRoot, s}

	// Creating Member1
	member1Key, err := kp.GeneratePrivateKey()
	s.NoError(err)
	member1PubKey, err := kp.ExportPublicKeyPEM(kp.ExtractPublicKey(member1Key))
	s.NoError(err)

	res1 := root.SignedCall(ctx, pm, *rootDomainRef, "CreateMember", *cb.Prototypes["member"], []interface{}{"Member1", member1PubKey})
	member1Ref := res1.(string)
	s.NotEqual("", member1Ref)

	// Creating Member2
	member2Key, err := kp.GeneratePrivateKey()
	s.NoError(err)
	member2PubKey, err := kp.ExportPublicKeyPEM(kp.ExtractPublicKey(member2Key))
	s.NoError(err)

	res2 := root.SignedCall(ctx, pm, *rootDomainRef, "CreateMember", *cb.Prototypes["member"], []interface{}{"Member2", member2PubKey})
	member2Ref := res2.(string)
	s.NotEqual("", member2Ref)

	// Transfer 1 coin from Member1 to Member2
	csMember1 := cryptography.NewKeyBoundCryptographyService(member1Key)
	member1 := Caller{member1Ref, lr, csMember1, s}
	resTransfer := member1.SignedCall(ctx, pm, *rootDomainRef, "Transfer", *cb.Prototypes["member"], []interface{}{1, member2Ref})
	s.Nil(resTransfer)

	// Verify Member1 balance
	res3 := root.SignedCall(ctx, pm, *rootDomainRef, "GetBalance", *cb.Prototypes["member"], []interface{}{member1Ref})
	s.Equal(999999999, int(res3.(uint64)))

	// Verify Member2 balance
	res4 := root.SignedCall(ctx, pm, *rootDomainRef, "GetBalance", *cb.Prototypes["member"], []interface{}{member2Ref})
	s.Equal(1000000001, int(res4.(uint64)))
}

func (s *LogicRunnerFuncSuite) TestFullValidationCycleError() {
	s.T().Skip("test is terribly wrong")
	if parallel {
		s.T().Parallel()
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
	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()

	err := cb.Build(map[string]string{"child": goChild, "contract": goContract})
	s.NoError(err)

	obj, prototype := s.getObjectInstance(ctx, am, cb, "contract")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "NewChilds", 1)
	s.NoError(err, "contract call")
	s.Equal(uint64(0), firstMethodRes(s.T(), resp))

	mb := lr.(*LogicRunner).MessageBus.(*testmessagebus.TestMessageBus)
	toValidate := make([]insolar.Parcel, 0)
	mb.ReRegister(insolar.TypeValidateCaseBind, func(ctx context.Context, m insolar.Parcel) (insolar.Reply, error) {
		toValidate = append(toValidate, m)
		return nil, nil
	})
	toExecute := make([]insolar.Parcel, 0)
	mb.ReRegister(insolar.TypeExecutorResults, func(ctx context.Context, m insolar.Parcel) (insolar.Reply, error) {
		toExecute = append(toExecute, m)
		return nil, nil
	})
	toCheckValidate := make([]insolar.Parcel, 0)
	mb.ReRegister(insolar.TypeValidationResults, func(ctx context.Context, m insolar.Parcel) (insolar.Reply, error) {
		toCheckValidate = append(toCheckValidate, m)
		return nil, nil
	})

	newPulse := insolar.Pulse{PulseNumber: 1231234, Entropy: insolar.Entropy{}}

	err = lr.(*LogicRunner).MessageBus.(insolar.PulseManager).Set(
		ctx, newPulse, true,
	)
	s.NoError(err)

	for _, m := range toValidate {
		lr.HandleValidateCaseBindMessage(ctx, m)
	}

	for _, m := range toExecute {
		lr.HandleExecutorResultsMessage(ctx, m)
	}

	for _, m := range toCheckValidate {
		lr.HandleValidationResultsMessage(ctx, m)
	}
}

func (s *LogicRunnerFuncSuite) TestConstructorReturnNilError() {
	if parallel {
		s.T().Parallel()
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
	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()

	err := cb.Build(map[string]string{
		"one": contractOneCode,
		"two": contractTwoCode,
	})
	s.NoError(err)

	obj, prototype := s.getObjectInstance(ctx, am, cb, "one")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Hello")
	s.NoError(err, "contract call")

	var result interface{}
	var contractErr *foundation.Error

	err = signer.UnmarshalParams(resp.(*reply.CallMethod).Result, &result, &contractErr)
	s.NoError(err, "unmarshal answer")
	s.NotNil(contractErr)
	s.Contains(contractErr.Error(), "[ FakeNew ] ( INSCONSTRUCTOR_* ) ( Generated Method ) Constructor returns nil")
}

func (s *LogicRunnerFuncSuite) TestRecursiveCallError() {
	if parallel {
		s.T().Parallel()
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
	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()

	err := cb.Build(map[string]string{
		"recursive": recursiveContractCode,
	})
	s.NoError(err)

	obj, prototype := s.getObjectInstance(ctx, am, cb, "recursive")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "Recursive")
	s.NoError(err, "contract call")

	var contractErr *foundation.Error
	err = signer.UnmarshalParams(resp.(*reply.CallMethod).Result, &contractErr)
	s.NoError(err, "unmarshal answer")
	s.NotNil(contractErr)
	s.Contains(contractErr.Error(), "loop detected")
}

func (s *LogicRunnerFuncSuite) TestNewAllowanceNotFromWalletError() {
	if parallel {
		s.T().Parallel()
	}
	var contractOneCode = `
package main
import (
	"fmt"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/application/proxy/allowance"
	"github.com/insolar/insolar/application/proxy/wallet"
	"github.com/insolar/insolar/insolar"
)
type One struct {
	foundation.BaseContract
}
func (r *One) CreateAllowance(member string) (error) {
	memberRef, refErr := insolar.NewReferenceFromBase58(member)
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
	contracts := []string{"member", "allowance", "wallet", "rootdomain"}
	contractCode := s.LoadBasicContracts(contracts)
	contractCode["one"] = contractOneCode

	ctx := context.TODO()
	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()
	err := cb.Build(contractCode)
	s.NoError(err)

	kp := platformpolicy.NewKeyProcessor()

	// Initializing Root Domain
	rootDomainID, err := am.RegisterRequest(ctx, insolar.GenesisRecord.Ref(), &message.Parcel{Msg: &message.GenesisRequest{Name: "4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.7ZQboaH24PH42sqZKUvoa7UBrpuuubRtShp6CKNuWGZa"}})
	s.NoError(err)
	rootDomainRef := getRefFromID(rootDomainID)
	rootDomainDesc, err := am.ActivateObject(
		ctx,
		insolar.Reference{},
		*rootDomainRef,
		insolar.GenesisRecord.Ref(),
		*cb.Prototypes["rootdomain"],
		false,
		goplugintestutils.CBORMarshal(s.T(), nil),
	)
	s.NoError(err, "create contract")
	s.NotEqual(rootDomainRef, nil, "contract created")

	// Creating Root member
	rootKey, err := kp.GeneratePrivateKey()
	s.NoError(err)
	rootPubKey, err := kp.ExportPublicKeyPEM(kp.ExtractPublicKey(rootKey))
	s.NoError(err)

	rootMemberID, err := am.RegisterRequest(
		ctx,
		insolar.GenesisRecord.Ref(),
		&message.Parcel{
			Msg: &message.GenesisRequest{
				Name: "4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.4FFB8zfQoGznSmzDxwv4njX1aR9ioL8GHSH17QXH2AFa",
			},
		},
	)
	s.NoError(err)
	rootMemberRef := getRefFromID(rootMemberID)

	m, err := member.New("root", string(rootPubKey))
	s.NoError(err)

	_, err = am.ActivateObject(
		ctx,
		insolar.Reference{},
		*rootMemberRef,
		*rootDomainRef,
		*cb.Prototypes["member"],
		false,
		goplugintestutils.CBORMarshal(s.T(), m),
	)
	s.NoError(err)

	// Updating root domain with root member
	_, err = am.UpdateObject(ctx, insolar.Reference{}, insolar.Reference{}, rootDomainDesc, goplugintestutils.CBORMarshal(s.T(), rootdomain.RootDomain{RootMember: *rootMemberRef}))
	s.NoError(err)

	cs := cryptography.NewKeyBoundCryptographyService(rootKey)
	root := Caller{rootMemberRef.String(), lr, cs, s}

	// Creating Member
	memberKey, err := kp.GeneratePrivateKey()
	s.NoError(err)
	memberPubKey, err := kp.ExportPublicKeyPEM(kp.ExtractPublicKey(memberKey))
	s.NoError(err)

	res1 := root.SignedCall(ctx, pm, *rootDomainRef, "CreateMember", *cb.Prototypes["member"], []interface{}{"Member", string(memberPubKey)})
	memberRef := res1.(string)
	s.NotEqual("", memberRef)

	// Call CreateAllowance method in custom contract
	domain, err := insolar.NewReferenceFromBase58("7ZQboaH24PH42sqZKUvoa7UBrpuuubRtShp6CKNuWGZa.7ZQboaH24PH42sqZKUvoa7UBrpuuubRtShp6CKNuWGZa")
	s.Require().NoError(err)
	contractID, err := am.RegisterRequest(ctx, insolar.GenesisRecord.Ref(), &message.Parcel{Msg: &message.CallConstructor{}})
	s.NoError(err)
	contract := getRefFromID(contractID)
	_, err = am.ActivateObject(
		ctx,
		*domain,
		*contract,
		insolar.GenesisRecord.Ref(),
		*cb.Prototypes["one"],
		false,
		goplugintestutils.CBORMarshal(s.T(), nil),
	)
	s.NoError(err, "create contract")
	s.NotEqual(contract, nil, "contract created")

	resp, err := executeMethod(ctx, lr, pm, *contract, *cb.Prototypes["one"], 0, "CreateAllowance", memberRef)
	s.NoError(err, "contract call")

	var contractErr *foundation.Error

	err = signer.UnmarshalParams(resp.(*reply.CallMethod).Result, &contractErr)
	s.NoError(err, "unmarshal answer")
	s.NotNil(contractErr)
	s.Contains(contractErr.Error(), "[ New Allowance ] : Can't create allowance from not wallet contract")

	// Verify Member balance
	res3 := root.SignedCall(ctx, pm, *rootDomainRef, "GetBalance", *cb.Prototypes["member"], []interface{}{memberRef})
	s.Equal(1000000000, int(res3.(uint64)))
}

func (s *LogicRunnerFuncSuite) TestGetParentError() {
	if parallel {
		s.T().Parallel()
	}
	var contractOneCode = `
package main
 import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import "github.com/insolar/insolar/application/proxy/two"
import "github.com/insolar/insolar/insolar"
 type One struct {
	foundation.BaseContract
}
 func (r *One) AddChildAndReturnMyselfAsParent() (insolar.Reference, error) {
	holder := two.New()
	friend, err := holder.AsChild(r.GetReference())
	if err != nil {
		return insolar.Reference{}, err
	}

 	return friend.GetParent()
}
`
	var contractTwoCode = `
package main
 import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)
 type Two struct {
	foundation.BaseContract
}
 func New() (*Two, error) {
	return &Two{}, nil
}
 func (r *Two) GetParent() (insolar.Reference, error) {
	return *r.GetContext().Parent, nil
}
 `
	ctx := context.Background()
	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()
	err := cb.Build(map[string]string{"one": contractOneCode, "two": contractTwoCode})
	s.NoError(err)

	obj, prototype := s.getObjectInstance(ctx, am, cb, "one")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "AddChildAndReturnMyselfAsParent")
	s.Equal(*obj, Ref{}.FromSlice(firstMethodRes(s.T(), resp).([]byte)))

	ValidateAllResults(s.T(), ctx, lr)
}

func (s *LogicRunnerFuncSuite) TestReleaseRequestsAfterPulseError() {
	s.T().Skip("Test for old architecture. Unskip when new queue mechanism will release.")
	if parallel {
		s.T().Parallel()
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
	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()

	err := cb.Build(map[string]string{
		"one": sleepContract,
	})
	s.NoError(err)

	obj, prototype := s.getObjectInstance(ctx, am, cb, "one")

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
			insolar.Pulse{PulseNumber: 1, Entropy: insolar.Entropy{}},
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
	s.Error(err, "contract call")

	s.Contains(err.Error(), "abort execution: new Pulse coming")
}

func getLogicRunnerWithoutValidation(lr insolar.LogicRunner) *LogicRunner {
	rlr := lr.(*LogicRunner)
	newmb := rlr.MessageBus.(*testmessagebus.TestMessageBus)

	emptyFunc := func(context.Context, insolar.Parcel) (res insolar.Reply, err error) {
		return nil, nil
	}

	newmb.ReRegister(insolar.TypeValidationResults, emptyFunc)
	newmb.ReRegister(insolar.TypeExecutorResults, emptyFunc)

	rlr.MessageBus = newmb

	return rlr
}

func (s *LogicRunnerFuncSuite) TestGinsiderMustDieAfterInsolardError() {
	if parallel {
		s.T().Parallel()
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
	lr, am, cb, _, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()

	err := cb.Build(map[string]string{
		"one": emptyMethodContract,
	})
	s.NoError(err)

	_, prototype := s.getObjectInstance(ctx, am, cb, "one")

	proto, err := am.GetObject(ctx, *prototype, nil, false)
	codeRef, err := proto.Code()

	s.NoError(err, "get contract code")

	rlr := lr.(*LogicRunner)
	gp, err := goplugin.NewGoPlugin(rlr.Cfg, rlr.MessageBus, rlr.ArtifactManager)

	callContext := &insolar.LogicCallContext{
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
		Arguments: goplugintestutils.CBORMarshal(s.T(), []interface{}{}),
	}

	// emulate death
	err = rlr.sock.Close()
	s.Require().NoError(err)

	client, err := gp.Downstream(ctx)

	// call method without waiting of it execution
	client.Go("RPC.CallMethod", req, res, nil)

	// wait for gorund try to send answer back, it will see closing connection, after that it needs to die
	// ping to goPlugin, it has to be dead
	for start := time.Now(); time.Since(start) < time.Minute; {
		time.Sleep(100 * time.Millisecond)
		_, err = rpc.Dial(gp.Cfg.GoPlugin.RunnerProtocol, gp.Cfg.GoPlugin.RunnerListen)
		if err != nil {
			break
		}

		log.Debug("TestGinsiderMustDieAfterInsolard: gorund still alive")
	}

	s.Require().Error(err, "rpc Dial")
	s.Contains(err.Error(), "connect: connection refused")
}

func (s *LogicRunnerFuncSuite) TestGetRemoteDataError() {
	if parallel {
		s.T().Parallel()
	}
	var contractOneCode = `
package main
 import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
 import "github.com/insolar/insolar/application/proxy/two"
 import "github.com/insolar/insolar/insolar"
 type One struct {
	foundation.BaseContract
 }
 func (r *One) GetChildCode() (insolar.Reference, error) {
	holder := two.New()
	child, err := holder.AsChild(r.GetReference())
	if err != nil {
		return insolar.Reference{}, err
	}

 	return child.GetCode()
 }

 func (r *One) GetChildPrototype() (insolar.Reference, error) {
	holder := two.New()
	child, err := holder.AsChild(r.GetReference())
	if err != nil {
		return insolar.Reference{}, err
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
	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()
	err := cb.Build(map[string]string{"one": contractOneCode, "two": contractTwoCode})
	s.NoError(err)

	obj, prototype := s.getObjectInstance(ctx, am, cb, "one")

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "GetChildCode")
	s.NoError(err, "contract call")
	s.Equal(*cb.Codes["two"], Ref{}.FromSlice(firstMethodRes(s.T(), resp).([]byte)), "Compare Code Refs")

	resp, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "GetChildPrototype")
	s.NoError(err, "contract call")
	s.Equal(*cb.Prototypes["two"], Ref{}.FromSlice(firstMethodRes(s.T(), resp).([]byte)), "Compare Code prototypes")
}

func (s *LogicRunnerFuncSuite) TestNoLoopsWhileNotificationCallError() {
	if parallel {
		s.T().Parallel()
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
	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()
	err := cb.Build(map[string]string{"one": contractOneCode, "two": contractTwoCode})
	s.NoError(err)

	obj, prototype := s.getObjectInstance(ctx, am, cb, "one")

	for i := 0; i < 100; i++ {

	}

	resp, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "GetChildCode", goplugintestutils.CBORMarshal(s.T(), []interface{}{}))
	s.NoError(err, "contract call")
	r := goplugintestutils.CBORUnMarshal(s.T(), resp.(*reply.CallMethod).Result)
	s.Equal([]interface{}{uint64(100), nil}, r)
}

func (s *LogicRunnerFuncSuite) TestPrototypeMismatchError() {
	if parallel {
		s.T().Parallel()
	}
	testContract := `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/application/proxy/first"
	"github.com/insolar/insolar/insolar"
)

type Contract struct {
	foundation.BaseContract
}

func (c *Contract) Test(firstRef *insolar.Reference) (string, error) {
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
	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()

	err := cb.Build(map[string]string{"test": testContract, "first": firstContract, "second": secondContract})
	s.NoError(err)

	testObj, testPrototype := s.getObjectInstance(ctx, am, cb, "test")
	secondObj, _ := s.getObjectInstance(ctx, am, cb, "second")

	s.NoError(err, "create contract")
	s.NotEqual(secondObj, nil, "contract created")

	resp, err := executeMethod(ctx, lr, pm, *testObj, *testPrototype, 0, "Test", *secondObj)
	s.NoError(err, "contract call")

	ch := new(codec.CborHandle)
	res := []interface{}{&foundation.Error{}}
	err = codec.NewDecoderBytes(resp.(*reply.CallMethod).Result, ch).Decode(&res)
	s.Equal(map[interface{}]interface{}(map[interface{}]interface{}{"S": "[ RouteCall ] on calling main API: CallMethod returns error: proxy call error: try to call method of prototype as method of another prototype"}), res[1])
}

func (s *LogicRunnerFuncSuite) getObjectInstance(ctx context.Context, am artifacts.Client, cb *goplugintestutils.ContractsBuilder, contractName string) (*insolar.Reference, *insolar.Reference) {
	domain, err := insolar.NewReferenceFromBase58("4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.7ZQboaH24PH42sqZKUvoa7UBrpuuubRtShp6CKNuWGZa")
	s.Require().NoError(err)
	contractID, err := am.RegisterRequest(
		ctx,
		insolar.GenesisRecord.Ref(),
		&message.Parcel{Msg: &message.CallConstructor{PrototypeRef: testutils.RandomRef()}},
	)
	s.NoError(err)
	objectRef := getRefFromID(contractID)

	_, err = am.ActivateObject(
		ctx,
		*domain,
		*objectRef,
		insolar.GenesisRecord.Ref(),
		*cb.Prototypes[contractName],
		false,
		goplugintestutils.CBORMarshal(s.T(), nil),
	)
	s.NoError(err, "create contract")
	s.NotEqual(objectRef, nil, "contract created")

	return objectRef, cb.Prototypes[contractName]
}

func TestLogicRunnerFunc(t *testing.T) {
	if err := log.SetLevel("debug"); err != nil {
		log.Error("Failed to set logLevel to debug: ", err.Error())
	}

	t.Parallel()
	suite.Run(t, new(LogicRunnerFuncSuite))
}

func (s *LogicRunnerFuncSuite) TestImmutableAnnotation() {
	if parallel {
		s.T().Parallel()
	}
	var codeOne = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import "github.com/insolar/insolar/application/proxy/two"

type One struct {
	foundation.BaseContract
}

func (r *One) ExternalImmutableCall() (int, error) {
	holder := two.New()
	objTwo, err := holder.AsChild(r.GetReference())
	if err != nil {
		return 0, err
	}
	return objTwo.ReturnNumberAsImmutable()
}

func (r *One) ExternalImmutableCallMakesExternalCall() (error) {
	holder := two.New()
	objTwo, err := holder.AsChild(r.GetReference())
	if err != nil {
		return err
	}
	return objTwo.Immutable()
}
`

	var codeTwo = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import "github.com/insolar/insolar/application/proxy/three"

type Two struct {
	foundation.BaseContract
}

func New() (*Two, error) {
	return &Two{}, nil
}

func (r *Two) ReturnNumber() (int, error) {
	return 42, nil
}

//ins:immutable
func (r *Two) Immutable() (error) {
	holder := three.New()
	objThree, err := holder.AsChild(r.GetReference())
	if err != nil {
		return err
	}
	return objThree.DoNothing()
}

`

	var codeThree = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type Three struct {
	foundation.BaseContract
}

func New() (*Three, error) {
	return &Three{}, nil
}

func (r *Three) DoNothing() (error) {
	return nil
}

`

	ctx := context.TODO()
	lr, am, cb, pm, cleaner := s.PrepareLrAmCbPm()
	defer cleaner()

	err := cb.Build(map[string]string{"one": codeOne, "two": codeTwo, "three": codeThree})
	s.NoError(err)

	obj, prototype := s.getObjectInstance(ctx, am, cb, "one")

	res, err := executeMethod(ctx, lr, pm, *obj, *prototype, 0, "ExternalImmutableCall")
	s.NoError(err)
	resParsed := goplugintestutils.CBORUnMarshalToSlice(s.T(), res.(*reply.CallMethod).Result)
	s.Equal(uint64(42), resParsed[0])
	s.Equal(nil, resParsed[1])

	res, err = executeMethod(ctx, lr, pm, *obj, *prototype, 0, "ExternalImmutableCallMakesExternalCall")
	s.NoError(err, "contract call")
	resParsed = goplugintestutils.CBORUnMarshalToSlice(s.T(), res.(*reply.CallMethod).Result)
	s.Equal(map[interface{}]interface{}{"S": "[ RouteCall ] on calling main API: Try to call route from immutable method"}, resParsed[0])
}
