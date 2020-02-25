// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package storage

import (
	"github.com/insolar/insolar/insolar"
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is returned when value was not found.
	ErrNotFound = errors.New("value not found")
	ErrBadPulse = errors.New("pulse should be bigger than latest")
)

type pulseKey insolar.PulseNumber

func (k pulseKey) Scope() Scope {
	return ScopePulse
}

func (k pulseKey) ID() []byte {
	return append([]byte{prefixPulse}, insolar.PulseNumber(k).Bytes()...)
}

// DB provides a simple key-value store interface for persisting data.
type DB interface {
	Get(key Key) (value []byte, err error)
	Set(key Key, value []byte) error
}

// Key represents a key for the key-value store. Scope is required to separate different DB clients and should be
// unique.
type Key interface {
	Scope() Scope
	ID() []byte
}

// Scope separates DB clients.
type Scope byte

// Bytes returns binary scope representation.
func (s Scope) Bytes() []byte {
	return []byte{byte(s)}
}

const (
	// ScopePulse is the scope for pulse storage.
	ScopePulse Scope = 1
)
