/*
 *    Copyright 2019 Insolar Technologies
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

package jet

import (
	"bytes"

	"github.com/insolar/insolar/core"
	"github.com/ugorji/go/codec"
)

// // ZeroJetID is value of an empty Jet ID
// var ZeroJetID = *NewID(0, nil)
//
// // NewID generates RecordID for jet.
// func NewID(depth uint8, prefix []byte) *core.JetID {
// 	var id core.JetID
// 	copy(id[:core.PulseNumberSize], core.PulseNumberJet.Bytes())
// 	id[core.PulseNumberSize] = depth
// 	copy(id[core.PulseNumberSize+1:], prefix)
// 	return &id
// }

// IDSet is an alias for map[ID]struct{}
type IDSet map[core.RecordID]struct{}

// Has checks if passed id is in IDSet set.
func (j IDSet) Has(id core.RecordID) bool {
	_, ok := j[id]
	return ok
}

// Bytes serializes pulse.
func (j IDSet) Bytes() []byte {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.CborHandle{})
	enc.MustEncode(j)
	return buf.Bytes()
}

// // Jet extracts depth and prefix from jet id.
// func Jet(id core.JetID) (uint8, []byte) {
// 	if core.RecordID(id).Pulse() != core.PulseNumberJet {
// 		panic(fmt.Sprintf("provided id %b is not a jet id", id))
// 	}
// 	return id[core.PulseNumberSize], id[core.PulseNumberSize+1:]
// }
//
// func Parent(id JetID) core.JetID {
// 	depth, prefix := Jet(id)
// 	if depth == 0 {
// 		return id
// 	}
//
// 	return *NewID(depth-1, ResetBits(prefix, depth-1))
// }
//
// // ResetBits returns a new byte slice with all bits in 'value' reset, starting from 'start' number of bit. If 'start'
// // is bigger than len(value), the original slice will be returned.
// func ResetBits(value []byte, start uint8) []byte {
// 	if int(start) >= len(value)*8 {
// 		return value
// 	}
//
// 	startByte := start / 8
// 	startBit := start % 8
//
// 	result := make([]byte, len(value))
// 	copy(result, value[:startByte])
//
// 	// Reset bits in starting byte.
// 	mask := byte(0xFF)
// 	mask <<= 8 - byte(startBit)
// 	result[startByte] = value[startByte] & mask
//
// 	return result
// }
