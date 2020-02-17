// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package insolar

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/pulse"
)

const (
	// PulseNumberSize is alias that was left for compatibility
	PulseNumberSize = pulse.NumberSize
	// EntropySize declares the number of bytes in the pulse entropy
	EntropySize = 64
	// OriginIDSize declares the number of bytes in the origin id
	OriginIDSize = 16
)

// Entropy is 64 random bytes used in every pseudo-random calculations.
type Entropy [EntropySize]byte

func (entropy Entropy) Marshal() ([]byte, error) { return entropy[:], nil }
func (entropy Entropy) MarshalTo(data []byte) (int, error) {
	copy(data, entropy[:])
	return EntropySize, nil
}
func (entropy *Entropy) Unmarshal(data []byte) error {
	if len(data) != EntropySize {
		return errors.New("Not enough bytes to unpack Entropy")
	}
	copy(entropy[:], data)
	return nil
}
func (entropy *Entropy) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, entropy)
}
func (entropy Entropy) Size() int { return EntropySize }
func (entropy Entropy) Compare(other Entropy) int {
	return bytes.Compare(entropy[:], other[:])
}
func (entropy Entropy) Equal(other Entropy) bool {
	return entropy.Compare(other) == 0
}

// PulseNumber is a sequential number of Pulse.
// Upper 2 bits are reserved for use in references (scope), must be zero otherwise.
// Valid Absolute PulseNumber must be >65536.
// If PulseNumber <65536 it is a relative PulseNumber
type PulseNumber = pulse.Number

// NewPulseNumber creates pulse number from bytes.
func NewPulseNumber(buf []byte) PulseNumber {
	return PulseNumber(binary.BigEndian.Uint32(buf))
}

func NewPulseNumberFromStr(pn string) (PulseNumber, error) {
	i, err := strconv.ParseUint(pn, 10, 32)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse pulse number")
	}
	return PulseNumber(i), nil
}

//go:generate minimock -i github.com/insolar/insolar/insolar.PulseManager -o ../testutils -s _mock.go -g

// PulseManager provides Ledger's methods related to Pulse.
type PulseManager interface {
	// Set set's new pulse and closes current jet drop. If dry is true, nothing will be saved to storage.
	Set(ctx context.Context, pulse Pulse) error
}

// PulseRange represents range of pulses.
type PulseRange struct {
	Begin PulseNumber
	End   PulseNumber
}

func (pr *PulseRange) String() string {
	return fmt.Sprintf("[%v:%v]", pr.Begin, pr.End)
}

// Pulse is base data structure for a pulse.
type Pulse struct {
	PulseNumber     PulseNumber
	PrevPulseNumber PulseNumber
	NextPulseNumber PulseNumber

	PulseTimestamp   int64
	EpochPulseNumber pulse.Epoch
	OriginID         [OriginIDSize]byte

	Entropy Entropy
	Signs   map[string]PulseSenderConfirmation
}

// PulseSenderConfirmation contains confirmations of the pulse from other pulsars
// Because the system is using BFT for consensus between pulsars, because of it
// All pulsar send to the chosen pulsar their confirmations
// Every node in the network can verify the signatures
type PulseSenderConfirmation struct {
	PulseNumber     PulseNumber
	ChosenPublicKey string
	Entropy         Entropy
	Signature       []byte
}

// GenesisPulse is a first pulse for the system
// because first 2 bits of pulse number and first 65536 pulses a are used by system needs and pulse numbers are related to the seconds of Unix time
// for calculation pulse numbers we use the formula = unix.Now() - firstPulseDate + 65536
var GenesisPulse = &Pulse{
	PulseNumber:      pulse.MinTimePulse,
	Entropy:          [EntropySize]byte{},
	EpochPulseNumber: pulse.MinTimePulse,
	PulseTimestamp:   pulse.UnixTimeOfMinTimePulse,
}

// EphemeralPulse is used for discovery network bootstrap
var EphemeralPulse = &Pulse{
	PulseNumber:      pulse.MinTimePulse,
	Entropy:          [EntropySize]byte{},
	EpochPulseNumber: pulse.EphemeralPulseEpoch,
	PulseTimestamp:   pulse.UnixTimeOfMinTimePulse,
}
