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

var convertTests = []struct {
	name string
	key  Key
	id   ID
}{
	{
		key: Key{Pulse: 10, Hash: str2Bytes("21853428b06925493bf23d2c5ba76ee86e3e3c1a13fe164307250193")},
		id:  str2ID("0000000a" + "21853428b06925493bf23d2c5ba76ee86e3e3c1a13fe164307250193"),
	},
}

func Test_KeyIDConversion(t *testing.T) {
	for _, tt := range convertTests {
		t.Run(tt.name, func(t *testing.T) {
			gotID := Key2ID(tt.key)
			gotKey := ID2Key(gotID)
			assert.Equal(t, tt.key, gotKey)
			assert.Equal(t, tt.id, gotID)
		})
	}
}

func Test_RecordByTypeIDPanic(t *testing.T) {
	assert.Panics(t, func() { getRecordByTypeID(0) })
}

var type2idTests = []struct {
	typ string
	rec Record
	id  TypeID
}{
	{"RequestRecord", &RequestRecord{}, requestRecordID},
	{"CallRequest", &CallRequest{}, callRequestID},
	{"LockUnlockRequest", &LockUnlockRequest{}, lockUnlockRequestID},
	{"ReadRecordRequest", &ReadRecordRequest{}, readRecordRequestID},
	{"ReadObject", &ReadObject{}, readObjectID},
	{"ReadObjectComposite", &ReadObjectComposite{}, readObjectCompositeID},

	// result records
	// case resultRecordID:
	{"WipeOutRecord", &WipeOutRecord{}, wipeOutRecordID},
	{"ReadRecordResult", &ReadRecordResult{}, readRecordResultID},
	{"StatelessCallResult", &StatelessCallResult{}, statelessCallResultID},
	{"StatelessExceptionResult", &StatelessExceptionResult{}, statelessExceptionResultID},
	{"ReadObjectResult", &ReadObjectResult{}, readObjectResultID},
	{"SpecialResult", &SpecialResult{}, specialResultID},
	{"LockUnlockResult", &LockUnlockResult{}, lockUnlockResultID},
	{"RejectionResult", &RejectionResult{}, rejectionResultID},
	{"ActivationRecord", &ActivationRecord{}, activationRecordID},
	{"ClassActivateRecord", &ClassActivateRecord{}, classActivateRecordID},
	{"ObjectActivateRecord", &ObjectActivateRecord{}, objectActivateRecordID},
	{"CodeRecord", &CodeRecord{}, codeRecordID},
	{"AmendRecord", &AmendRecord{}, amendRecordID},
	{"ClassAmendRecord", &ClassAmendRecord{}, classAmendRecordID},
	{"MemoryMigrationCode", &MemoryMigrationCode{}, memoryMigrationCodeID},
	{"DeactivationRecord", &DeactivationRecord{}, deactivationRecordID},
	{"ObjectAmendRecord", &ObjectAmendRecord{}, objectAmendRecordID},
	{"StatefulCallResult", &StatefulCallResult{}, statefulCallResultID},
	{"StatefulExceptionResult", &StatefulExceptionResult{}, statefulExceptionResultID},
	{"EnforcedObjectAmendRecord", &EnforcedObjectAmendRecord{}, enforcedObjectAmendRecordID},
	{"ObjectAppendRecord", &ObjectAppendRecord{}, objectAppendRecordID},
}

func Test_TypeIDConversion(t *testing.T) {
	for _, tt := range type2idTests {
		t.Run(tt.typ, func(t *testing.T) {
			gotRecTypeID := getTypeIDbyRecord(tt.rec)
			gotRecord := getRecordByTypeID(tt.id)
			assert.Equal(t, "*record."+tt.typ, fmt.Sprintf("%T", gotRecord))
			assert.Equal(t, tt.id, gotRecTypeID)
		})
	}
}

var serializeTests = []struct {
	name         string
	rec          Record
	expectTypeID TypeID
}{
	{
		"RequestRecord_WithRequester",
		&RequestRecord{
			Requester: Reference{
				Domain: str2ID("0000000a" + "21853428b06925493bf23d2c5ba76ee86e3e3c1a13fe164307250193"),
			},
		},
		requestRecordID,
	},
	{
		"RequestRecord_WithTarget",
		&RequestRecord{
			Target: Reference{
				Domain: str2ID("0000000a" + "21853428b06925493bf23d2c5ba76ee86e3e3c1a13fe164307250193"),
			},
		},
		requestRecordID,
	},
}

func Test_EncodeToRaw(t *testing.T) {
	for _, tt := range serializeTests {
		t.Run(tt.name, func(t *testing.T) {
			raw, err := EncodeToRaw(tt.rec)
			if err != nil {
				panic(err)
			}
			// fmt.Println(tt.name, "got", spew.Sdump(raw))
			assert.Equal(t, tt.expectTypeID, raw.Type)
		})
	}
}
