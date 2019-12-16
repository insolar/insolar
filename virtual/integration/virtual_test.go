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
// +build slowtest

package integration

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/semaphore"

	"github.com/insolar/insolar/application/genesisrefs"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/machinesmanager"
	"github.com/insolar/insolar/testutils"
)

func TestVirtual_BasicOperations(t *testing.T) {
	t.Parallel()

	cfg := DefaultVMConfig()

	t.Run("happy path", func(t *testing.T) {
		ctx := inslogger.TestContext(t)

		expectedRes := struct {
			blip string
		}{
			blip: "blop",
		}

		mle := testutils.NewMachineLogicExecutorMock(t).CallMethodMock.Set(
			func(_ context.Context, _ *insolar.LogicCallContext, _ insolar.Reference, _ []byte, _ string, _ insolar.Arguments) ([]byte, insolar.Arguments, error) {
				return insolar.MustSerialize(expectedRes), insolar.MustSerialize(expectedRes), nil
			},
		)

		mm := machinesmanager.NewMachinesManagerMock(t).GetExecutorMock.Set(
			func(_ insolar.MachineType) (insolar.MachineLogicExecutor, error) {
				return mle, nil
			},
		).RegisterExecutorMock.Set(
			func(_ insolar.MachineType, _ insolar.MachineLogicExecutor) error {
				return nil
			},
		)

		s, err := NewVirtualServer(t, ctx, cfg).SetMachinesManager(mm).PrepareAndStart()
		require.NoError(t, err)
		defer s.Stop(ctx)

		// Prepare environment (mimic) for first call
		var objectID *insolar.ID
		{
			codeID, err := s.mimic.AddCode(ctx, []byte{})
			require.NoError(t, err)

			prototypeID, err := s.mimic.AddObject(ctx, *codeID, true, []byte{})
			require.NoError(t, err)

			objectID, err = s.mimic.AddObject(ctx, *prototypeID, false, []byte{})
			require.NoError(t, err)
		}
		t.Logf("iniitialization done")

		objectRef := insolar.NewReference(*objectID)

		res, requestRef, err := CallContract(
			s, objectRef, "good.CallMethod", nil, s.pulseGenerator.GetLastPulse().PulseNumber,
		)
		require.NoError(t, err)

		assert.NotEmpty(t, requestRef)
		assert.Equal(t, &reply.CallMethod{
			Object: objectRef,
			Result: insolar.MustSerialize(expectedRes),
		}, res)
	})

	t.Run("create user test", func(t *testing.T) {
		ctx := inslogger.TestContext(t)

		s, err := NewVirtualServer(t, ctx, cfg).WithGenesis().PrepareAndStart()
		require.NoError(t, err)
		defer s.Stop(ctx)

		user, err := NewUserWithKeys()
		if err != nil {
			panic("failed to create new user: " + err.Error())
		}

		callMethodReply, _, err := s.BasicAPICall(ctx, "member.create", nil, genesisrefs.ContractRootMember, user)
		if err != nil {
			panic(err.Error())
		}

		var result map[string]interface{}
		if err := insolar.Deserialize(callMethodReply.(*reply.CallMethod).Result, &result); err != nil {
			panic(err.Error())
		}

		assert.Nil(t, result["Error"])
		assert.NotNil(t, result["Returns"].([]interface{})[0].(map[string]interface{})["reference"])
	})

	t.Run("create and transfer test", func(t *testing.T) {
		ctx := inslogger.TestContext(t)

		s, err := NewVirtualServer(t, ctx, cfg).WithGenesis().PrepareAndStart()
		require.NoError(t, err)
		defer s.Stop(ctx)

		user1, err := NewUserWithKeys()
		if err != nil {
			panic("failed to create new user: " + err.Error())
		}
		user2, err := NewUserWithKeys()
		if err != nil {
			panic("failed to create new user: " + err.Error())
		}

		var walletReference1 insolar.Reference
		{
			callMethodReply, _, err := s.BasicAPICall(ctx, "member.create", nil, genesisrefs.ContractRootMember, user1)
			if err != nil {
				panic(err.Error())
			}

			var result map[string]interface{}
			if err := insolar.Deserialize(callMethodReply.(*reply.CallMethod).Result, &result); err != nil {
				panic(err.Error())
			}

			assert.Nil(t, result["Error"])

			walletReferenceString := result["Returns"].([]interface{})[0].(map[string]interface{})["reference"]
			assert.NotNil(t, walletReferenceString)
			assert.IsType(t, "", walletReferenceString)

			walletReference, err := insolar.NewReferenceFromString(walletReferenceString.(string))
			assert.NoError(t, err)

			walletReference1 = *walletReference
		}

		var walletReference2 insolar.Reference
		{
			callMethodReply, _, err := s.BasicAPICall(ctx, "member.create", nil, genesisrefs.ContractRootMember, user2)
			if err != nil {
				panic(err.Error())
			}

			var result map[string]interface{}
			if err := insolar.Deserialize(callMethodReply.(*reply.CallMethod).Result, &result); err != nil {
				panic(err.Error())
			}

			assert.Nil(t, result["Error"])
			assert.NotNil(t, result["Returns"].([]interface{})[0].(map[string]interface{})["reference"])

			walletReferenceString := result["Returns"].([]interface{})[0].(map[string]interface{})["reference"]
			assert.NotNil(t, walletReferenceString)
			assert.IsType(t, "", walletReferenceString)

			walletReference, err := insolar.NewReferenceFromString(walletReferenceString.(string))
			assert.NoError(t, err)

			walletReference2 = *walletReference
		}

		var feeWalletBalance string
		{
			callParams := map[string]interface{}{"reference": FeeWalletUser.Reference.String()}
			callMethodReply, _, err := s.BasicAPICall(ctx, "member.getBalance", callParams, FeeWalletUser.Reference, FeeWalletUser)
			if err != nil {
				panic(err.Error())
			}

			var result map[string]interface{}
			if err := insolar.Deserialize(callMethodReply.(*reply.CallMethod).Result, &result); err != nil {
				panic(err.Error())
			}
			require.Nil(t, result["Error"])

			fmt.Printf("%#v\n", result)
			feeWalletBalance = result["Returns"].([]interface{})[0].(map[string]interface{})["balance"].(string)
		}

		{
			callParams := map[string]interface{}{"amount": "10000", "toMemberReference": walletReference2.String()}
			callMethodReply, _, err := s.BasicAPICall(ctx, "member.transfer", callParams, walletReference1, user1)
			if err != nil {
				panic(err.Error())
			}

			var result map[string]interface{}
			if err := insolar.Deserialize(callMethodReply.(*reply.CallMethod).Result, &result); err != nil {
				panic(err.Error())
			}

			assert.Nil(t, result["Error"])
		}

		{
			for i := 1; i < 30; i++ {
				callParams := map[string]interface{}{"reference": FeeWalletUser.Reference.String()}
				callMethodReply, _, err := s.BasicAPICall(ctx, "member.getBalance", callParams, FeeWalletUser.Reference, FeeWalletUser)
				if err != nil {
					panic(err.Error())
				}

				var result map[string]interface{}
				if err := insolar.Deserialize(callMethodReply.(*reply.CallMethod).Result, &result); err != nil {
					panic(err.Error())
				}
				require.Nil(t, result["Error"])

				fmt.Printf("%#v\n", result)
				newBalance := result["Returns"].([]interface{})[0].(map[string]interface{})["balance"].(string)

				if newBalance != feeWalletBalance {
					break
				}

				time.Sleep(100 * time.Millisecond)
				if i == 29 {
					assert.FailNow(t, "failed to wait money in feeWallet")
				}
			}
		}
	})
}

