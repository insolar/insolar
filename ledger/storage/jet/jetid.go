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

package jet

import (
	"bytes"

	"github.com/insolar/insolar/core"
	"github.com/ugorji/go/codec"
)

// ZeroJetID is value of an empty Jet ID
var ZeroJetID = *NewID(0, nil)

// NewID generates RecordID for jet.
func NewID(depth uint8, prefix []byte) *core.RecordID {
	var id core.RecordID
	copy(id[:core.PulseNumberSize], core.PulseNumberJet.Bytes())
	id[core.PulseNumberSize] = depth
	copy(id[core.PulseNumberSize+1:], prefix)
	return &id
}

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

// Jet extracts depth and prefix from jet id.
func Jet(id core.RecordID) (uint8, []byte) {
	if id.Pulse() != core.PulseNumberJet {
		panic("provided id in not a jet id")
	}
	return id[core.PulseNumberSize], id[core.PulseNumberSize+1:]
}
