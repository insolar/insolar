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

	"github.com/insolar/insolar/logicrunner/builtin/contract/member/signer"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"

	hwProxy "github.com/insolar/insolar/logicrunner/builtin/proxy/helloworld"
)

// HelloWorld contract
type HelloWorld struct {
	foundation.BaseContract
	Greeted int
}

var INSATTR_Greet_API = true

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

func (hw *HelloWorld) CreateChild() (interface{}, error) {
	hwHolder := hwProxy.New()
	chw, err := hwHolder.AsChild(hw.GetReference())
	if err != nil {
		return nil, errors.Wrap(err, "[ HelloWorld.CreateChild ] Can't save as child")
	}
	return chw.GetReference().String(), nil
}

func (hw *HelloWorld) CountChild() (interface{}, error) {
	count := 0

	iterator, err := hw.NewChildrenTypedIterator(hwProxy.GetPrototype())
	if err != nil {
		return nil, fmt.Errorf("[ CountChild ] Can't get children: %s", err.Error())
	}

	for iterator.HasNext() {
		cref, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("[ CountChild ] Can't get next child: %s", err.Error())
		}

		m := hwProxy.GetObject(cref)

		childCountI, err := m.Count()
		if err != nil {
			return nil, fmt.Errorf("[ CountChild ] Can't get count of child: %s", err.Error())
		}

		childCount, ok := childCountI.(uint64)
		if !ok {
			return nil, fmt.Errorf("[ CountChild ] Bad childCount format, expected int got %T", childCountI)
		}

		count = count + int(childCount)
	}

	return count, nil
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
	PublicKey  string      `json:"memberPublicKey"`
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
		return hw.Greet(request.Params.CallParams.(map[string]interface{})["name"].(string))
	case "Count":
		return hw.Count()
	case "Errored":
		return hw.Errored()
	case "CreateChild":
		return hw.CreateChild()
	case "CountChild":
		return hw.CountChild()
	default:
		return nil, errors.New("Unknown method " + request.Params.CallSite)
	}
}

// New returns a new empty contract
func New() (*HelloWorld, error) {
	return &HelloWorld{
		Greeted: 0,
	}, nil
}
