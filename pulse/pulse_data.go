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

package pulse

import (
	"fmt"
	"strings"
	"time"

	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
)

const InvalidPulseEpoch uint32 = 0
const EphemeralPulseEpoch = InvalidPulseEpoch + 1

var _ DataReader = &Data{}

type Data struct {
	PulseNumber Number
	DataExt
}

type DataHolder interface {
	GetPulseNumber() Number
	GetPulseData() Data
	GetPulseDataDigest() cryptkit.DigestHolder
}

type DataExt struct {
	// ByteSize=44
	PulseEpoch     uint32
	PulseEntropy   longbits.Bits256
	NextPulseDelta uint16
	PrevPulseDelta uint16
	Timestamp      uint32
}

type DataReader interface {
	GetPulseNumber() Number
	GetStartOfEpoch() Number
	// GetPulseEntropy()	[4]uint64
	GetNextPulseDelta() uint16
	GetPrevPulseDelta() uint16
	GetTimestamp() uint64
	IsExpectedPulse() bool
	IsFromEphemeral() bool
}

func NewFirstPulsarData(delta uint16, entropy longbits.Bits256) Data {
	return newPulsarData(OfNow(), delta, entropy)
}

func NewPulsarData(pn Number, deltaNext uint16, deltaPrev uint16, entropy longbits.Bits256) Data {
	r := newPulsarData(pn, deltaNext, entropy)
	r.PrevPulseDelta = deltaPrev
	return r
}

func NewFirstEphemeralData() Data {
	return newEphemeralData(MinTimePulse)
}

type EntropyFunc func() longbits.Bits256

func (r Data) String() string {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprint(r.PulseNumber))

	ep := OfUint32(r.PulseEpoch)
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

func newPulsarData(pn Number, delta uint16, entropy longbits.Bits256) Data {
	if delta == 0 {
		panic("delta cant be zero")
	}
	return Data{
		PulseNumber: pn,
		DataExt: DataExt{
			PulseEpoch:     pn.AsUint32(),
			PulseEntropy:   entropy,
			Timestamp:      uint32(time.Now().Unix()),
			NextPulseDelta: delta,
			PrevPulseDelta: 0,
		},
	}
}

func newEphemeralData(pn Number) Data {
	s := Data{
		PulseNumber: pn,
		DataExt: DataExt{
			PulseEpoch:     EphemeralPulseEpoch,
			Timestamp:      0,
			NextPulseDelta: 1,
			PrevPulseDelta: 0,
		},
	}
	fixedPulseEntropy(&s.PulseEntropy, s.PulseNumber)
	return s
}

/* This function has a fixed implementation and MUST remain unchanged as some elements of Consensus rely on identical behavior of this functions. */
func fixedPulseEntropy(v *longbits.Bits256, pn Number) {
	longbits.FillBitsWithStaticNoise(uint32(pn), (*v)[:])
}

func (r Data) EnsurePulseData() {
	if !r.PulseNumber.IsTimePulse() {
		panic("incorrect pulse number")
	}
	if !OfUint32(r.PulseEpoch).IsSpecialOrTimePulse() {
		panic("incorrect pulse epoch")
	}
	if r.NextPulseDelta == 0 {
		panic("next delta can't be zero")
	}
}

func (r Data) IsValidPulseData() bool {
	if !r.PulseNumber.IsTimePulse() {
		return false
	}
	if !OfUint32(r.PulseEpoch).IsSpecialOrTimePulse() {
		return false
	}
	if r.NextPulseDelta == 0 {
		return false
	}
	return true
}

func (r Data) IsEmpty() bool {
	return r.PulseNumber.IsUnknown()
}

func (r Data) IsEmptyWithEpoch(epoch uint32) bool {
	return r.PulseNumber.IsUnknown() && r.PulseEpoch == epoch
}

func (r Data) IsValidExpectedPulseData() bool {
	if !r.PulseNumber.IsTimePulse() {
		return false
	}
	if !OfUint32(r.PulseEpoch).IsSpecialOrTimePulse() {
		return false
	}
	if r.PrevPulseDelta != 0 {
		return false
	}
	return true
}

