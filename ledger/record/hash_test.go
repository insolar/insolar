/*
 *    Copyright 2018 INS Ecosystem
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

package record

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type recordgen func() Record

var emptyRecordsGens = []recordgen{
	func() Record { return &RequestRecord{} },
	func() Record { return &CallRequest{} },
	func() Record { return &LockUnlockRequest{} },
	func() Record { return &ReadRecordRequest{} },
	func() Record { return &ReadObject{} },
	func() Record { return &ReadObjectComposite{} },
	// result records
	func() Record { return &ResultRecord{} },
	func() Record { return &WipeOutRecord{} },
	func() Record { return &ReadRecordResult{} },
	func() Record { return &StatelessCallResult{} },
	func() Record { return &StatelessExceptionResult{} },
	func() Record { return &ReadObjectResult{} },
	func() Record { return &SpecialResult{} },
	func() Record { return &LockUnlockResult{} },
	func() Record { return &RejectionResult{} },
	func() Record { return &ActivationRecord{} },
	func() Record { return &ClassActivateRecord{} },
	func() Record { return &ObjectActivateRecord{} },
	func() Record { return &CodeRecord{} },
	func() Record { return &AmendRecord{} },
	func() Record { return &ClassAmendRecord{} },
	func() Record { return &DeactivationRecord{} },
	func() Record { return &ObjectAmendRecord{} },
	func() Record { return &StatefulCallResult{} },
	func() Record { return &StatefulExceptionResult{} },
	func() Record { return &EnforcedObjectAmendRecord{} },
	func() Record { return &ObjectAppendRecord{} },
}

func Test_HashesNotTheSameOnDifferentTypes(t *testing.T) {
	found := make(map[string]string)
	for _, recFn := range emptyRecordsGens {
		rec := recFn()
		recType := fmt.Sprintf("%T", rec)
		hashBytes := SHA3Hash224(rec)

		hashHex := fmt.Sprintf("%x", hashBytes)
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
		hashHex := fmt.Sprintf("%x", SHA3Hash224(rec))
		hashes[i] = hashHex
	}

	// same struct with different should produce the same hashes
	for i, recFn := range emptyRecordsGens {
		rec := recFn()
		hashHex := fmt.Sprintf("%x", SHA3Hash224(rec))
		assert.Equal(t, hashes[i], hashHex)
	}
}

var hashtestsRecordsMutate = []struct {
	typ     string
	records []Record
}{
	{
		"ReadObject",
		[]Record{
			&ReadObject{},
			&ReadObject{ProjectionType: 1},
			&ReadObject{ProjectionType: 2},
		},
	},
	{
		"ReadObjectComposite",
		[]Record{
			&ReadObjectComposite{},
			&ReadObjectComposite{ReadObject: ReadObject{ProjectionType: 1}},
			&ReadObjectComposite{ReadObject: ReadObject{ProjectionType: 2}},
			&ReadObjectComposite{CompositeType: Reference{
				Domain: str2ID("0000000a" + "21853428b06925493bf23d2c5ba76ee86e3e3c1a13fe164307250193"),
				Record: str2ID("0000000a" + "21853428b06925493bf23d2c5ba76ee86e3e3c1a13fe164307250193"),
			}},
			&ReadObjectComposite{CompositeType: Reference{
				Domain: str2ID("0000000a" + "21853428b06925493bf23d2c5ba76ee86e3e3c1a13fe164307250193"),
			}},
			&ReadObjectComposite{CompositeType: Reference{
				Record: str2ID("0000000a" + "21853428b06925493bf23d2c5ba76ee86e3e3c1a13fe164307250193"),
			}},
		},
	},
}

func Test_CBORhashesMutation(t *testing.T) {
	for _, tt := range hashtestsRecordsMutate {
		found := make(map[string]string)
		for _, rec := range tt.records {
			h := SHA3Hash224(rec)
			hHex := fmt.Sprintf("%x", h)

			typ, ok := found[hHex]
			if !ok {
				found[hHex] = tt.typ
				continue
			}
			t.Errorf("%s failed: found %s hash for \"%s\" test, should not repeats in sha3hash224tests", tt.typ, hHex, typ)
		}
	}
}
