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
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/helloworld"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/migrationshard"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/pkshard"
)

// RootDomain is smart contract representing entrance point to system.
type RootDomain struct {
	foundation.BaseContract
	MigrationAddressShards [insolar.GenesisAmountMigrationAddressShards]insolar.Reference
	PublicKeyShards        [insolar.GenesisAmountPublicKeyShards]insolar.Reference
	NodeDomain             insolar.Reference
}

// GetMemberByPublicKey gets member reference by public key.
// ins:immutable
func (rd RootDomain) GetMemberByPublicKey(publicKey string) (*insolar.Reference, error) {
	trimmedPublicKey := foundation.TrimPublicKey(publicKey)
	i := foundation.GetShardIndex(trimmedPublicKey, insolar.GenesisAmountPublicKeyShards)
	if i >= len(rd.PublicKeyShards) {
		return nil, fmt.Errorf("incorrect shard index")
	}
	s := pkshard.GetObject(rd.PublicKeyShards[i])
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
// ins:immutable
func (rd RootDomain) GetMemberByMigrationAddress(migrationAddress string) (*insolar.Reference, error) {
	trimmedMigrationAddress := foundation.TrimAddress(migrationAddress)
	i := foundation.GetShardIndex(trimmedMigrationAddress, insolar.GenesisAmountMigrationAddressShards)
	if i >= len(rd.MigrationAddressShards) {
		return nil, fmt.Errorf("incorrect shard index")
	}
	s := migrationshard.GetObject(rd.MigrationAddressShards[i])
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

// GetNodeDomainRef returns reference of NodeDomain instance
// ins:immutable
func (rd RootDomain) GetNodeDomainRef() (insolar.Reference, error) {
	return rd.NodeDomain, nil
}

var INSATTR_Info_API = true

// AddMigrationAddresses adds migration addresses to list.
// ins:immutable
func (rd *RootDomain) AddMigrationAddresses(migrationAddresses []string) error {
	newMA := [insolar.GenesisAmountMigrationAddressShards][]string{}
	for _, ma := range migrationAddresses {
		trimmedMigrationAddress := foundation.TrimAddress(ma)
		i := foundation.GetShardIndex(trimmedMigrationAddress, insolar.GenesisAmountMigrationAddressShards)
		if i >= len(newMA) {
			return fmt.Errorf("incorrect migration shard index")
		}
		newMA[i] = append(newMA[i], trimmedMigrationAddress)
	}

	for i, ma := range newMA {
		if len(ma) == 0 {
			continue
		}
		s := migrationshard.GetObject(rd.MigrationAddressShards[i])
		err := s.AddFreeMigrationAddresses(ma)
		if err != nil {
			return errors.New("failed to add migration addresses to shard")
		}
	}

	return nil
}

// AddMigrationAddress adds migration address to list.
// ins:immutable
func (rd *RootDomain) AddMigrationAddress(migrationAddress string) error {
	trimmedMigrationAddress := foundation.TrimAddress(migrationAddress)
	i := foundation.GetShardIndex(trimmedMigrationAddress, insolar.GenesisAmountMigrationAddressShards)
	if i >= len(rd.MigrationAddressShards) {
		return fmt.Errorf("incorrect migration shard index")
	}
	s := migrationshard.GetObject(rd.MigrationAddressShards[i])
	err := s.AddFreeMigrationAddresses([]string{trimmedMigrationAddress})
	if err != nil {
		return errors.New("failed to add migration address to shard")
	}

	return nil
}

// ins:immutable
func (rd *RootDomain) GetFreeMigrationAddress(publicKey string) (string, error) {
	trimmedPublicKey := foundation.TrimPublicKey(publicKey)
	shardIndex := foundation.GetShardIndex(trimmedPublicKey, insolar.GenesisAmountPublicKeyShards)
	if shardIndex >= len(rd.MigrationAddressShards) {
		return "", fmt.Errorf("incorrect migration address shard index")
	}

	for i := shardIndex; i < len(rd.MigrationAddressShards); i++ {
		mas := migrationshard.GetObject(rd.MigrationAddressShards[i])
		ma, err := mas.GetFreeMigrationAddress()

		if err == nil {
			return ma, nil
		}

		if err != nil {
			if !strings.Contains(err.Error(), "no more migration address left") {
				return "", errors.Wrap(err, "failed to set reference in migration address shard")
			}
		}
	}

	for i := 0; i < shardIndex; i++ {
		mas := migrationshard.GetObject(rd.MigrationAddressShards[i])
		ma, err := mas.GetFreeMigrationAddress()

		if err == nil {
			return ma, nil
		}

		if err != nil {
			if !strings.Contains(err.Error(), "no more migration address left") {
				return "", errors.Wrap(err, "failed to set reference in migration address shard")
			}
		}
	}

	return "", errors.New("no more migration addresses left in any shard")
}

// AddNewMemberToMaps adds new member to PublicKeyMap and MigrationAddressMap.
// ins:immutable
func (rd *RootDomain) AddNewMemberToMaps(publicKey string, migrationAddress string, memberRef insolar.Reference) error {
	trimmedPublicKey := foundation.TrimPublicKey(publicKey)
	shardIndex := foundation.GetShardIndex(trimmedPublicKey, insolar.GenesisAmountPublicKeyShards)
	if shardIndex >= len(rd.PublicKeyShards) {
		return fmt.Errorf("incorrect public key shard index")
	}
	pks := pkshard.GetObject(rd.PublicKeyShards[shardIndex])
	err := pks.SetRef(trimmedPublicKey, memberRef.String())
	if err != nil {
		return errors.Wrap(err, "failed to set reference in public key shard")
	}

	trimmedMigrationAddress := foundation.TrimAddress(migrationAddress)
	shardIndex = foundation.GetShardIndex(trimmedMigrationAddress, insolar.GenesisAmountPublicKeyShards)
	if shardIndex >= len(rd.MigrationAddressShards) {
		return fmt.Errorf("incorrect migration address shard index")
	}
	mas := migrationshard.GetObject(rd.MigrationAddressShards[shardIndex])
	err = mas.SetRef(migrationAddress, memberRef.String())
	if err != nil {
		return errors.Wrap(err, "failed to set reference in migration address shard")
	}

	return nil
}

// AddNewMemberToPublicKeyMap adds new member to PublicKeyMap.
// ins:immutable
func (rd *RootDomain) AddNewMemberToPublicKeyMap(publicKey string, memberRef insolar.Reference) error {
	trimmedPublicKey := foundation.TrimPublicKey(publicKey)
	i := foundation.GetShardIndex(trimmedPublicKey, insolar.GenesisAmountPublicKeyShards)
	if i >= len(rd.PublicKeyShards) {
		return fmt.Errorf("incorrect public key shard index")
	}
	s := pkshard.GetObject(rd.PublicKeyShards[i])
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
