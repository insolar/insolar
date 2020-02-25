// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package storage

import (
	"sync"

	"github.com/insolar/insolar/insolar"
)

//go:generate minimock -i github.com/insolar/insolar/network/storage.CloudHashStorage -o ../../testutils/network -s _mock.go -g

// CloudHashStorage provides methods for accessing CloudHash.
type CloudHashStorage interface {
	ForPulseNumber(pulse insolar.PulseNumber) ([]byte, error)
	Append(pulse insolar.PulseNumber, cloudHash []byte) error
}

// newCloudHashStorage constructor creates cloudHashStorage
func newCloudHashStorage() *cloudHashStorage { // nolint
	return &cloudHashStorage{}
}

// NewMemoryCloudHashStorage constructor creates cloudHashStorage
func NewMemoryCloudHashStorage() *MemoryCloudHashStorage {
	return &MemoryCloudHashStorage{
		entries: make(map[insolar.PulseNumber][]byte),
	}
}

type cloudHashStorage struct { // nolint
	DB   DB `inject:""`
	lock sync.RWMutex
}

func (c *cloudHashStorage) ForPulseNumber(pulse insolar.PulseNumber) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	result, err := c.DB.Get(pulseKey(pulse))
	if err != nil {
		return nil, err
	}
	return result, err
}

func (c *cloudHashStorage) Append(pulse insolar.PulseNumber, cloudHash []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	return c.DB.Set(pulseKey(pulse), cloudHash)
}

type MemoryCloudHashStorage struct {
	lock    sync.RWMutex
	entries map[insolar.PulseNumber][]byte
}

func (m *MemoryCloudHashStorage) ForPulseNumber(pulse insolar.PulseNumber) ([]byte, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if s, ok := m.entries[pulse]; ok {
		return s, nil
	}
	return nil, ErrNotFound
}

func (m *MemoryCloudHashStorage) Append(pulse insolar.PulseNumber, cloudHash []byte) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.entries[pulse] = cloudHash
	return nil
}
