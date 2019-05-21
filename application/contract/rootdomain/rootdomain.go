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

package rootdomain

import (
	"encoding/json"
	"fmt"
	"github.com/insolar/insolar/application/proxy/member"
	"github.com/insolar/insolar/application/proxy/wallet"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

// RootDomain is smart contract representing entrance point to system
type RootDomain struct {
	foundation.BaseContract
	RootMember    insolar.Reference
	OracleMembers map[string]insolar.Reference
	MDAdminMember insolar.Reference
	MDWallet      insolar.Reference
	NodeDomain    insolar.Reference
}

var INSATTR_CreateMember_API = true

func (rd *RootDomain) GetMDAdminMemberRef() (*insolar.Reference, error) {
	return &rd.MDAdminMember, nil
}

func (rd *RootDomain) GetMDWalletRef() (*insolar.Reference, error) {
	return &rd.MDWallet, nil
}

func (rd *RootDomain) GetOracleMembers() (map[string]insolar.Reference, error) {
	return rd.OracleMembers, nil
}

func (rd *RootDomain) GetRootMemberRef() (*insolar.Reference, error) {
	return &rd.RootMember, nil
}

// CreateMember processes create member request
func (rd *RootDomain) CreateMember(name string, key string) (string, error) {
	memberHolder := member.New(name, key)
	m, err := memberHolder.AsChild(rd.GetReference())
	if err != nil {
		return "", fmt.Errorf("[ CreateMember ] Can't save as child: %s", err.Error())
	}

	wHolder := wallet.New(1000 * 1000 * 1000)
	_, err = wHolder.AsDelegate(m.GetReference())
	if err != nil {
		return "", fmt.Errorf("[ CreateMember ] Can't save as delegate: %s", err.Error())
	}

	return m.GetReference().String(), nil
}

var INSATTR_Info_API = true

// Info returns information about basic objects
func (rd *RootDomain) Info() (interface{}, error) {
	oracleMembersOut := map[string]string{}
	for name, ref := range rd.OracleMembers {
		oracleMembersOut[name] = ref.String()
	}

	res := map[string]interface{}{
		"root_member":     rd.RootMember.String(),
		"oracle_members":  oracleMembersOut,
		"md_admin_member": rd.MDAdminMember.String(),
		"node_domain":     rd.NodeDomain.String(),
	}
	resJSON, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("[ Info ] Can't marshal res: %s", err.Error())
	}
	return resJSON, nil
}

// GetNodeDomainRef returns reference of NodeDomain instance
func (rd *RootDomain) GetNodeDomainRef() (insolar.Reference, error) {
	return rd.NodeDomain, nil
}

// NewRootDomain creates new RootDomain
func NewRootDomain() (*RootDomain, error) {
	return &RootDomain{}, nil
}

// DumpAllUsers processes dump all users request
func (rd *RootDomain) DumpAllUsers() (*proxyctx.ChildrenTypedIterator, error) {
	return rd.NewChildrenTypedIterator(member.GetPrototype())
}
