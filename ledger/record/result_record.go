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

type ReasonCode uint

type ResultRecord struct {
	AppDataRecord

	RequestRecord RecordReference
}

type WipeOutRecord struct {
	ResultRecord

	Replacement RecordReference
	WipedHash   RecordHash
}

type StatelessResult struct {
	ResultRecord
}

type ReadRecordResult struct {
	StatelessResult

	RecordBody []byte
}

type StatelessCallResult struct {
	StatelessResult

	resultMemory Memory
}

type StatelessExceptionResult struct {
	StatelessCallResult

	ExceptionType RecordReference
}

type ReadObjectResult struct {
	StatelessResult

	State            int
	MomoryProjection Memory
}

type SpecialResult struct {
	ResultRecord

	ReasonCode ReasonCode
}

type LockUnlockResult struct {
	SpecialResult
}

type RejectionResult struct {
	SpecialResult
}

type StatefulResult struct {
	ResultRecord
}

type ActivationRecord struct {
	StatefulResult

	GoverningDomain RecordReference
}

type ClassActivateRecord struct {
	ActivationRecord

	CodeRecord    RecordReference
	DefaultMemory Memory
}

type ObjectActivateRecord struct {
	ActivationRecord

	ClassActivateRecord RecordReference
	Memory              Memory
}

type StorageRecord struct {
	StatefulResult
}

type CodeRecord struct {
	StorageRecord

	Interfaces   []RecordReference
	TargetedCode [][]byte // []MachineBinaryCode
	SourceCode   string   // ObjectSourceCode
}

type AmendRecord struct {
	StatefulResult

	BaseRecord    RecordReference
	AmendedRecord RecordReference
}

type ClassAmendRecord struct {
	AmendRecord

	NewCode []byte // ObjectBinaryCode
}

func (r *ClassAmendRecord) MigrationCodes() []*MemoryMigrationCode {
	panic("not implemented")
}

type MemoryMigrationCode struct {
	ClassAmendRecord

	GeneratedByClassRecord RecordReference
	MigrationCodeRecord    RecordReference
}

type DeactivationRecord struct {
	AmendRecord
}

type ObjectAmendRecord struct {
	AmendRecord

	NewMemory Memory
}

type StatefulCallResult struct {
	ObjectAmendRecord

	ResultMemory Memory
}

type StatefulExceptionResult struct {
	StatefulCallResult

	ExceptionType RecordReference
}

type EnforcedObjectAmendRecord struct {
	ObjectAmendRecord
}

type ObjectAppendRecord struct {
	AmendRecord

	AppendMemory Memory
}
