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

	"github.com/insolar/insolar/application/contract/member/signer"
	"github.com/insolar/insolar/application/proxy/rootdomain"
	"github.com/insolar/insolar/application/proxy/wallet"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
	"github.com/ugorji/go/codec"
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

func UnmarshalParams(data []byte, to ...interface{}) error {
	ch := new(codec.CborHandle)
	return codec.NewDecoderBytes(data, ch).Decode(&to)
}

func (m *Member) Call(ref core.RecordRef, method string, params []byte, seed []byte, sign []byte) (interface{}, *foundation.Error) {

	args, err := core.MarshalArgs(
		m.GetReference(),
		method,
		params,
		seed)
	if err != nil {
		return nil, &foundation.Error{S: err.Error()}
	}
	verified, err := ecdsa.Verify(args, sign, m.GetPublicKey())
	if err != nil {
		return nil, &foundation.Error{S: err.Error()}
	}
	if !verified {
		return nil, &foundation.Error{S: "Incorrect signature"}
	}

	rootDomain := rootdomain.GetObject(ref)

	switch method {
	case "CreateMember":
		var name string
		var key string
		err := UnmarshalParams(params, &name, &key)
		if err != nil {
			return nil, &foundation.Error{S: err.Error()}
		}
		fmt.Println("PARAMS", name, key)
		member := rootDomain.CreateMember(name, key)
		return member, nil
	case "GetMyBalance":
		balance := wallet.GetImplementationFrom(m.GetReference()).GetTotalBalance()
		return balance, nil
	case "GetBalance":
		var member string
		err := UnmarshalParams(params, &member)
		if err != nil {
			return nil, &foundation.Error{S: err.Error()}
		}
		balance := wallet.GetImplementationFrom(core.NewRefFromBase58(member)).GetTotalBalance()
		return balance, nil
	case "Transfer":
		var amount float64
		var toStr string
		err := UnmarshalParams(params, &amount, &toStr)
		if err != nil {
			return nil, &foundation.Error{S: err.Error()}
		}
		to := core.NewRefFromBase58(toStr)
		wallet.GetImplementationFrom(m.GetReference()).Transfer(uint(amount), &to)
		return nil, nil
	}
	return nil, &foundation.Error{S: "Unknown method"}
}

func (m *Member) AuthorizedCall(ref core.RecordRef, delegate core.RecordRef, method string, params []byte, seed []byte, sign []byte) ([]byte, *foundation.Error) {
	serialized, err := signer.Serialize(ref[:], delegate[:], method, params, seed)
	if err != nil {
		return nil, &foundation.Error{S: err.Error()}
	}
	verified, err := ecdsa.Verify(serialized, sign, m.PublicKey)
	if err != nil {
		return nil, &foundation.Error{S: err.Error()}
	}
	if !verified {
		return nil, &foundation.Error{S: "Incorrect signature"}
	}

	var contract core.RecordRef
	if !delegate.Equal(core.RecordRef{}) {
		contract, err = foundation.GetImplementationFor(ref, delegate)
		if err != nil {
			return nil, &foundation.Error{S: err.Error()}
		}
	} else {
		contract = ref
	}
	ret, err := proxyctx.Current.RouteCall(contract, true, method, params)
	if err != nil {
		return nil, &foundation.Error{S: err.Error()}
	}
	return ret, nil
}
