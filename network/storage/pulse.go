//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package storage

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
)

//go:generate minimock -i github.com/insolar/insolar/network/storage.PulseAccessor -o ../../testutils/network -s _mock.go -g

// PulseAccessor provides methods for accessing pulses.
type PulseAccessor interface {
	ForPulseNumber(context.Context, insolar.PulseNumber) (insolar.Pulse, error)
	Latest(ctx context.Context) (insolar.Pulse, error)
}

//go:generate minimock -i github.com/insolar/insolar/network/storage.PulseAppender -o ../../testutils/network -s _mock.go -g

// PulseAppender provides method for appending pulses to storage.
type PulseAppender interface {
	Append(ctx context.Context, pulse insolar.Pulse) error
}

//go:generate minimock -i github.com/insolar/insolar/network/storage.PulseCalculator -o ../../testutils/network -s _mock.go -g

// PulseCalculator performs calculations for pulses.
type PulseCalculator interface {
	Forwards(ctx context.Context, pn insolar.PulseNumber, steps int) (insolar.Pulse, error)
	Backwards(ctx context.Context, pn insolar.PulseNumber, steps int) (insolar.Pulse, error)
}

//go:generate minimock -i github.com/insolar/insolar/network/storage.PulseRangeHasher -o ../../testutils/network -s _mock.go -g

// PulseRangeHasher provides methods for hashing and validate Pulse chain
type PulseRangeHasher interface {
	GetRangeHash(insolar.PulseRange) ([]byte, error)
	ValidateRangeHash(insolar.PulseRange, []byte) (bool, error)
}

// NewPulseStorage constructor creates PulseStorage
func NewPulseStorage() *PulseStorage {
	return &PulseStorage{}
}

type PulseStorage struct {
	DB   DB `inject:""`
	lock sync.RWMutex
}

func (p *PulseStorage) GetRangeHash(insolar.PulseRange) ([]byte, error) {
	panic("implement me")
}

func (p *PulseStorage) ValidateRangeHash(insolar.PulseRange, []byte) (bool, error) {
	panic("implement me")
}

// Forwards calculates steps pulses forwards from provided Pulse. If calculated Pulse does not exist, ErrNotFound will
// be returned.
func (p *PulseStorage) Forwards(ctx context.Context, pn insolar.PulseNumber, steps int) (pulse insolar.Pulse, err error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	node, err := p.get(pn)
	if err != nil {
		return
	}

	iterator := node
	for i := 0; i < steps; i++ {
		if iterator.Next == nil {
			err = insolar.ErrNotFound
			return
		}
		iterator, err = p.get(*iterator.Next)
		if err != nil {
			return
		}
	}

	return iterator.Pulse, nil
}

// Backwards calculates steps pulses backwards from provided Pulse. If calculated Pulse does not exist, ErrNotFound will
// be returned.
func (p *PulseStorage) Backwards(ctx context.Context, pn insolar.PulseNumber, steps int) (pulse insolar.Pulse, err error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	node, err := p.get(pn)
	if err != nil {
		return
	}

	iterator := node
	for i := 0; i < steps; i++ {
		if iterator.Prev == nil {
			err = insolar.ErrNotFound
			return
		}
		iterator, err = p.get(*iterator.Prev)
		if err != nil {
			return
		}
	}

	return iterator.Pulse, nil
}

// Append appends provided Pulse to current storage. Pulse number should be greater than currently saved for preserving
// Pulse consistency. If provided Pulse does not meet the requirements, ErrBadPulse will be returned.
func (p *PulseStorage) Append(ctx context.Context, pulse insolar.Pulse) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	var insertWithHead = func(head insolar.PulseNumber) error {
		oldHead, err := p.get(head)
		if err != nil {
			return err
		}
		oldHead.Next = &pulse.PulseNumber

		// Set new Pulse.
		err = p.set(pulse.PulseNumber, dbNode{
			Prev:  &oldHead.Pulse.PulseNumber,
			Pulse: pulse,
		})
		if err != nil {
			return err
		}
		// Set old updated head.
		err = p.set(oldHead.Pulse.PulseNumber, oldHead)
		if err != nil {
			return err
		}
		// Set head meta record.
		return p.setHead(pulse.PulseNumber)
	}
	var insertWithoutHead = func() error {
		// Set new Pulse.
		err := p.set(pulse.PulseNumber, dbNode{
			Pulse: pulse,
		})
		if err != nil {
			return err
		}
		// Set head meta record.
		return p.setHead(pulse.PulseNumber)
	}

	head, err := p.head()
	if err == ErrNotFound {
		return insertWithoutHead()
	}

	if pulse.PulseNumber <= head {
		return ErrBadPulse
	}
	return insertWithHead(head)
}

func (p *PulseStorage) ForPulseNumber(ctx context.Context, pn insolar.PulseNumber) (pulse insolar.Pulse, err error) {
	nd, err := p.get(pn)
	if err != nil {
		return
	}
	return nd.Pulse, nil
}

func (p *PulseStorage) Latest(ctx context.Context) (insolar.Pulse, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	head, err := p.head()
	if err != nil {
		return insolar.Pulse{}, err
	}
	nd, err := p.get(head)
	if err != nil {
		return insolar.Pulse{}, err
	}
	return nd.Pulse, nil
}

type metaKey byte

func (k metaKey) Scope() Scope {
	return ScopePulse
}

func (k metaKey) ID() []byte {
	return []byte{prefixMeta, byte(k)}
}

type dbNode struct {
	Pulse      insolar.Pulse
	Prev, Next *insolar.PulseNumber
}

var (
	prefixPulse byte = 1
	prefixMeta  byte = 2
)

var (
	keyHead metaKey = 1
)

func (p *PulseStorage) get(pn insolar.PulseNumber) (nd dbNode, err error) {
	buf, err := p.DB.Get(pulseKey(pn))
	if err == ErrNotFound {
		err = ErrNotFound
		return
	}
	if err != nil {
		return
	}
	nd = deserialize(buf)
	return
}

func (p *PulseStorage) set(pn insolar.PulseNumber, nd dbNode) error {
	return p.DB.Set(pulseKey(pn), serialize(nd))
}

func (p *PulseStorage) head() (pn insolar.PulseNumber, err error) {
	buf, err := p.DB.Get(keyHead)
	if err == ErrNotFound {
		err = ErrNotFound
		return
	}
	if err != nil {
		return
	}
	pn = insolar.NewPulseNumber(buf)
	return
}

func (p *PulseStorage) setHead(pn insolar.PulseNumber) error {
	return p.DB.Set(keyHead, pn.Bytes())
}

func serialize(nd dbNode) []byte {
	return insolar.MustSerialize(nd)
}

func deserialize(buf []byte) (nd dbNode) {
	insolar.MustDeserialize(buf, &nd)
	return nd
}