type ServerHelper struct {
	s *Server
}

func (h *ServerHelper) createUser(ctx context.Context) (*User, error) {
	user, err := NewUserWithKeys()
	if err != nil {
		return nil, errors.Errorf("failed to create new user: " + err.Error())
	}

	{
		callMethodReply, _, err := h.s.BasicAPICall(ctx, "member.create", nil, genesisrefs.ContractRootMember, user)
		if err != nil {
			return nil, errors.Wrap(err, "failed to call member.create")
		}

		var result map[string]interface{}
		if cm, ok := callMethodReply.(*reply.CallMethod); !ok {
			return nil, errors.Wrapf(err, "unexpected type of return value %T", callMethodReply)
		} else if err := insolar.Deserialize(cm.Result, &result); err != nil {
			return nil, errors.Wrap(err, "failed to deserialize result")
		}

		r0, ok := result["Returns"]
		if ok && r0 != nil {
			if r1, ok := r0.([]interface{}); !ok {
				return nil, errors.Errorf("bad response: bad type of 'Returns' [%#v]", r0)
			} else if len(r1) != 2 {
				return nil, errors.Errorf("bad response: bad length of 'Returns' [%#v]", r0)
			} else if r2, ok := r1[0].(map[string]interface{}); !ok {
				return nil, errors.Errorf("bad response: bad type of first value [%#v]", r1)
			} else if r3, ok := r2["reference"]; !ok {
				return nil, errors.Errorf("bad response: absent reference field [%#v]", r2)
			} else if walletReferenceString, ok := r3.(string); !ok {
				return nil, errors.Errorf("bad response: reference field expected to be a string [%#v]", r3)
			} else if walletReference, err := insolar.NewReferenceFromString(walletReferenceString); err != nil {
				return nil, errors.Wrap(err, "bad response: got bad reference")
			} else {
				user.Reference = *walletReference
			}

			return user, nil
		}

		r0, ok = result["Error"]
		if ok && r0 != nil {
			return nil, errors.Errorf("%T: %#v", r0, r0)
		}

		panic("unreachable")
	}
}

