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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/helloworld"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/shard"
)

// RootDomain is smart contract representing entrance point to system.
type RootDomain struct {
	foundation.BaseContract
	RootMember             insolar.Reference
	MigrationDaemonMembers [insolar.GenesisAmountActiveMigrationDaemonMembers]insolar.Reference
	MigrationAdminMember   insolar.Reference
	MigrationWallet        insolar.Reference
	CostCenter             insolar.Reference
	FeeWallet              insolar.Reference
	MigrationAddressShards [insolar.GenesisAmountMigrationAddressShards]insolar.Reference
	PublicKeyShards        [insolar.GenesisAmountPublicKeyShards]insolar.Reference
	FreeBurnAddresses      []string
	NodeDomain             insolar.Reference
}

// GetMigrationAdminMemberRef gets migration admin member reference.
func (rd RootDomain) GetCostCenterRef() (insolar.Reference, error) {
	return rd.MigrationAdminMember, nil
}

// GetFeeWalletRef gets fee wallet reference.
func (rd RootDomain) GetFeeWalletRef() (insolar.Reference, error) {
	return rd.FeeWallet, nil
}

// GetMigrationWalletRef gets migration wallet reference.
func (rd RootDomain) GetMigrationWalletRef() (insolar.Reference, error) {
	return rd.MigrationWallet, nil
}

// GetMigrationAdminMember gets migration admin member reference.
func (rd RootDomain) GetMigrationAdminMember() (insolar.Reference, error) {
	return rd.MigrationAdminMember, nil
}

// GetActiveMigrationDaemonMembers gets migration daemon members references.
func (rd RootDomain) GetActiveMigrationDaemonMembers() ([3]insolar.Reference, error) {
	return rd.MigrationDaemonMembers, nil
}

// GetRootMemberRef gets root member reference.
func (rd RootDomain) GetRootMemberRef() (insolar.Reference, error) {
	return rd.RootMember, nil
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
func (rd RootDomain) GetMemberByPublicKey(publicKey string) (*insolar.Reference, error) {
	trimmedPublicKey := trimPublicKey(publicKey)
	i := foundation.GetShardIndex(trimmedPublicKey, insolar.GenesisAmountPublicKeyShards)
	if i >= len(rd.PublicKeyShards) {
		return nil, fmt.Errorf("incorect shard index")
	}
	s := shard.GetObject(rd.PublicKeyShards[i])
	refStr, err := s.GetRef(trimmedPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get reference in shard")
	}
	ref, err := insolar.NewReferenceFromBase58(refStr)
	if err != nil {
		return nil, errors.Wrap(err, "bad member reference for this public key")
	}

	return ref, nil
}

// GetMemberByMigrationAddress gets member reference by burn address.
func (rd RootDomain) GetMemberByMigrationAddress(migrationAddress string) (*insolar.Reference, error) {
	trimmedMigrationAddress := trimMigrationAddress(migrationAddress)
	i := foundation.GetShardIndex(trimmedMigrationAddress, insolar.GenesisAmountMigrationAddressShards)
	if i >= len(rd.MigrationAddressShards) {
		return nil, fmt.Errorf("incorect shard index")
	}
	s := shard.GetObject(rd.MigrationAddressShards[i])
	refStr, err := s.GetRef(trimmedMigrationAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get reference in shard")
	}
	ref, err := insolar.NewReferenceFromBase58(refStr)
	if err != nil {
		return nil, errors.Wrap(err, "bad member reference for this migration address")
	}

	return ref, nil
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

// AddNewMemberToMaps adds new member to PublicKeyMap and MigrationAddressMap.
func (rd *RootDomain) AddNewMemberToMaps(publicKey string, migrationAddress string, memberRef insolar.Reference) error {
	trimmedPublicKey := trimPublicKey(publicKey)
	i := foundation.GetShardIndex(trimmedPublicKey, insolar.GenesisAmountPublicKeyShards)
	if i >= len(rd.PublicKeyShards) {
		return fmt.Errorf("incorect public key shard index")
	}
	s := shard.GetObject(rd.PublicKeyShards[i])
	err := s.SetRef(trimmedPublicKey, memberRef.String())
	if err != nil {
		return errors.Wrap(err, "failed to set reference in public key shard")
	}

	trimmedMigrationAddress := trimMigrationAddress(migrationAddress)
	i = foundation.GetShardIndex(trimmedMigrationAddress, insolar.GenesisAmountMigrationAddressShards)
	if i >= len(rd.MigrationAddressShards) {
		return fmt.Errorf("incorect migration shard index")
	}
	s = shard.GetObject(rd.MigrationAddressShards[i])
	err = s.SetRef(trimmedMigrationAddress, memberRef.String())
	if err != nil {
		return errors.Wrap(err, "failed to set reference in migration address shard")
	}

	return nil
}

// AddNewMemberToPublicKeyMap adds new member to PublicKeyMap.
func (rd *RootDomain) AddNewMemberToPublicKeyMap(publicKey string, memberRef insolar.Reference) error {
	trimmedPublicKey := trimPublicKey(publicKey)
	i := foundation.GetShardIndex(trimmedPublicKey, insolar.GenesisAmountPublicKeyShards)
	if i >= len(rd.PublicKeyShards) {
		return fmt.Errorf("incorect public key shard index")
	}
	s := shard.GetObject(rd.PublicKeyShards[i])
	err := s.SetRef(trimmedPublicKey, memberRef.String())
	if err != nil {
		return errors.Wrap(err, "failed to set reference in public key shard")
	}

	return nil
}

func (rd *RootDomain) CreateHelloWorld() (string, error) {
	helloWorldHolder := helloworld.New()
	m, err := helloWorldHolder.AsChild(rd.GetReference())
	if err != nil {
		return "", fmt.Errorf("failed to save as child: %s", err.Error())
	}

	return m.GetReference().String(), nil
}

func trimPublicKey(publicKey string) string {
	return trimMigrationAddress(between(publicKey, "KEY-----", "-----END"))
}

func trimMigrationAddress(burnAddress string) string {
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
