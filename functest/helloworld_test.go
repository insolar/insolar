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

// +build functest

package functest

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/testutils/launchnet"
)

type HelloWorldInstance struct {
	Ref *insolar.Reference
}

func NewHelloWorld(ctx context.Context) (*HelloWorldInstance, error) {
	seed, err := requester.GetSeed(launchnet.TestRPCUrl)
	if err != nil {
		return nil, err
	}

	rootCfg, err := requester.CreateUserConfig(launchnet.Root.Ref, launchnet.Root.PrivKey, launchnet.Root.PubKey)
	res, err := requester.SendWithSeed(ctx, launchnet.TestRPCUrl, rootCfg, &requester.Params{
		CallSite:   "CreateHelloWorld",
		CallParams: make(foundation.StableMap),
		PublicKey:  rootCfg.PublicKey},
		seed)
	if err != nil {
		return nil, err
	}

	var result requester.ContractResponse
	err = json.Unmarshal(res, &result)
	if err != nil {
		return nil, err
	} else if result.Error != nil {
		return nil, errors.Errorf("[ NewHelloWorld ] Failed to execute: %s", result.Error.Message)
	}

	rv, ok := result.Result.CallResult.(string)
	if !ok {
		return nil, errors.Errorf("[ NewHelloWorld ] Failed to decode: expected string, got %T", result.Result)
	}

	i := HelloWorldInstance{}
	i.Ref, err = insolar.NewReferenceFromBase58(rv)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (i *HelloWorldInstance) Greet(ctx context.Context, name string) (string, error) {
	seed, err := requester.GetSeed(launchnet.TestRPCUrl)
	if err != nil {
		return "", err
	}

	rootCfg, err := requester.CreateUserConfig(i.Ref.String(), launchnet.Root.PrivKey, launchnet.Root.PubKey)
	callParams := make(foundation.StableMap)
	callParams["name"] = name
	res, err := requester.SendWithSeed(ctx, launchnet.TestRPCUrl, rootCfg, &requester.Params{
		CallSite:   "Greet",
		CallParams: callParams,
		PublicKey:  rootCfg.PublicKey},
		seed)
	if err != nil {
		return "", err
	}

	var result requester.ContractResponse
	err = json.Unmarshal(res, &result)
	if err != nil {
		return "", err
	} else if result.Error != nil {
		return "", errors.Errorf("[ Greet ] Failed to execute: %s", result.Error.Message)
	}

	rv, ok := result.Result.CallResult.(string)
	if !ok {
		return "", errors.Errorf("[ Greet ] Failed to decode: expected string, got %T", result.Result)
	}
	return rv, nil
}

func (i *HelloWorldInstance) Count(ctx context.Context) (int, error) {
	seed, err := requester.GetSeed(launchnet.TestRPCUrl)
	if err != nil {
		return 0, err
	}

	rootCfg, err := requester.CreateUserConfig(i.Ref.String(), launchnet.Root.PrivKey, launchnet.Root.PubKey)
	res, err := requester.SendWithSeed(ctx, launchnet.TestRPCUrl, rootCfg, &requester.Params{
		CallSite:   "Count",
		CallParams: make(foundation.StableMap),
		PublicKey:  rootCfg.PublicKey},
		seed)
	if err != nil {
		return 0, err
	}

	var result requester.ContractResponse
	err = json.Unmarshal(res, &result)
	if err != nil {
		return 0, err
	} else if result.Error != nil {
		return 0, errors.Errorf("[ Count ] Failed to execute: %s", result.Error.Message)
	}

	rv, ok := result.Result.CallResult.(float64)
	if !ok {
		return 0, errors.Errorf("[ Count ] Failed to decode: expected float64, got %T", result.Result)
	}
	return int(rv), nil
}

func (i *HelloWorldInstance) PulseNumber(t *testing.T, ctx context.Context) (int, error) {
	member := &launchnet.User{i.Ref.String(), launchnet.Root.PrivKey, launchnet.Root.PubKey}
	result, err := signedRequest(t, launchnet.TestRPCUrl, member, "PulseNumber", nil)
	if err != nil {
		return 0, err
	}
	rv, ok := result.(float64)
	if !ok {
		return 0, errors.Errorf("failed to decode: expected float64, got %T", result)
	}
	return int(rv), nil
}

func (i *HelloWorldInstance) CreateChild(ctx context.Context) (*HelloWorldInstance, error) {
	seed, err := requester.GetSeed(launchnet.TestRPCUrl)
	if err != nil {
		return nil, err
	}

	rootCfg, err := requester.CreateUserConfig(i.Ref.String(), launchnet.Root.PrivKey, launchnet.Root.PubKey)
	res, err := requester.SendWithSeed(ctx, launchnet.TestRPCUrl, rootCfg, &requester.Params{
		CallSite:   "CreateChild",
		CallParams: make(foundation.StableMap),
		PublicKey:  rootCfg.PublicKey},
		seed)
	if err != nil {
		return nil, err
	}

	var result requester.ContractResponse
	err = json.Unmarshal(res, &result)
	if err != nil {
		return nil, err
	} else if result.Error != nil {
		return nil, errors.Errorf("[ CreateChild ] Failed to execute: %s", result.Error.Message)
	}

	rv, ok := result.Result.CallResult.(string)
	if !ok {
		return nil, errors.Errorf("[ CreateChild ] Failed to decode: expected string, got %T", result.Result)
	}

	child := HelloWorldInstance{}
	child.Ref, err = insolar.NewReferenceFromBase58(rv)
	if err != nil {
		return nil, err
	}

	return &child, nil
}

func (i *HelloWorldInstance) ReturnObj(ctx context.Context) (map[string]interface{}, error) {
	seed, err := requester.GetSeed(launchnet.TestRPCUrl)
	if err != nil {
		return nil, err
	}

	rootCfg, err := requester.CreateUserConfig(i.Ref.String(), launchnet.Root.PrivKey, launchnet.Root.PubKey)
	res, err := requester.SendWithSeed(ctx, launchnet.TestRPCUrl, rootCfg, &requester.Params{
		CallSite:   "ReturnObj",
		CallParams: make(foundation.StableMap),
		PublicKey:  rootCfg.PublicKey},
		seed)
	if err != nil {
		return nil, err
	}

	var result requester.ContractResponse
	err = json.Unmarshal(res, &result)
	if err != nil {
		return nil, err
	} else if result.Error != nil {
		return nil, errors.Errorf("[ ReturnObj ] Failed to execute: %s", result.Error.Message)
	}

	rv, ok := result.Result.CallResult.(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("[ ReturnObj ] Failed to decode result: expected map[string]interface{}, got %T", result.Result.CallResult)
	}
	return rv, nil
}

func TestCallHelloWorld(t *testing.T) {
	a, r := assert.New(t), require.New(t)
	ctx := context.TODO()

	hw, err := NewHelloWorld(ctx)
	r.NoError(err, "Unexpected error")
	a.NotEmpty(hw.Ref, "Ref doesn't exists")

	for i := 0; i < 100; i++ {
		val, err := hw.Greet(ctx, "Simon")
		r.NoError(err, "Unexpected error was thrown on Greet")
		a.Contains(val, "Simon'", "Returned message doesn't contains Simon")
	}

	count, err := hw.Count(ctx)
	r.NoError(err)
	// tip: right now deduplication is not presented in our system, so number of created
	//      requests should be less or equal to result count of registered requests
	a.LessOrEqual(100, count)
}

func TestCallPulseNumber(t *testing.T) {
	a, r := assert.New(t), require.New(t)
	ctx := context.TODO()

	hw, err := NewHelloWorld(ctx)
	r.NoError(err, "Unexpected error")
	a.NotEmpty(hw.Ref, "Ref doesn't exists")
	pulseNum, err := hw.PulseNumber(t, ctx)
	r.NoError(err)

	r.True(pulseNum > 0)
}

func TestCallHelloWorldReturnObj(t *testing.T) {
	a, r := assert.New(t), require.New(t)
	ctx := context.TODO()

	hw, err := NewHelloWorld(ctx)
	r.NoError(err, "Unexpected error")
	a.NotEmpty(hw.Ref, "Ref doesn't exists")

	val, err := hw.ReturnObj(ctx)
	r.NoError(err)
	r.Equal(val["message"].(map[string]interface{})["someText"], "Hello world")
}
