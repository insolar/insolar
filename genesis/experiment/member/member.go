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

package member

import (
	"fmt"
	"reflect"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/genesis/experiment/nodedomain/utils"
	"github.com/insolar/insolar/genesis/proxy/member"
	"github.com/insolar/insolar/genesis/proxy/rootdomain"
	"github.com/insolar/insolar/genesis/proxy/wallet"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

type Member struct {
	foundation.BaseContract
	Name      string
	PublicKey string
}

func (m *Member) GetName() string {
	return m.Name
}
func (m *Member) GetPublicKey() string {
	return m.PublicKey
}

func New(name string, key string) *Member {
	return &Member{
		Name:      name,
		PublicKey: key,
	}
}

/*func (m *Member) Check(val []byte, sign []byte) bool {
	res, err := utils.Verify(val, sign, m.PublicKey)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return res
}*/

type Msg struct {
	Ref    string
	Method string
	Params []interface{}
	Seed   []byte
}

func (m *Member) AuthorizedCall(ref string, method string, params []interface{}, seed []byte, sign []byte) ([]interface{}, *foundation.Error) {
	args := Msg{ref, method, params, seed}
	var serialized []byte
	err := proxyctx.Current.Serialize(args, &serialized)
	if err != nil {
		return nil, &foundation.Error{err.Error()}
	}
	verified, err := utils.Verify(serialized, sign, m.PublicKey)
	if err != nil {
		return nil, &foundation.Error{err.Error()}
	}
	if !verified {
		return nil, &foundation.Error{"Incorrect signature"}
	}

	switch method {
	case "CreateMember":
		domain := rootdomain.GetObject(core.NewRefFromBase58(ref))
		name, ok := params[0].(string)
		if !ok {
			return nil, &foundation.Error{"First parameter must be string"}
		}
		key, ok := params[1].(string)
		if !ok {
			return nil, &foundation.Error{"Second parameter must be string"}
		}
		return []interface{}{domain.CreateMember(name, key)}, nil
	case "GetName":
		membr := member.GetObject(core.NewRefFromBase58(ref))
		return []interface{}{membr.GetName()}, nil
	case "GetPublicKey":
		membr := member.GetObject(core.NewRefFromBase58(ref))
		return []interface{}{membr.GetPublicKey()}, nil
	case "GetBalance":
		wallet := wallet.GetImplementationFrom(core.NewRefFromBase58(ref))
		return []interface{}{wallet.GetTotalBalance()}, nil
	case "SendMoney":
		fmt.Println("SENDING", params, reflect.TypeOf(params[0]), reflect.TypeOf(params[1]))
		wallet := wallet.GetImplementationFrom(core.NewRefFromBase58(ref))
		amount, ok := params[0].(uint64)
		if !ok {
			fmt.Println(params[0])
			return nil, &foundation.Error{"First parameter must be uint"}
		}
		to, ok := params[1].(string)
		if !ok {
			return nil, &foundation.Error{"Second parameter must be string"}
		}
		v := core.NewRefFromBase58(to)
		fmt.Println("HERE", wallet, amount, to)
		wallet.Transfer(uint(amount), &v)
		return nil, nil
	case "DumpAllUsers":
		domain := rootdomain.GetObject(core.NewRefFromBase58(ref))
		return []interface{}{domain.DumpAllUsers()}, nil
	}
	return nil, &foundation.Error{"Unknown method"}
}
