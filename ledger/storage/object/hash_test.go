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

package object

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/insolar/insolar/insolar/record"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/assert"
)

type recordgen func() record.VirtualRecord

var emptyRecordsGens = []recordgen{
	// request records
	func() record.VirtualRecord { return &RequestRecord{} },
	// result records
	func() record.VirtualRecord { return &ActivateRecord{} },
	func() record.VirtualRecord { return &CodeRecord{} },
	func() record.VirtualRecord { return &DeactivationRecord{} },
	func() record.VirtualRecord { return &AmendRecord{} },
	func() record.VirtualRecord { return &TypeRecord{} },
	func() record.VirtualRecord { return &ChildRecord{} },
	func() record.VirtualRecord { return &GenesisRecord{insolar.GenesisRecord} },
}

func getRecordHashData(rec record.VirtualRecord) []byte {
	buff := bytes.NewBuffer(nil)
	rec.WriteHashData(buff)
	return buff.Bytes()
}

func Test_HashesNotTheSameOnDifferentTypes(t *testing.T) {
	found := make(map[string]string)
	for _, recFn := range emptyRecordsGens {
		rec := recFn()
		recType := fmt.Sprintf("%T", rec)
		hashHex := fmt.Sprintf("%x", getRecordHashData(rec))
		// fmt.Println(recType, "=>", hashHex)
		typename, ok := found[hashHex]
		if !ok {
			found[hashHex] = recType
			continue
		}
		t.Errorf("same hashes for %s and %s types, empty struct with different types should not be the same", recType, typename)
	}
}

func Test_HashesTheSame(t *testing.T) {
	hashes := make([]string, len(emptyRecordsGens))
	for i, recFn := range emptyRecordsGens {
		rec := recFn()
		hashHex := fmt.Sprintf("%x", getRecordHashData(rec))
		hashes[i] = hashHex
	}

	// same struct with different should produce the same hashes
	for i, recFn := range emptyRecordsGens {
		rec := recFn()
		hashHex := fmt.Sprintf("%x", getRecordHashData(rec))
		assert.Equal(t, hashes[i], hashHex)
	}
}

var pcs = platformpolicy.NewPlatformCryptographyScheme()
var hashtestsRecordsMutate = []struct {
	typ     string
	records []record.VirtualRecord
}{
	{
		"CodeRecord",
		[]record.VirtualRecord{
			&CodeRecord{},
			&CodeRecord{Code: CalculateIDForBlob(pcs, insolar.GenesisPulse.PulseNumber, []byte{1, 2, 3})},
			&CodeRecord{
				Code: CalculateIDForBlob(pcs, insolar.GenesisPulse.PulseNumber, []byte{1, 2, 3}),
				SideEffectRecord: SideEffectRecord{
					Domain: insolar.Reference{1, 2, 3},
				},
			},
		},
	},
}

func Test_CBORhashesMutation(t *testing.T) {
	for _, tt := range hashtestsRecordsMutate {
		found := make(map[string]string)
		for _, rec := range tt.records {
			h := getRecordHashData(rec)
			hHex := fmt.Sprintf("%x", h)

			typ, ok := found[hHex]
			if !ok {
				found[hHex] = tt.typ
				continue
			}
			t.Errorf("%s failed: found %s hash for \"%s\" test", tt.typ, hHex, typ)
		}
	}
}
