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

package mashard

import (
	"github.com/pkg/errors"

	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

// MAShard - shard contract for migration addresses.
type MAShard struct {
	foundation.BaseContract
	Map    foundation.StableMap
	FreeMA []string
}

// New creates new member.
func New() (*MAShard, error) {
	return &MAShard{
		Map:    map[string]string{},
		FreeMA: []string{},
	}, nil
}

// GetMigrationAddressesAmount gets amount of free migration addresses
func (s MAShard) GetMigrationAddressesAmount(migrationAddresses []string) (int, error) {
	return len(s.FreeMA), nil
}

// AddFreeMigrationAddresses add new addresses to the array of free migration addresses
func (s *MAShard) AddFreeMigrationAddresses(migrationAddresses []string) error {
	s.FreeMA = append(s.FreeMA, migrationAddresses...)
	return nil
}

// GetFreeMigrationAddress gets free migration address from list
func (s *MAShard) GetFreeMigrationAddress() (string, error) {
	if len(s.FreeMA) <= 0 {
		return "", errors.New("no more migration address left")
	}
	ma := s.FreeMA[0]
	s.FreeMA = s.FreeMA[1:]

	return ma, nil
}

// GetRef gets ref by key.
func (s MAShard) GetRef(key string) (string, error) {
	if ref, ok := s.Map[key]; !ok {
		return "", errors.New("failed to find reference by key")
	} else {
		return ref, nil
	}
}

// SetRef sets ref with migration address as a key.
func (s *MAShard) SetRef(ma string, ref string) error {
	if _, ok := s.Map[ma]; ok {
		return errors.New("can't set reference because this key already exists")
	}
	s.Map[ma] = ref
	return nil
}
