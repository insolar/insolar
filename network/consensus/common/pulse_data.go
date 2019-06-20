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

package common

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const InvalidPulseEpoch uint32 = 0
const EphemeralPulseEpoch = InvalidPulseEpoch + 1

var _ PulseDataReader = &PulseData{}

type PulseData struct {
	PulseNumber PulseNumber
	PulseDataExt
}

type PulseDataExt struct {
	// ByteSize=44
	PulseEpoch     uint32
	PulseEntropy   Bits256
	NextPulseDelta uint16
	PrevPulseDelta uint16
	Timestamp      uint32
}

type PulseDataReader interface {
	GetPulseNumber() PulseNumber
	GetStartOfEpoch() PulseNumber
	// GetPulseEntropy()	[4]uint64
	GetNextPulseDelta() uint16
	GetPrevPulseDelta() uint16
	GetTimestamp() uint64
	IsExpectedPulse() bool
	IsFromEphemeral() bool
}

func NewFirstPulsarData(delta uint16) *PulseData {
	return newPulsarData(PulseNumberOfNow(), delta)
}

func NewPulsarData(pn PulseNumber, deltaNext uint16, deltaPrev uint16) *PulseData {
	r := newPulsarData(pn, deltaNext)
	r.PrevPulseDelta = deltaPrev
	return r
}

func NewFirstEphemeralData() *PulseData {
	return newEphemeralData(MinRegularPulseNumber)
}

func (r PulseData) String() string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprint(r.PulseNumber))

	ep := PulseNumberOfUint32(r.PulseEpoch)
	if ep != r.PulseNumber && ep != 0 {
		buf.WriteString(fmt.Sprintf("@%d", ep))
	}
	if r.NextPulseDelta == r.PrevPulseDelta {
		buf.WriteString(fmt.Sprintf(",±%d", r.NextPulseDelta))
	} else {
		if r.NextPulseDelta > 0 {
			buf.WriteString(fmt.Sprintf(",+%d", r.NextPulseDelta))
		}
		if r.PrevPulseDelta > 0 {
			buf.WriteString(fmt.Sprintf(",-%d", r.PrevPulseDelta))
		}
	}
	return buf.String()
}

func randBits256() Bits256 {
	v := Bits256{}
	_, _ = rand.Read(v[:])
	return v
}

func newPulsarData(pn PulseNumber, delta uint16) *PulseData {
	if delta == 0 {
		panic("delta cant be zero")
	}
	s := PulseData{
		PulseNumber: pn,
		PulseDataExt: PulseDataExt{
			PulseEntropy:   randBits256(),
			Timestamp:      uint32(time.Now().Unix()),
			NextPulseDelta: delta,
			PrevPulseDelta: 0,
		},
	}
	s.PulseEpoch = s.PulseNumber.AsUint32()
	return &s
}

func newEphemeralData(pn PulseNumber) *PulseData {
	s := PulseData{
		PulseNumber: pn,
		PulseDataExt: PulseDataExt{
			PulseEpoch:     EphemeralPulseEpoch,
			Timestamp:      0,
			NextPulseDelta: 1,
			PrevPulseDelta: 0,
		},
	}
	fixedPulseEntropy(&s.PulseEntropy, s.PulseNumber)
	return &s
}

/* This function has a fixed implementation and MUST remain unchanged as some elements of Consesnsus rely on identical behavior of this functions. */
func fixedPulseEntropy(v *Bits256, pn PulseNumber) {

	FillBitsWithStaticNoise(uint32(pn), (*v)[:])
}

func (r *PulseData) EnsurePulseData() {
	if !r.PulseNumber.IsTimePulse() {
		panic("incorrect pulse number")
	}
	if !PulseNumberOfUint32(r.PulseEpoch).IsSpecialOrTimePulse() {
		panic("incorrect pulse epoch")
	}
	if r.NextPulseDelta == 0 {
		panic("next delta can't be zero")
	}
}

func (r *PulseData) IsValidPulseData() bool {
	if !r.PulseNumber.IsTimePulse() {
		return false
	}
	if !PulseNumberOfUint32(r.PulseEpoch).IsSpecialOrTimePulse() {
		return false
	}
	if r.NextPulseDelta == 0 {
		return false
	}
	return true
}

func (r *PulseData) IsEmpty() bool {
	return r.PulseNumber.IsUnknown()
}

func (r *PulseData) IsValidExpectedPulseData() bool {
	if !r.PulseNumber.IsTimePulse() {
		return false
	}
	if !PulseNumberOfUint32(r.PulseEpoch).IsSpecialOrTimePulse() {
		return false
	}
	if r.PrevPulseDelta != 0 {
		return false
	}
	return true
}

