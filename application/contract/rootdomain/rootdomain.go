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
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

// RootDomain is smart contract representing entrance point to system
type RootDomain struct {
	foundation.BaseContract
	RootMember        insolar.Reference
	OracleMembers     map[string]insolar.Reference
	MDAdminMember     insolar.Reference
	MDWallet          insolar.Reference
	BurnAddressMap    map[string]insolar.Reference
	PublicKeyMap      map[string]insolar.Reference
	FreeBurnAddresses []string
	NodeDomain        insolar.Reference
}

var INSATTR_CreateMember_API = true

// NewRootDomain creates new RootDomain
func NewRootDomain() (*RootDomain, error) {
	return &RootDomain{}, nil
}

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

// GetNodeDomainRef returns reference of NodeDomain instance
func (rd *RootDomain) GetNodeDomainRef() (insolar.Reference, error) {
	return rd.NodeDomain, nil
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

func (rd *RootDomain) AddBurnAddresses(burnAddresses []string) error {
	rd.FreeBurnAddresses = append(rd.FreeBurnAddresses, burnAddresses...)

	return nil
}

func (rd *RootDomain) AddBurnAddress(burnAddress string) error {
	rd.FreeBurnAddresses = append(rd.FreeBurnAddresses, burnAddress)

	return nil
}

func (rd *RootDomain) GetBurnAddress() (string, error) {
	if len(rd.FreeBurnAddresses) == 0 {
		return "", fmt.Errorf("[ GetBurnAddress ] No more burn address left")
	}

	return rd.FreeBurnAddresses[0], nil
}

func (rd *RootDomain) AddNewMemberToMaps(publicKey string, burnAddress string, memberRef insolar.Reference) error {
	rd.PublicKeyMap[publicKey] = memberRef
	rd.BurnAddressMap[burnAddress] = memberRef

	return nil
}

func (rd *RootDomain) GetMemberByBurnAddress(burnAddress string) (insolar.Reference, error) {
	return rd.BurnAddressMap[burnAddress], nil
}
