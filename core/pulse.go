/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package core

import (
	"encoding/binary"
	"time"
)

const (
	// PulseNumberSize declares the number of bytes in the pulse number
	PulseNumberSize = 4
	// EntropySize declares the number of bytes in the pulse entropy
	EntropySize = 64
)

// Entropy is 64 random bytes used in every pseudo-random calculations.
type Entropy [EntropySize]byte

// PulseNumber is a sequential number of Pulse.
// Upper 2 bits are reserved for use in references (scope), must be zero otherwise.
// Valid Absolute PulseNum must be >65536.
// If PulseNum <65536 it is a relative PulseNum
type PulseNumber uint32

// Bytes serializes pulse number.
func (pn PulseNumber) Bytes() []byte {
	buff := make([]byte, PulseNumberSize)
	binary.BigEndian.PutUint32(buff, uint32(pn))
	return buff
}

// Bytes2PulseNumber deserialize pulse number.
func Bytes2PulseNumber(buf []byte) PulseNumber {
	return PulseNumber(binary.BigEndian.Uint32(buf))
}

// Pulse is base data structure for a pulse.
type Pulse struct {
	PulseNumber     PulseNumber
	NextPulseNumber PulseNumber
	Entropy         Entropy
	Signs           map[string]PulseSenderConfirmation
}

func (p *Pulse) PulseDuration() time.Duration {
	return time.Second * time.Duration(p.NextPulseNumber-p.PulseNumber)
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

// FirstPulseDate is the hardcoded date of the first pulse
const firstPulseDate = 1535760000 //09/01/2018 @ 12:00am (UTC)
// FirstPulseNumber is the hardcoded first pulse number. Because first 65536 numbers are saved for the system's needs
const FirstPulseNumber = 65537

// GenesisPulse is a first pulse for the system
// because first 2 bits of pulse number and first 65536 pulses a are used by system needs and pulse numbers are related to the seconds of Unix time
// for calculation pulse numbers we use the formula = unix.Now() - firstPulseDate + 65536
var GenesisPulse = &Pulse{PulseNumber: FirstPulseNumber, Entropy: [EntropySize]byte{}}

// CalculatePulseNumber is helper for calculating next pulse number, when a network is being started
func CalculatePulseNumber(now time.Time) PulseNumber {
	return PulseNumber(now.Unix() - firstPulseDate + FirstPulseNumber)
}
