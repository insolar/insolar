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
	"math/rand"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/insolar"
)

type HelloWorldInstance struct {
	Ref *insolar.Reference
}

func NewHelloWorld(ctx context.Context) (*HelloWorldInstance, error) {
	seed, err := requester.GetSeed(TestAPIURL)
	if err != nil {
		return nil, err
	}

	rootCfg, err := requester.CreateUserConfig(root.ref, root.privKey, root.pubKey)
	res, err := requester.SendWithSeed(ctx, TestCallUrl, rootCfg, &requester.Request{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "api.call",
		Params:  requester.Params{CallSite: "CreateHelloWorld", CallParams: map[string]interface{}{}, PublicKey: rootCfg.PublicKey},
	}, seed)
	if err != nil {
		return nil, err
	}

	var result requester.ContractAnswer
	err = json.Unmarshal(res, &result)
	if err != nil {
		return nil, err
	} else if result.Error != nil {
		return nil, errors.Errorf("[ NewHelloWorld ] Failed to execute: %s", result.Error.Message)
	}

	rv, ok := result.Result.ContractResult.(string)
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
	seed, err := requester.GetSeed(TestAPIURL)
	if err != nil {
		return "", err
	}

	rootCfg, err := requester.CreateUserConfig(i.Ref.String(), root.privKey, root.pubKey)
	res, err := requester.SendWithSeed(ctx, TestCallUrl, rootCfg, &requester.Request{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "api.call",
		Params:  requester.Params{CallSite: "Greet", CallParams: map[string]interface{}{"name": name}, PublicKey: rootCfg.PublicKey},
	}, seed)
	if err != nil {
		return "", err
	}

	var result requester.ContractAnswer
	err = json.Unmarshal(res, &result)
	if err != nil {
		return "", err
	} else if result.Error != nil {
		return "", errors.Errorf("[ Greet ] Failed to execute: %s", result.Error.Message)
	}

	rv, ok := result.Result.ContractResult.(string)
	if !ok {
		return "", errors.Errorf("[ Greet ] Failed to decode: expected string, got %T", result.Result)
	}
	return rv, nil
}

func (i *HelloWorldInstance) Count(ctx context.Context) (int, error) {
	seed, err := requester.GetSeed(TestAPIURL)
	if err != nil {
		return 0, err
	}

	rootCfg, err := requester.CreateUserConfig(i.Ref.String(), root.privKey, root.pubKey)
	res, err := requester.SendWithSeed(ctx, TestCallUrl, rootCfg, &requester.Request{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "api.call",
		Params:  requester.Params{CallSite: "Count", CallParams: map[string]interface{}{}, PublicKey: rootCfg.PublicKey},
	}, seed)
	if err != nil {
		return 0, err
	}

	var result requester.ContractAnswer
	err = json.Unmarshal(res, &result)
	if err != nil {
		return 0, err
	} else if result.Error != nil {
		return 0, errors.Errorf("[ Count ] Failed to execute: %s", result.Error.Message)
	}

	rv, ok := result.Result.ContractResult.(float64)
	if !ok {
		return 0, errors.Errorf("[ Count ] Failed to decode: expected float64, got %T", result.Result)
	}
	return int(rv), nil
}

func (i *HelloWorldInstance) CreateChild(ctx context.Context) (*HelloWorldInstance, error) {
	seed, err := requester.GetSeed(TestAPIURL)
	if err != nil {
		return nil, err
	}

	rootCfg, err := requester.CreateUserConfig(i.Ref.String(), root.privKey, root.pubKey)
	res, err := requester.SendWithSeed(ctx, TestCallUrl, rootCfg, &requester.Request{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "api.call",
		Params:  requester.Params{CallSite: "CreateChild", CallParams: map[string]interface{}{}, PublicKey: rootCfg.PublicKey},
	}, seed)
	if err != nil {
		return nil, err
	}

	var result requester.ContractAnswer
	err = json.Unmarshal(res, &result)
	if err != nil {
		return nil, err
	} else if result.Error != nil {
		return nil, errors.Errorf("[ CreateChild ] Failed to execute: %s", result.Error.Message)
	}

	rv, ok := result.Result.ContractResult.(string)
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

func (i *HelloWorldInstance) CountChild(ctx context.Context) (int, error) {
	seed, err := requester.GetSeed(TestAPIURL)
	if err != nil {
		return 0, err
	}

	rootCfg, err := requester.CreateUserConfig(i.Ref.String(), root.privKey, root.pubKey)
	res, err := requester.SendWithSeed(ctx, TestCallUrl, rootCfg, &requester.Request{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "api.call",
		Params:  requester.Params{CallSite: "CreateChild", CallParams: map[string]interface{}{}, PublicKey: rootCfg.PublicKey},
	}, seed)
	if err != nil {
		return 0, err
	}

	var result requester.ContractAnswer
	err = json.Unmarshal(res, &result)
	if err != nil {
		return 0, err
	} else if result.Error != nil {
		return 0, errors.Errorf("[ CountChild ] Failed to execute: %s", result.Error.Message)
	}

	rv, ok := result.Result.ContractResult.(float64)
	if !ok {
		return 0, errors.Errorf("[ CountChild ] Failed to decode result: expected float64, got %T", result.Result)
	}
	return int(rv), nil
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

func TestCallHelloWorldChild(t *testing.T) {
	t.Skip("Feature 'child of contract object' is not stable right now")

	a, r := assert.New(t), require.New(t)
	ctx := context.TODO()

	hw, err := NewHelloWorld(ctx)
	r.NoError(err, "Unexpected error")
	a.NotEmpty(hw.Ref, "Ref doesn't exists")

	var children []*HelloWorldInstance
	var childrenCntArray []int
	var childrenCnt int
	for i := 0; i < 10; i++ {
		hwt, err := hw.CreateChild(ctx)
		r.NoError(err)
		r.NotEmpty(hwt.Ref)
		children = append(children, hwt)

		cnt := rand.Int() % 13
		for i := 0; i < cnt; i++ {
			val, err := hwt.Greet(ctx, "Martha")
			r.NoError(err, "Unexpected error was thrown on Greet")
			a.Contains(val, "Martha'", "Returned message doesn't contains Martha")
		}
		childrenCntArray = append(childrenCntArray, cnt)
		childrenCnt = childrenCnt + cnt
	}

	for i := 0; i < 10; i++ {
		count, err := children[i].Count(ctx)
		r.NoError(err)
		a.Equal(childrenCntArray[i], count)
	}

	countOverall, err := hw.CountChild(ctx)
	r.NoError(err)
	a.Equal(countOverall, childrenCnt)
}
