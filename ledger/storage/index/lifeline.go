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

package index

import (
	"bytes"
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage/record"
	"github.com/ugorji/go/codec"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/storage/index.Accessor -o ./ -s _mock.go

// Accessor provides info about Index-values from storage.
type Accessor interface {
	// ForID returns Index for provided id.
	ForID(ctx context.Context, id core.RecordID) (ObjectLifeline, error)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/storage/index.Modifier -o ./ -s _mock.go

// Modifier provides methods for setting Index-values to storage.
type Modifier interface {
	// Set saves new Index-value in storage.
	Set(ctx context.Context, id core.RecordID, index ObjectLifeline) error
}

// ObjectLifeline represents meta information for record object.
type ObjectLifeline struct {
	LatestState         *core.RecordID // Amend or activate record.
	LatestStateApproved *core.RecordID // State approved by VM.
	ChildPointer        *core.RecordID // Meta record about child activation.
	Parent              core.RecordRef
	Delegates           map[core.RecordRef]core.RecordRef
	State               record.State
	LatestUpdate        core.PulseNumber
	JetID               core.JetID
}

// Encode converts lifeline index into binary format.
func Encode(index ObjectLifeline) []byte {
	buff := bytes.NewBuffer(nil)
	enc := codec.NewEncoder(buff, &codec.CborHandle{})
	enc.MustEncode(index)

	return buff.Bytes()
}

// Decode converts byte array into lifeline index struct.
func Decode(buff []byte) (index ObjectLifeline) {
	dec := codec.NewDecoderBytes(buff, &codec.CborHandle{})
	dec.MustDecode(&index)

	return
}

// Clone returns copy of argument idx value.
func Clone(idx ObjectLifeline) ObjectLifeline {
	if idx.LatestState != nil {
		tmp := *idx.LatestState
		idx.LatestState = &tmp
	}

	if idx.LatestStateApproved != nil {
		tmp := *idx.LatestStateApproved
		idx.LatestStateApproved = &tmp
	}

	if idx.ChildPointer != nil {
		tmp := *idx.ChildPointer
		idx.ChildPointer = &tmp
	}

	if idx.Delegates != nil {
		cp := make(map[core.RecordRef]core.RecordRef)
		for k, v := range idx.Delegates {
			cp[k] = v
		}
		idx.Delegates = cp
	}

	return idx
}
