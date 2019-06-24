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
	"fmt"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/application/contract/member/signer"
	"github.com/insolar/insolar/insolar"
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

func (hw *HelloWorld) Call(rootDomain insolar.Reference, method string, params []byte, seed []byte, sign []byte) (interface{}, error) {
	var name string
	switch method {
	case "Greet":
		if err := signer.UnmarshalParams(params, &name); err != nil {
			return nil, fmt.Errorf("[ registerNodeCall ] Can't unmarshal params: %s", err.Error())
		}
		return hw.Greet(name)
	case "Count":
		return hw.Count()
	case "Errored":
		return hw.Errored()
	case "CreateChild":
		return hw.CreateChild()
	case "CountChild":
		return hw.CountChild()
	default:
		return nil, errors.New("Unknown method")
	}
}

// New returns a new empty contract
func New() (*HelloWorld, error) {
	return &HelloWorld{
		Greeted: 0,
	}, nil
}
