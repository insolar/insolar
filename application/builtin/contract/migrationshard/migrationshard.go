// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package migrationshard

import (
	"github.com/pkg/errors"

	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

// MigrationShard - shard contract for migration addresses.
type MigrationShard struct {
	foundation.BaseContract
	Map                    foundation.StableMap
	FreeMigrationAddresses []string
}

// New creates new member.
func New(migrationAddresses []string) (*MigrationShard, error) {
	return &MigrationShard{
		Map:                    make(foundation.StableMap),
		FreeMigrationAddresses: migrationAddresses,
	}, nil
}

// GetMigrationAddressesAmount gets amount of free migration addresses
// ins:immutable
func (s *MigrationShard) GetMigrationAddressesAmount() (int, error) {
	return len(s.FreeMigrationAddresses), nil
}

// AddFreeMigrationAddresses add new addresses to the array of free migration addresses
func (s *MigrationShard) AddFreeMigrationAddresses(migrationAddresses []string) error {
	s.FreeMigrationAddresses = append(s.FreeMigrationAddresses, migrationAddresses...)
	return nil
}

// GetFreeMigrationAddress gets free migration address from list
func (s *MigrationShard) GetFreeMigrationAddress() (string, error) {
	if len(s.FreeMigrationAddresses) <= 0 {
		return "", errors.New("no more migration address left")
	}
	ma := s.FreeMigrationAddresses[0]
	s.FreeMigrationAddresses = s.FreeMigrationAddresses[1:]

	return ma, nil
}

// GetRef gets ref by key.
// ins:immutable
func (s *MigrationShard) GetRef(key string) (string, error) {
	if ref, ok := s.Map[key]; !ok {
		return "", errors.New("failed to find reference by key")
	} else {
		return ref, nil
	}
}

// SetRef sets ref with migration address as a key.
func (s *MigrationShard) SetRef(ma string, ref string) error {
	if _, ok := s.Map[ma]; ok {
		return errors.New("can't set reference because this key already exists")
	}
	s.Map[ma] = ref
	return nil
}
