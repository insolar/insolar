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
	"strings"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/helloworld"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

// RootDomain is smart contract representing entrance point to system.
type RootDomain struct {
	foundation.BaseContract
	RootMember             insolar.Reference
	MigrationDaemonMembers []insolar.Reference
	MigrationAdminMember   insolar.Reference
	MigrationWallet        insolar.Reference
	CostCenter             insolar.Reference
	CommissionWallet       insolar.Reference
	BurnAddressMap         map[string]insolar.Reference
	PublicKeyMap           map[string]insolar.Reference
	FreeBurnAddresses      []string
	NodeDomain             insolar.Reference
}

// GetMigrationAdminMemberRef gets migration admin member reference.
func (rd RootDomain) GetMigrationAdminMemberRef() (*insolar.Reference, error) {
	return &rd.MigrationAdminMember, nil
}

// GetMigrationWalletRef gets migration wallet reference.
func (rd RootDomain) GetMigrationWalletRef() (*insolar.Reference, error) {
	return &rd.MigrationWallet, nil
}

// GetMigrationDaemonMembers gets migration daemon members references.
func (rd RootDomain) GetMigrationDaemonMembers() ([]insolar.Reference, error) {
	return rd.MigrationDaemonMembers, nil
}

// GetRootMemberRef gets root member reference.
func (rd RootDomain) GetRootMemberRef() (*insolar.Reference, error) {
	return &rd.RootMember, nil
}

// GetBurnAddress pulls out burn address from list.
func (rd *RootDomain) GetBurnAddress() (string, error) {
	if len(rd.FreeBurnAddresses) == 0 {
		return "", fmt.Errorf("no more burn addresses left")
	}

	result := rd.FreeBurnAddresses[0]
	rd.FreeBurnAddresses = rd.FreeBurnAddresses[1:]

	return result, nil
}

// GetMemberByPublicKey gets member reference by public key.
func (rd RootDomain) GetMemberByPublicKey(publicKey string) (insolar.Reference, error) {
	return rd.PublicKeyMap[trimPublicKey(publicKey)], nil
}

// GetMemberByBurnAddress gets member reference by burn address.
func (rd RootDomain) GetMemberByBurnAddress(burnAddress string) (insolar.Reference, error) {
	return rd.BurnAddressMap[trimBurnAddress(burnAddress)], nil
}

// GetCostCenter gets cost center reference.
func (rd RootDomain) GetCostCenter() (insolar.Reference, error) {
	return rd.CostCenter, nil
}

// GetNodeDomainRef returns reference of NodeDomain instance
func (rd RootDomain) GetNodeDomainRef() (insolar.Reference, error) {
	return rd.NodeDomain, nil
}

var INSATTR_Info_API = true

// Info returns information about basic objects
func (rd RootDomain) Info() (interface{}, error) {
	migrationDaemonsMembersOut := []string{}
	for _, ref := range rd.MigrationDaemonMembers {
		migrationDaemonsMembersOut = append(migrationDaemonsMembersOut, ref.String())
	}

	res := map[string]interface{}{
		"rootDomain":             rd.GetReference().String(),
		"rootMember":             rd.RootMember.String(),
		"migrationDaemonMembers": migrationDaemonsMembersOut,
		"migrationAdminMember":   rd.MigrationAdminMember.String(),
		"nodeDomain":             rd.NodeDomain.String(),
	}
	resJSON, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %s", err.Error())
	}
	return resJSON, nil
}

// AddBurnAddresses adds burn addresses to list.
func (rd *RootDomain) AddBurnAddresses(burnAddresses []string) error {
	rd.FreeBurnAddresses = append(rd.FreeBurnAddresses, burnAddresses...)

	return nil
}

// AddBurnAddress adds burn address to list.
func (rd *RootDomain) AddBurnAddress(burnAddress string) error {
	rd.FreeBurnAddresses = append(rd.FreeBurnAddresses, burnAddress)

	return nil
}

// AddNewMemberToMaps adds new member to PublicKeyMap and BurnAddressMap.
func (rd *RootDomain) AddNewMemberToMaps(publicKey string, burnAddress string, memberRef insolar.Reference) error {
	if _, ok := rd.PublicKeyMap[trimPublicKey(publicKey)]; ok {
		return fmt.Errorf("member for this publicKey already exist")
	}
	rd.PublicKeyMap[trimPublicKey(publicKey)] = memberRef

	if _, ok := rd.PublicKeyMap[trimPublicKey(burnAddress)]; ok {
		return fmt.Errorf("member for this burnAddress already exist")
	}
	rd.BurnAddressMap[trimBurnAddress(burnAddress)] = memberRef

	return nil
}

// AddNewMemberToPublicKeyMap adds new member to PublicKeyMap
func (rd *RootDomain) AddNewMemberToPublicKeyMap(publicKey string, memberRef insolar.Reference) error {
	if _, ok := rd.PublicKeyMap[trimPublicKey(publicKey)]; ok {
		return fmt.Errorf("member for this publicKey already exist")
	}
	rd.PublicKeyMap[trimPublicKey(publicKey)] = memberRef

	return nil
}

func (rd *RootDomain) CreateHelloWorld() (map[string]interface{}, error) {
	helloWorldHolder := helloworld.New()
	m, err := helloWorldHolder.AsChild(rd.GetReference())
	if err != nil {
		return nil, fmt.Errorf("failed to save as child: %s", err.Error())
	}

	return map[string]interface{}{"reference": m.GetReference().String()}, nil
}

func trimPublicKey(publicKey string) string {
	return trimBurnAddress(between(publicKey, "KEY-----", "-----END"))
}

func trimBurnAddress(burnAddress string) string {
	return strings.ToLower(strings.Join(strings.Split(strings.TrimSpace(burnAddress), "\n"), ""))
}

func between(value string, a string, b string) string {
	// Get substring between two strings.
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirst := pos + len(a)
	if posFirst >= posLast {
		return ""
	}
	return value[posFirst:posLast]
}
