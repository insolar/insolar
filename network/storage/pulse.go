/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification,
 * are permitted (subject to the limitations in the disclaimer below) provided that
 * the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *  * Neither the name of Insolar Technologies nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
 * BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
 * CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING,
 * BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package storage

import (
	"bytes"
	"context"

	"github.com/insolar/insolar/core"
	"github.com/ugorji/go/codec"
	"sync"
)

//go:generate minimock -i github.com/insolar/insolar/network/storage.PulseAccessor -o ../../testutils/network -s _mock.go

// PulseAccessor provides methods for accessing pulses.
type PulseAccessor interface {
	ForPulseNumber(context.Context, core.PulseNumber) (core.Pulse, error)
	Latest(ctx context.Context) (core.Pulse, error)
}

//go:generate minimock -i github.com/insolar/insolar/network/storage.PulseAppender -o ../../testutils/network -s _mock.go

// PulseAppender provides method for appending pulses to storage.
type PulseAppender interface {
	Append(ctx context.Context, pulse core.Pulse) error
}

//go:generate minimock -i github.com/insolar/insolar/network/storage.PulseCalculator -o ../../testutils/network -s _mock.go

// PulseCalculator performs calculations for pulses.
type PulseCalculator interface {
	Forwards(ctx context.Context, pn core.PulseNumber, steps int) (core.Pulse, error)
	Backwards(ctx context.Context, pn core.PulseNumber, steps int) (core.Pulse, error)
}

//go:generate minimock -i github.com/insolar/insolar/network/storage.PulseRangeHasher -o ../../testutils/network -s _mock.go

// PulseRangeHasher provides methods for hashing and validate pulse chain
type PulseRangeHasher interface {
	GetRangeHash(core.PulseRange) ([]byte, error)
	ValidateRangeHash(core.PulseRange, []byte) (bool, error)
}

// NewPulseStorage constructor creates PulseStorage
func NewPulseStorage() *PulseStorage {
	return &PulseStorage{}
}

type PulseStorage struct {
	DB   DB `inject:""`
	lock sync.RWMutex
}

func (p *PulseStorage) GetRangeHash(core.PulseRange) ([]byte, error) {
	panic("implement me")
}

func (p *PulseStorage) ValidateRangeHash(core.PulseRange, []byte) (bool, error) {
	panic("implement me")
}

// Forwards calculates steps pulses forwards from provided pulse. If calculated pulse does not exist, ErrNotFound will
// be returned.
func (p *PulseStorage) Forwards(ctx context.Context, pn core.PulseNumber, steps int) (pulse core.Pulse, err error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	node, err := p.get(pn)
	if err != nil {
		return
	}

	iterator := node
	for i := 0; i < steps; i++ {
		if iterator.next == nil {
			err = core.ErrNotFound
			return
		}
		iterator, err = p.get(*iterator.next)
		if err != nil {
			return
		}
	}

	return iterator.pulse, nil
}

// Backwards calculates steps pulses backwards from provided pulse. If calculated pulse does not exist, ErrNotFound will
// be returned.
func (p *PulseStorage) Backwards(ctx context.Context, pn core.PulseNumber, steps int) (pulse core.Pulse, err error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	node, err := p.get(pn)
	if err != nil {
		return
	}

	iterator := node
	for i := 0; i < steps; i++ {
		if iterator.prev == nil {
			err = core.ErrNotFound
			return
		}
		iterator, err = p.get(*iterator.prev)
		if err != nil {
			return
		}
	}

	return iterator.pulse, nil
}

// Append appends provided pulse to current storage. Pulse number should be greater than currently saved for preserving
// pulse consistency. If provided pulse does not meet the requirements, ErrBadPulse will be returned.
func (p *PulseStorage) Append(ctx context.Context, pulse core.Pulse) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	var insertWithHead = func(head core.PulseNumber) error {
		oldHead, err := p.get(head)
		if err != nil {
			return err
		}
		oldHead.next = &pulse.PulseNumber

		// Set new pulse.
		err = p.set(pulse.PulseNumber, dbNode{
			prev:  &oldHead.pulse.PulseNumber,
			pulse: pulse,
		})
		if err != nil {
			return err
		}
		// Set old updated head.
		err = p.set(oldHead.pulse.PulseNumber, oldHead)
		if err != nil {
			return err
		}
		// Set head meta record.
		return p.setHead(pulse.PulseNumber)
	}
	var insertWithoutHead = func() error {
		// Set new Pulse.
		err := p.set(pulse.PulseNumber, dbNode{
			pulse: pulse,
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

func (p *PulseStorage) ForPulseNumber(ctx context.Context, pn core.PulseNumber) (pulse core.Pulse, err error) {
	nd, err := p.get(pn)
	if err != nil {
		return
	}
	return nd.pulse, nil
}

func (p *PulseStorage) Latest(ctx context.Context) (core.Pulse, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	head, err := p.head()
	if err != nil {
		return core.Pulse{}, err
	}
	nd, err := p.get(head)
	if err != nil {
		return core.Pulse{}, err
	}
	return nd.pulse, nil
}

type pulseKey core.PulseNumber

func (k pulseKey) Scope() Scope {
	return ScopePulse
}

func (k pulseKey) ID() []byte {
	return append([]byte{prefixPulse}, core.PulseNumber(k).Bytes()...)
}

type metaKey byte

func (k metaKey) Scope() Scope {
	return ScopePulse
}

func (k metaKey) ID() []byte {
	return []byte{prefixMeta, byte(k)}
}

type dbNode struct {
	pulse      core.Pulse
	prev, next *core.PulseNumber
}

var (
	prefixPulse byte = 1
	prefixMeta  byte = 2
)

var (
	keyHead metaKey = 1
)

func (p *PulseStorage) get(pn core.PulseNumber) (nd dbNode, err error) {
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

func (p *PulseStorage) set(pn core.PulseNumber, nd dbNode) error {
	return p.DB.Set(pulseKey(pn), serialize(nd))
}

func (p *PulseStorage) head() (pn core.PulseNumber, err error) {
	buf, err := p.DB.Get(keyHead)
	if err == ErrNotFound {
		err = ErrNotFound
		return
	}
	if err != nil {
		return
	}
	pn = core.NewPulseNumber(buf)
	return
}

func (p *PulseStorage) setHead(pn core.PulseNumber) error {
	return p.DB.Set(keyHead, pn.Bytes())
}

func serialize(nd dbNode) []byte {
	buff := bytes.NewBuffer(nil)
	enc := codec.NewEncoder(buff, &codec.CborHandle{})
	enc.MustEncode(nd)
	return buff.Bytes()
}

func deserialize(buf []byte) (nd dbNode) {
	dec := codec.NewDecoderBytes(buf, &codec.CborHandle{})
	dec.MustDecode(&nd)
	return nd
}
