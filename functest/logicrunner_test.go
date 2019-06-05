///
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
///

// +build functest

package functest

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSingleContractError(t *testing.T) {
	var contractCode = `
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
	objectRef := callConstructor(t, uploadContractOnce(t, "test", contractCode))

	// be careful - jsonUnmarshal convert json numbers to float64
	result := callMethod(t, objectRef, "Get")
	require.Equal(t, float64(0), result)

	result = callMethod(t, objectRef, "Inc")
	require.Equal(t, float64(1), result)

	result = callMethod(t, objectRef, "Get")
	require.Equal(t, float64(1), result)

	result = callMethod(t, objectRef, "Dec")
	require.Equal(t, float64(0), result)

	result = callMethod(t, objectRef, "Get")
	require.Equal(t, float64(0), result)
}

func TestContractCallingContractError(t *testing.T) {
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

	uploadContractOnce(t, "two", contractTwoCode)
	objectRef := callConstructor(t, uploadContractOnce(t, "one", contractOneCode))

	resp := callMethod(t, objectRef, "Hello", "ins")
	require.Equal(t, "Hi, ins! Two said: Hello you too, ins. 1 times!", resp)

	for i := 2; i <= 5; i++ {
		resp = callMethod(t, objectRef, "Again", "ins")
		assert.Equal(t, fmt.Sprintf("Hi, ins! Two said: Hello you too, ins. %d times!", i), resp)
	}

	// TODO разобраться с кастом рефов
	//resp = callMethod(t, objectRef, "GetFriend")
	//
	//// this is base64 instead of base58
	//fmt.Println(">>> resp ", resp)
	//fmt.Println(">>> CBORUnMarshal ", goplugintestutils.CBORUnMarshal(t, []byte(resp.(string))))
	//fmt.Println(">>> resp.([]byte) ", []byte(resp.(string)))
	//two := insolar.Reference{}.FromSlice([]byte(resp.(string)))

	//r0 := resp.([]uint8)
	//var two insolar.Reference
	//for i := 0; i < 64; i++ {
	//	two[i] = r0[i]
	//}

	//for i := 6; i <= 9; i++ {
	//	resp = callMethod(t, &two, "Hello", "Insolar")
	//	assert.Equal(t, fmt.Sprintf("Hello you too, Insolar. %d times!", i), resp)
	//}
	//
	//resp = callMethod(t, &two, "TestPayload")
	//res := resp.(map[interface{}]interface{})["Str"]
	//assert.Equal(t,"HiHere", res)
}

func TestInjectingDelegateError(t *testing.T) {
	var contractOneCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import two "github.com/insolar/insolar/application/proxy/injection_delegate_two"

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

	uploadContractOnce(t, "injection_delegate_two", contractTwoCode)
	obj := callConstructor(t, uploadContractOnce(t, "injection_delegate_one", contractOneCode))

	resp := callMethod(t, obj, "Hello", "ins")
	require.Equal(t,"Hi, ins! Two said: Hello you too, ins. 644 times!", resp)

	resp = callMethod(t, obj, "HelloFromDelegate", "ins")
	require.Equal(t,"Hello you too, ins. 1288 times!", resp)
}

func TestBasicNotificationCallError(t *testing.T) {
	var contractOneCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import two "github.com/insolar/insolar/application/proxy/basic_notification_call_two"

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
	uploadContractOnce(t, "basic_notification_call_two", contractTwoCode)
	obj := callConstructor(t, uploadContractOnce(t, "basic_notification_call_one", contractOneCode))

	// no error
	callMethod(t, obj, "Hello")

	resp := callMethod(t, obj, "Value")
	require.Equal(t, float64(644), resp)
}

func TestContextPassingError(t *testing.T) {
	var contractOneCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type One struct {
	foundation.BaseContract
}

func (r *One) Hello() (string, error) {
	return r.GetPrototype().String(), nil
}
`
	prototype := uploadContractOnce(t, "context_passing", contractOneCode)
	obj := callConstructor(t, prototype)

	resp := callMethod(t, obj, "Hello")
	require.Equal(t, prototype.String(), resp)
}
