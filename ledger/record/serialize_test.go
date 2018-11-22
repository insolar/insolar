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

package record

import (
	"fmt"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/assert"
)

func Test_RecordByTypeIDPanic(t *testing.T) {
	assert.Panics(t, func() { getRecordByTypeID(0) })
}

var type2idTests = []struct {
	typ string
	rec Record
	id  TypeID
}{
	// request records
	{"CallRequest", &CallRequest{}, typeCallRequest},

	// result records
	{"ObjectActivateRecord", &ObjectActivateRecord{}, typeActivate},
	{"CodeRecord", &CodeRecord{}, typeCode},
	{"DeactivationRecord", &DeactivationRecord{}, typeDeactivate},
	{"ObjectAmendRecord", &ObjectAmendRecord{}, typeAmend},
	{"TypeRecord", &TypeRecord{}, typeType},
	{"ChildRecord", &ChildRecord{}, typeChild},
	{"GenesisRecord", &GenesisRecord{}, typeGenesis},
}

func Test_TypeIDConversion(t *testing.T) {
	for _, tt := range type2idTests {
		t.Run(tt.typ, func(t *testing.T) {
			gotRecTypeID := tt.rec.Type()
			gotRecord := getRecordByTypeID(tt.id)
			assert.Equal(t, "*record."+tt.typ, fmt.Sprintf("%T", gotRecord))
			assert.Equal(t, tt.id, gotRecTypeID)
		})
	}
}

func TestSerializeDeserializeRecord(t *testing.T) {
	cs := platformpolicy.NewPlatformCryptographyScheme()

	rec := ObjectActivateRecord{
		ObjectStateRecord: ObjectStateRecord{
			Memory: CalculateIDForBlob(cs, core.GenesisPulse.PulseNumber, []byte{1, 2, 3}),
		},
	}
	serialized := SerializeRecord(&rec)
	deserialized := DeserializeRecord(serialized)
	assert.Equal(t, rec, *deserialized.(*ObjectActivateRecord))
}