func (h *ServerHelper) transferMoney(ctx context.Context, from User, to User, amount string) (string, error) {
	callParams := map[string]interface{}{
		"amount":            amount,
		"toMemberReference": to.Reference.String(),
	}
	callMethodReply, _, err := h.s.BasicAPICall(ctx, "member.transfer", callParams, from.Reference, &from)
	if err != nil {
		return "", errors.Wrap(err, "failed to call member.transfer")
	}

	var result map[string]interface{}
	if cm, ok := callMethodReply.(*reply.CallMethod); !ok {
		return "", errors.Wrapf(err, "unexpected type of return value %T", callMethodReply)
	} else if err := insolar.Deserialize(cm.Result, &result); err != nil {
		return "", errors.Wrap(err, "failed to deserialize result")
	}

	r0, ok := result["Returns"]
	if ok && r0 != nil {
		if r1, ok := r0.([]interface{}); !ok {
			return "", errors.Errorf("bad response: bad type of 'Returns' [%#v]", r0)
		} else if len(r1) != 2 {
			return "", errors.Errorf("bad response: bad length of 'Returns' [%#v]", r0)
		} else if r2, ok := r1[0].(map[string]interface{}); !ok {
			return "", errors.Errorf("bad response: bad type of first value [%#v]", r1)
		} else if r3, ok := r2["Fee"]; !ok {
			return "", errors.Errorf("bad response: absent Fee field [%#v]", r2)
		} else if fee, ok := r3.(string); !ok {
			return "", errors.Errorf("bad response: Fee field expected to be a string [%#v]", r3)
		} else {
			return fee, nil
		}
	}

	r0, ok = result["Error"]
	if ok && r0 != nil {
		return "", errors.Errorf("%T: %#v", r0, r0)
	}

	panic("unreachable")
}

func (h *ServerHelper) getBalance(ctx context.Context, user User) (string, error) {
	callParams := map[string]interface{}{
		"reference": user.Reference.String(),
	}
	callMethodReply, _, err := h.s.BasicAPICall(ctx, "member.getBalance", callParams, user.Reference, &user)
	if err != nil {
		return "", errors.Wrap(err, "failed to call member.getBalance")
	}

	var result map[string]interface{}
	if cm, ok := callMethodReply.(*reply.CallMethod); !ok {
		return "", errors.Wrapf(err, "unexpected type of return value %T", callMethodReply)
	} else if err := insolar.Deserialize(cm.Result, &result); err != nil {
		return "", errors.Wrap(err, "failed to deserialize result")
	}

	r0, ok := result["Returns"]
	if ok && r0 != nil {
		if r1, ok := r0.([]interface{}); !ok {
			return "", errors.Errorf("bad response: bad type of 'Returns' [%#v]", r0)
		} else if len(r1) != 2 {
			return "", errors.Errorf("bad response: bad length of 'Returns' [%#v]", r0)
		} else if r2, ok := r1[0].(map[string]interface{}); !ok {
			return "", errors.Errorf("bad response: bad type of first value [%#v]", r1)
		} else if r3, ok := r2["balance"]; !ok {
			return "", errors.Errorf("bad response: absent balance field [%#v]", r2)
		} else if balance, ok := r3.(string); !ok {
			return "", errors.Errorf("bad response: balance field expected to be a string [%#v]", r3)
		} else {
			return balance, nil
		}
	}

	r0, ok = result["Error"]
	if ok && r0 != nil {
		return "", errors.Errorf("%T: %#v", r0, r0)
	}

	panic("unreachable")
}

func BenchmarkSimple(b *testing.B) {
	ctx := context.Background()
	cfg := DefaultVMConfig()

	s, err := NewVirtualServer(b, ctx, cfg).WithGenesis().PrepareAndStart()
	require.NoError(b, err)
	defer s.Stop(ctx)

	var (
		iterations = b.N
		helper     = ServerHelper{s}
		syncAssert = assert.New(&testutils.SyncT{TB: b})
		sema       = semaphore.NewWeighted(100)
	)

	b.ResetTimer()

	var wg sync.WaitGroup
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func() {
			defer wg.Done()

			ctx, _ := inslogger.WithTraceField(ctx, uuid.New().String())

			err := sema.Acquire(ctx, 1)
			if err == nil {
				defer sema.Release(1)
			} else {
				panic(fmt.Sprintf("unexpected: %s", err.Error()))
			}

			_, err = helper.createUser(ctx)
			syncAssert.NoError(err, "failed to create user")
		}()
	}
	wg.Wait()
}
