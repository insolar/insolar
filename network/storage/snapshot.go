// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package storage

import (
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/node"
	"github.com/pkg/errors"
)

//go:generate minimock -i github.com/insolar/insolar/network/storage.SnapshotStorage -o ../../testutils/network -s _mock.go -g

// SnapshotStorage provides methods for accessing Snapshot.
type SnapshotStorage interface {
	ForPulseNumber(insolar.PulseNumber) (*node.Snapshot, error)
	Append(pulse insolar.PulseNumber, snapshot *node.Snapshot) error
}

// newSnapshotStorage constructor creates PulseStorage
func newSnapshotStorage() *snapshotStorage { // nolint
	return &snapshotStorage{}
}

type snapshotStorage struct { // nolint
	DB   DB `inject:""`
	lock sync.RWMutex
}

func (s *snapshotStorage) Append(pulse insolar.PulseNumber, snapshot *node.Snapshot) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	buff, err := snapshot.Encode()
	if err != nil {
		return errors.Wrap(err, "[snapshotStorage] Failed to append snapshot")
	}
	return s.DB.Set(pulseKey(pulse), buff)
}

func (s *snapshotStorage) ForPulseNumber(pulse insolar.PulseNumber) (*node.Snapshot, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	buf, err := s.DB.Get(pulseKey(pulse))
	if err != nil {
		return nil, errors.Wrap(err, "[snapshotStorage] Failed to get snapshot from DB")
	}
	result := &node.Snapshot{}
	err = result.Decode(buf)
	if err != nil {
		return nil, errors.Wrap(err, "[snapshotStorage] Failed to decode snapshot")
	}
	return result, nil
}
