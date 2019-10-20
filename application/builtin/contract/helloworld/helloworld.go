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

package helloworld

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/application/builtin/contract/member/signer"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"

	hwProxy "github.com/insolar/insolar/application/builtin/proxy/helloworld"
)

// HelloWorld contract
type HelloWorld struct {
	foundation.BaseContract
	Greeted int
}

var INSATTR_Greet_API = true

type Text struct {
	SomeText string `json:"someText"`
}

type HwMessage struct {
	Message Text `json:"message"`
}

func (hw *HelloWorld) ReturnObj() (interface{}, error) {
	return hwProxy.HwMessage{Message: hwProxy.Text{SomeText: "Hello world"}}, nil
}

// Greet greats the caller
func (hw *HelloWorld) Greet(name string) (interface{}, error) {
	hw.Greeted++
	return fmt.Sprintf("Hello %s' world", name), nil
}

func (hw *HelloWorld) Count() (interface{}, error) {
	return hw.Greeted, nil
}

func (hw *HelloWorld) Errored() (interface{}, error) {
	return nil, errors.New("TestError")
}

// Get number pulse from foundation.
func (hw *HelloWorld) PulseNumber() (insolar.PulseNumber, error) {
	return foundation.GetPulseNumber()
}

func (hw *HelloWorld) CreateChild() (interface{}, error) {
	hwHolder := hwProxy.New()
	chw, err := hwHolder.AsChild(hw.GetReference())
	if err != nil {
		return nil, errors.Wrap(err, "[ HelloWorld.CreateChild ] Can't save as child")
	}
	return chw.GetReference().String(), nil
}

type Request struct {
	JsonRpc  string `json:"jsonrpc"`
	Id       int    `json:"id"`
	Method   string `json:"method"`
	Params   Params `json:"params"`
	LogLevel string `json:"logLevel,omitempty"`
}

type Params struct {
	Seed       string      `json:"seed"`
	CallSite   string      `json:"callSite"`
	CallParams interface{} `json:"callParams"`
	Reference  string      `json:"reference"`
	PublicKey  string      `json:"publicKey"`
}

func (hw *HelloWorld) Call(signedRequest []byte) (interface{}, error) {
	var signature string
	var pulseTimeStamp int64
	var rawRequest []byte

	err := signer.UnmarshalParams(signedRequest, &rawRequest, &signature, &pulseTimeStamp)
	if err != nil {
		return nil, fmt.Errorf(" Failed to decode: %s", err.Error())
	}

	request := Request{}
	err = json.Unmarshal(rawRequest, &request)
	if err != nil {
		return nil, fmt.Errorf(" Failed to unmarshal: %s", err.Error())
	}

	switch request.Params.CallSite {
	case "Greet":
		callParams, err := foundation.NewStableMapFromInterface(request.Params.CallParams)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse CallParams")
		}
		name, ok := callParams["name"]
		if !ok {
			return hw.Greet("Anonymous")
		}
		return hw.Greet(name)
	case "Count":
		return hw.Count()
	case "Errored":
		return hw.Errored()
	case "CreateChild":
		return hw.CreateChild()
	case "ReturnObj":
		return hw.ReturnObj()
	case "PulseNumber":
		return hw.PulseNumber()
	default:
		return nil, errors.New("unknown method " + request.Params.CallSite)
	}
}

// New returns a new empty contract
func New() (*HelloWorld, error) {
	return &HelloWorld{
		Greeted: 0,
	}, nil
}
