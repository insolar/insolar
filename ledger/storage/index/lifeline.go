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

package index

import (
	"bytes"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage/record"
	"github.com/ugorji/go/codec"
)

// ObjectLifeline represents meta information for record object.
type ObjectLifeline struct {
	LatestState         *core.RecordID // Amend or activate record.
	LatestStateApproved *core.RecordID // State approved by VM.
	ChildPointer        *core.RecordID // Meta record about child activation.
	Parent              core.RecordRef
	Delegates           map[core.RecordRef]core.RecordRef
	State               record.State
}

// EncodeObjectLifeline converts lifeline index into binary format.
func EncodeObjectLifeline(index *ObjectLifeline) ([]byte, error) {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.CborHandle{})
	err := enc.Encode(index)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// DecodeObjectLifeline converts byte array into lifeline index struct.
func DecodeObjectLifeline(buf []byte) (*ObjectLifeline, error) {
	dec := codec.NewDecoder(bytes.NewReader(buf), &codec.CborHandle{})
	var index ObjectLifeline
	err := dec.Decode(&index)
	if err != nil {
		return nil, err
	}
	return &index, nil
}
