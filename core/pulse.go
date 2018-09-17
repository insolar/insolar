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
)

const (
	PulseNumberSize = 4
	EntropySize     = 64
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

// Bytes2PulseNumber deserializes pulse number.
func Bytes2PulseNumber(buf []byte) PulseNumber {
	return PulseNumber(binary.BigEndian.Uint32(buf))
}

// Pulse is base data structure for a pulse.
type Pulse struct {
	PulseNumber PulseNumber
	Entropy     Entropy
}

type PulseManager interface {
	// Current returns current pulse structure.
	Current() (*Pulse, error)

	// Set set's new pulse and closes current jet drop.
	Set(pulse Pulse) error
}
