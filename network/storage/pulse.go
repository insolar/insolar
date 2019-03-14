/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package storage

import (
	"context"
	"github.com/insolar/insolar/core"
)

// PulseAccessor provides methods for accessing pulses.
type PulseAccessor interface {
	ForPulseNumber(context.Context, core.PulseNumber) (core.PulseNumber, error)
	Latest(ctx context.Context) (core.Pulse, error)
}

// PulseAppender provides method for appending pulses to storage.
type PulseAppender interface {
	Append(ctx context.Context, pulse core.Pulse) error
}

// PulseCalculator performs calculations for pulses.
type PulseCalculator interface {
	Forwards(ctx context.Context, pn core.PulseNumber, steps int) (core.Pulse, error)
	Backwards(ctx context.Context, pn core.PulseNumber, steps int) (core.Pulse, error)
}

// PulseRangeHasher provides methods for hashing and validate pulse chain
type PulseRangeHasher interface {
	GetRangeHash(core.PulseRange) ([]byte, error)
	ValidateRangeHash(core.PulseRange, []byte) (bool, error)
}

// PulseChainHasher provides methods for hashing and validate pulse chain
type PulseChainHasher interface {
	GetRangeHash(chain []core.PulseNumber) ([]byte, error)
	ValidateRangeHash(chain []core.PulseNumber, hash []byte) (bool, error)
}

// NewPulseStorage constructor creates PulseStorage
func NewPulseStorage() *PulseStorage {
	return &PulseStorage{}
}

type PulseStorage struct {
}

func (p *PulseStorage) Forwards(ctx context.Context, pn core.PulseNumber, steps int) (core.Pulse, error) {
	panic("implement me")
}

func (p *PulseStorage) Backwards(ctx context.Context, pn core.PulseNumber, steps int) (core.Pulse, error) {
	panic("implement me")
}

func (p *PulseStorage) Append(ctx context.Context, pulse core.Pulse) error {
	panic("implement me")
}

func (p *PulseStorage) ForPulseNumber(context.Context, core.PulseNumber) (core.PulseNumber, error) {
	panic("implement me")
}

func (p *PulseStorage) Latest(ctx context.Context) (core.Pulse, error) {
	panic("implement me")
}