func (r Data) EnsurePulsarData() {
	if !OfUint32(r.PulseEpoch).IsTimePulse() {
		panic("incorrect pulse epoch by pulsar")
	}
	r.EnsurePulseData()
}

func (r Data) IsValidPulsarData() bool {
	if !OfUint32(r.PulseEpoch).IsTimePulse() {
		return false
	}
	return r.IsValidPulseData()
}

func (r Data) EnsureEphemeralData() {
	if r.PulseEpoch != EphemeralPulseEpoch {
		panic("incorrect pulse epoch")
	}
	r.EnsurePulseData()
}

func (r Data) IsValidEphemeralData() bool {
	if r.PulseEpoch != EphemeralPulseEpoch {
		return false
	}
	return r.IsValidPulseData()
}

func (r Data) IsFromPulsar() bool {
	return r.PulseNumber.IsTimePulse() && OfUint32(r.PulseEpoch).IsTimePulse()
}

func (r Data) IsFromEphemeral() bool {
	return r.PulseNumber.IsTimePulse() && r.PulseEpoch == EphemeralPulseEpoch
}

func (r Data) GetStartOfEpoch() Number {
	ep := OfUint32(r.PulseEpoch)
	if r.PulseNumber.IsTimePulse() {
		return ep
	}
	return r.PulseNumber
}

func (r Data) CreateNextPulse(entropyGen EntropyFunc) Data {
	if r.IsFromEphemeral() {
		return r.createNextEphemeralPulse()
	}
	return r.createNextPulsarPulse(r.NextPulseDelta, entropyGen)
}

func (r Data) IsValidNext(n Data) bool {
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

func (r Data) IsValidPrev(p Data) bool {
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

func (r Data) GetNextPulseNumber() Number {
	if r.IsExpectedPulse() {
		panic("illegal state")
	}
	return r.PulseNumber.Next(r.NextPulseDelta)
}

func (r Data) GetPrevPulseNumber() Number {
	if r.IsFirstPulse() {
		panic("illegal state")
	}
	return r.PulseNumber.Prev(r.PrevPulseDelta)
}

func (r Data) CreateNextExpected() Data {
	s := Data{
		PulseNumber: r.GetNextPulseNumber(),
		DataExt: DataExt{
			PrevPulseDelta: r.NextPulseDelta,
			NextPulseDelta: 0,
		},
	}
	if r.IsFromEphemeral() {
		s.PulseEpoch = r.PulseEpoch
	}
	return s
}

func (r Data) CreateNextEphemeralPulse() Data {
	if !r.IsFromEphemeral() {
		panic("prev is not ephemeral")
	}
	return r.createNextEphemeralPulse()
}

func (r Data) createNextEphemeralPulse() Data {
	s := newEphemeralData(r.GetNextPulseNumber())
	s.PrevPulseDelta = r.NextPulseDelta
	return s
}

func (r Data) CreateNextPulsarPulse(delta uint16, entropyGen EntropyFunc) Data {
	if r.IsFromEphemeral() {
		panic("prev is ephemeral")
	}
	return r.createNextPulsarPulse(delta, entropyGen)
}

func (r Data) createNextPulsarPulse(delta uint16, entropyGen EntropyFunc) Data {
	s := newPulsarData(r.GetNextPulseNumber(), delta, entropyGen())
	s.PrevPulseDelta = r.NextPulseDelta
	return s
}

func (r Data) GetPulseNumber() Number {
	return r.PulseNumber
}

func (r Data) GetNextPulseDelta() uint16 {
	return r.NextPulseDelta
}

func (r Data) GetPrevPulseDelta() uint16 {
	return r.PrevPulseDelta
}

func (r Data) GetTimestamp() uint64 {
	return uint64(r.Timestamp)
}

func (r Data) IsExpectedPulse() bool {
	return r.PulseNumber.IsTimePulse() && r.NextPulseDelta == 0
}

func (r Data) IsFirstPulse() bool {
	return r.PulseNumber.IsTimePulse() && r.PrevPulseDelta == 0
}
