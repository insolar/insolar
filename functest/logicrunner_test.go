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
	result, err := callMethod(t, objectRef, "Get")
	require.Empty(t, err)
	require.Equal(t, float64(0), result)

	result, err = callMethod(t, objectRef, "Inc")
	require.Empty(t, err)
	require.Equal(t, float64(1), result)

	result, err = callMethod(t, objectRef, "Get")
	require.Empty(t, err)
	require.Equal(t, float64(1), result)

	result, err = callMethod(t, objectRef, "Dec")
	require.Empty(t, err)
	require.Equal(t, float64(0), result)

	result, err = callMethod(t, objectRef, "Get")
	require.Empty(t, err)
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

	resp, err := callMethod(t, objectRef, "Hello", "ins")
	require.Empty(t, err)
	require.Equal(t, "Hi, ins! Two said: Hello you too, ins. 1 times!", resp)

	for i := 2; i <= 5; i++ {
		resp, err = callMethod(t, objectRef, "Again", "ins")
		require.Empty(t, err)
		assert.Equal(t, fmt.Sprintf("Hi, ins! Two said: Hello you too, ins. %d times!", i), resp)
	}

	// TODO разобраться с кастом рефов
	// TODO может возвращать не реф а скасченный к стрингу реф? Reference.ToString() ?
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

	resp, err := callMethod(t, obj, "Hello", "ins")
	require.Empty(t, err)
	require.Equal(t, "Hi, ins! Two said: Hello you too, ins. 644 times!", resp)

	resp, err = callMethod(t, obj, "HelloFromDelegate", "ins")
	require.Empty(t, err)
	require.Equal(t, "Hello you too, ins. 1288 times!", resp)
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

	_, err := callMethod(t, obj, "Hello")
	require.Empty(t, err)

	resp, err := callMethod(t, obj, "Value")
	require.Empty(t, err)
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

	resp, err := callMethod(t, obj, "Hello")
	require.Empty(t, err)
	require.Equal(t, prototype.String(), resp)
}

func TestDeactivationError(t *testing.T) {
	var contractOneCode = `
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

	obj := callConstructor(t, uploadContractOnce(t, "deactivation", contractOneCode))

	_, err := callMethod(t, obj, "Kill")
	require.Empty(t, err)
}

// TODO вернуться позже или забить или сделать отдельный тест не через общую ручку а напрямую поймать ошибку
//func TestPanicError(t *testing.T) {
//	var contractOneCode = `
//package main
//
//import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
//import "errors"
//
//type One struct {
//	foundation.BaseContract
//}
//
//func (r *One) Panic() error {
//	return errors.New("test")
//}
//func (r *One) NotPanic() error {
//	return nil
//}
//`
//	obj := callConstructor(t, uploadContractOnce(t, "panic", contractOneCode))
//
//	_, err := callMethod(t, obj, "Panic") // need to check error
//	require.Equal(t, errors.New("test"), err.S)
//
//	_, err = callMethod(t, obj, "NotPanic") // no error
//	require.Empty(t, err)
//}

func TestGetChildrenError(t *testing.T) {
	goContract := `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	child "github.com/insolar/insolar/application/proxy/get_children_child"
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

	uploadContractOnce(t, "get_children_child", goChild)
	obj := callConstructor(t, uploadContractOnce(t, "get_children_one", goContract))

	resp, err := callMethod(t, obj, "SumChildsByIterator")
	require.Empty(t, err, "empty children")
	require.Equal(t, float64(0), resp)

	resp, err = callMethod(t, obj, "NewChilds", 10)
	require.Empty(t, err, "add children")
	require.Equal(t, float64(45), resp)

	resp, err = callMethod(t, obj, "SumChildsByIterator")
	require.Empty(t, err, "sum real children")
	require.Equal(t, float64(45), resp)
}


// TODO return 400, expects 200
//func TestErrorInterfaceError(t *testing.T) {
//	var contractOneCode = `
//package main
//
//import (
//	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
//	two "github.com/insolar/insolar/application/proxy/error_interface_two"
//)
//
//type One struct {
//	foundation.BaseContract
//}
//
//func (r *One) AnError() error {
//	holder := two.New()
//	friend, err := holder.AsChild(r.GetReference())
//	if err != nil {
//		return err
//	}
//
//	return friend.AnError()
//}
//
//func (r *One) NoError() error {
//	holder := two.New()
//	friend, err := holder.AsChild(r.GetReference())
//	if err != nil {
//		return err
//	}
//
//	return friend.NoError()
//}
//`
//
//	var contractTwoCode = `
//package main
//
//import (
//	"errors"
//
//	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
//)
//
//type Two struct {
//	foundation.BaseContract
//}
//func New() (*Two, error) {
//	return &Two{}, nil
//}
//func (r *Two) AnError() error {
//	return errors.New("an error")
//}
//func (r *Two) NoError() error {
//	return nil
//}
//`
//	uploadContractOnce(t, "error_interface_two", contractTwoCode)
//	obj := callConstructor(t, uploadContractOnce(t, "error_interface_one", contractOneCode))
//
//	resp, err := callMethod(t, obj, "AnError")
//	require.Equal(t, &foundation.Error{S: "an error"}, err)
//
//	resp, err = callMethod(t, obj, "NoError")
//	require.Nil(t, resp)
//}

func TestNilResultError(t *testing.T) {
	var contractOneCode = `
package main

import (
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	two "github.com/insolar/insolar/application/proxy/nil_result_two"
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

	uploadContractOnce(t, "nil_result_two", contractTwoCode)
	obj := callConstructor(t, uploadContractOnce(t, "nil_result_one", contractOneCode))

	resp, err := callMethod(t, obj, "Hello")
	require.Empty(t, err)
	require.Nil(t, resp)
}