func (r *PulseData) EnsurePulsarData() {
	if !PulseNumberOfUint32(r.PulseEpoch).IsTimePulse() {
		panic("incorrect pulse epoch by pulsar")
	}
	r.EnsurePulseData()
}

func (r *PulseData) IsValidPulsarData() bool {
	if !PulseNumberOfUint32(r.PulseEpoch).IsTimePulse() {
		return false
	}
	return r.IsValidPulseData()
}

func (r *PulseData) EnsureEphemeralData() {
	if r.PulseEpoch != EphemeralPulseEpoch {
		panic("incorrect pulse epoch")
	}
	r.EnsurePulseData()
}

func (r *PulseData) IsValidEphemeralData() bool {
	if r.PulseEpoch != EphemeralPulseEpoch {
		return false
	}
	return r.IsValidPulseData()
}

func (r *PulseData) IsFromPulsar() bool {
	return r.PulseNumber.IsTimePulse() && PulseNumber(r.PulseEpoch).IsTimePulse()
}

func (r *PulseData) IsFromEphemeral() bool {
	return r.PulseNumber.IsTimePulse() && r.PulseEpoch == EphemeralPulseEpoch
}

func (r *PulseData) GetStartOfEpoch() PulseNumber {
	ep := PulseNumberOfUint32(r.PulseEpoch)
	if ep.IsTimePulse() {
		return ep
	}
	return r.PulseNumber
}

func (r *PulseData) CreateNextPulse() *PulseData {
	if r.IsFromEphemeral() {
		return r.createNextEphemeralPulse()
	}
	return r.createNextPulsarPulse(r.NextPulseDelta)
}

func (r *PulseData) IsValidNext(n *PulseData) bool {
	if r.IsExpectedPulse() || r.GetNextPulseNumber() != n.PulseNumber || r.NextPulseDelta != n.PrevPulseDelta {
		return false
	}
	switch {
	case r.IsFromPulsar():
		return n.IsValidPulsarData()
	case r.IsFromEphemeral():
		return n.IsValidEphemeralData()
	}
	return n.IsValidPulseData()
}

func (r *PulseData) IsValidPrev(p *PulseData) bool {
	switch {
	case r.IsFirstPulse() || p.IsExpectedPulse() || p.GetNextPulseNumber() != r.PulseNumber || p.NextPulseDelta != r.PrevPulseDelta:
		return false
	case r.IsFromPulsar():
		return p.IsValidPulsarData()
	case r.IsFromEphemeral():
		return p.IsValidEphemeralData()
	default:
		return p.IsValidPulseData()
	}
}

func (r *PulseData) GetNextPulseNumber() PulseNumber {
	if r.IsExpectedPulse() {
		panic("illegal state")
	}
	return r.PulseNumber.Next(r.NextPulseDelta)
}

func (r *PulseData) CreateNextExpected() *PulseData {
	s := PulseData{
		PulseNumber: r.GetNextPulseNumber(),
		PulseDataExt: PulseDataExt{
			PrevPulseDelta: r.NextPulseDelta,
			NextPulseDelta: 0,
		},
	}
	if r.IsFromEphemeral() {
		s.PulseEpoch = r.PulseEpoch
	}
	return &s
}

func (r *PulseData) CreateNextEphemeralPulse() *PulseData {
	if !r.IsFromEphemeral() {
		panic("prev is not ephemeral")
	}
	return r.createNextEphemeralPulse()
}

func (r *PulseData) createNextEphemeralPulse() *PulseData {
	s := newEphemeralData(r.GetNextPulseNumber())
	s.PrevPulseDelta = r.NextPulseDelta
	return s
}

func (r *PulseData) CreateNextPulsarPulse(delta uint16) *PulseData {
	if r.IsFromEphemeral() {
		panic("prev is ephemeral")
	}
	return r.createNextPulsarPulse(delta)
}

func (r *PulseData) createNextPulsarPulse(delta uint16) *PulseData {
	s := newPulsarData(r.GetNextPulseNumber(), delta)
	s.PrevPulseDelta = r.NextPulseDelta
	return s
}

func (r *PulseData) GetPulseNumber() PulseNumber {
	return r.PulseNumber
}

func (r *PulseData) GetNextPulseDelta() uint16 {
	return r.NextPulseDelta
}

func (r *PulseData) GetPrevPulseDelta() uint16 {
	return r.PrevPulseDelta
}

func (r *PulseData) GetTimestamp() uint64 {
	return uint64(r.Timestamp)
}

func (r *PulseData) IsExpectedPulse() bool {
	return r.PulseNumber.IsTimePulse() && r.NextPulseDelta == 0
}

func (r *PulseData) IsFirstPulse() bool {
	return r.PulseNumber.IsTimePulse() && r.PrevPulseDelta == 0
}
