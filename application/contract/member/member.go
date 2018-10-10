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
	"github.com/insolar/insolar/application/contract/member/signer"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptohelpers/ecdsa"
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

func (m *Member) AuthorizedCall(ref string, delegate string, method string, params []byte, seed []byte, sign []byte) ([]byte, *foundation.Error) {
	serialized, err := signer.Serialize(ref, delegate, method, params, seed)
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
	if delegate != "" {
		contract, err = foundation.GetImplementationFor(core.NewRefFromBase58(ref), core.NewRefFromBase58(delegate))
		if err != nil {
			return nil, &foundation.Error{S: err.Error()}
		}
	} else {
		contract = core.NewRefFromBase58(ref)
	}
	ret, err := proxyctx.Current.RouteCall(contract, true, method, params)
	if err != nil {
		return nil, &foundation.Error{S: err.Error()}
	}
	return ret, nil
}
