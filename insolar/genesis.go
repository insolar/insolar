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

package insolar

import (
	"io"
)

type genesisBinary []byte

// GenesisRecord is initial chain record.
var GenesisRecord genesisBinary = []byte{0xAC}

// WriteHashData implements record.VirtualRecord.
func (r genesisBinary) WriteHashData(w io.Writer) (int, error) {
	return w.Write(r)
}

// ID returns genesis record id.
func (r genesisBinary) ID() ID {
	return *NewID(GenesisPulse.PulseNumber, r)
}

// Ref returns genesis record reference.
func (r genesisBinary) Ref() Reference {
	id := r.ID()
	return *NewReference(id, id)
}
