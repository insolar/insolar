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
	"encoding/binary"
	"io"
)

// ReasonCode is an error reason code.
type ReasonCode uint32

// ResultRecord is a common type for all results.
type ResultRecord struct {
	RequestRecord Reference
}

// WriteHash implements hash.Writer interface.
func (r *ResultRecord) WriteHash(w io.Writer) {
	// hash own fields
	r.RequestRecord.WriteHash(w)
	var data = []interface{}{
		resultRecordID,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// WipeOutRecord is a special record that takes place of another record
// when we need to completely wipe out some information from storage
// (think GDPR).
type WipeOutRecord struct {
	ResultRecord

	Replacement Reference
	WipedHash   Hash
}

// WriteHash implements hash.Writer interface.
func (r *WipeOutRecord) WriteHash(w io.Writer) {
	// hash parent
	r.ResultRecord.WriteHash(w)

	// hash own fields
	r.Replacement.WriteHash(w)
	var data = []interface{}{
		wipeOutRecordID,
		r.WipedHash,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// StatelessResult is a result type that does not need to be stored.
type StatelessResult struct {
	ResultRecord
}

// ReadRecordResult just contains necessary record from storage.
type ReadRecordResult struct {
	StatelessResult

	RecordBody []byte
}

// WriteHash implements hash.Writer interface.
func (r *ReadRecordResult) WriteHash(w io.Writer) {
	// hash parent
	r.StatelessResult.WriteHash(w)

	// hash own fields
	var data = []interface{}{
		readRecordResultID,
		r.RecordBody,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// StatelessCallResult is a contract call result that didn't produce new state.
type StatelessCallResult struct {
	StatelessResult

	ResultMemory Memory
}

// WriteHash implements hash.Writer interface.
func (r *StatelessCallResult) WriteHash(w io.Writer) {
	// hash parent
	r.StatelessResult.WriteHash(w)

	// hash own fields
	var data = []interface{}{
		statelessCallResultID,
		r.ResultMemory,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// Write allows to write to Request's paramMemory.
func (r *StatelessCallResult) Write(p []byte) (n int, err error) {
	r.ResultMemory = make([]byte, len(p))
	return copy(r.ResultMemory, p), nil
}

// Read allows to read Result's resultMemory.
func (r *StatelessCallResult) Read(p []byte) (n int, err error) {
	return copy(p, r.ResultMemory), nil
}

// StatelessExceptionResult is an exception result that does not need to be stored.
type StatelessExceptionResult struct {
	StatelessCallResult

	ExceptionType Reference
}

// WriteHash implements hash.Writer interface.
func (r *StatelessExceptionResult) WriteHash(w io.Writer) {
	// hash parent
	r.StatelessCallResult.WriteHash(w)

	// hash own fields
	r.ExceptionType.WriteHash(w)
	var data = []interface{}{
		statelessExceptionResultID,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// ReadObjectResult contains necessary object's memory.
type ReadObjectResult struct {
	StatelessResult

	State            int32
	MemoryProjection Memory
}

// WriteHash implements hash.Writer interface.
func (r *ReadObjectResult) WriteHash(w io.Writer) {
	// hash parent
	r.StatelessResult.WriteHash(w)

	// hash own fields
	var data = []interface{}{
		readObjectResultID,
		r.State,
		r.MemoryProjection,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// SpecialResult is a result type for special situations.
type SpecialResult struct {
	ResultRecord

	ReasonCode ReasonCode
}

// WriteHash implements hash.Writer interface.
func (r *SpecialResult) WriteHash(w io.Writer) {
	// hash parent
	r.ResultRecord.WriteHash(w)

	// hash own fields
	var data = []interface{}{
		specialResultID,
		r.ReasonCode,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// LockUnlockResult is a result of lock/unlock attempts.
type LockUnlockResult struct {
	SpecialResult
}

// WriteHash implements hash.Writer interface.
func (r *LockUnlockResult) WriteHash(w io.Writer) {
	// hash parent
	r.SpecialResult.WriteHash(w)

	// hash own fields
	var data = []interface{}{
		lockUnlockResultID,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// RejectionResult is a result type for failed attempts.
type RejectionResult struct {
	SpecialResult
}

// WriteHash implements hash.Writer interface.
func (r *RejectionResult) WriteHash(w io.Writer) {
	// hash parent
	r.SpecialResult.WriteHash(w)

	// hash own fields
	var data = []interface{}{
		rejectionResultID,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// StatefulResult is a result type which contents need to be persistently stored.
type StatefulResult struct {
	ResultRecord
}

// ActivationRecord is an activation record.
type ActivationRecord struct {
	StatefulResult

	GoverningDomain Reference
}

// WriteHash implements hash.Writer interface.
func (r *ActivationRecord) WriteHash(w io.Writer) {
	// hash parent
	r.StatefulResult.WriteHash(w)

	// hash own fields
	r.GoverningDomain.WriteHash(w)
	var data = []interface{}{
		activationRecordID,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// ClassActivateRecord is produced when we "activate" new contract class.
type ClassActivateRecord struct {
	ActivationRecord

	CodeRecord    Reference
	DefaultMemory Memory
}

// WriteHash implements hash.Writer interface.
func (r *ClassActivateRecord) WriteHash(w io.Writer) {
	// hash parent
	r.ActivationRecord.WriteHash(w)

	// hash own fields
	r.CodeRecord.WriteHash(w)
	var data = []interface{}{
		classActivateRecordID,
		r.DefaultMemory,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// ObjectActivateRecord is produced when we instantiate new object from an available class.
type ObjectActivateRecord struct {
	ActivationRecord

	ClassActivateRecord Reference
	Memory              Memory
}

// WriteHash implements hash.Writer interface.
func (r *ObjectActivateRecord) WriteHash(w io.Writer) {
	// hash parent
	r.ActivationRecord.WriteHash(w)

	// hash own fields
	r.ClassActivateRecord.WriteHash(w)
	var data = []interface{}{
		objectActivateRecordID,
		r.Memory,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// StorageRecord is produced when we store something in ledger. Code, data etc.
type StorageRecord struct {
	StatefulResult
}

// CodeRecord is a code storage record.
type CodeRecord struct {
	StorageRecord

	Interfaces   []Reference
	TargetedCode [][]byte // []MachineBinaryCode
	SourceCode   string   // ObjectSourceCode
}

// WriteHash implements hash.Writer interface.
func (r *CodeRecord) WriteHash(w io.Writer) {
	// hash parent
	r.StorageRecord.WriteHash(w)

	// hash own fields
	for _, v := range r.Interfaces {
		v.WriteHash(w)
	}
	var data = []interface{}{
		codeRecordID,
		[]byte(r.SourceCode),
	}
	for _, v := range r.TargetedCode {
		data = append(data, interface{}(v))
	}

	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// AmendRecord is produced when we modify another record in ledger.
type AmendRecord struct {
	StatefulResult

	BaseRecord    Reference
	AmendedRecord Reference
}

// WriteHash implements hash.Writer interface.
func (r *AmendRecord) WriteHash(w io.Writer) {
	// hash parent
	r.StatefulResult.WriteHash(w)

	// hash own fields
	r.BaseRecord.WriteHash(w)
	r.AmendedRecord.WriteHash(w)
	var data = []interface{}{
		amendRecordID,
	}

	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// ClassAmendRecord is an amendment record for classes.
type ClassAmendRecord struct {
	AmendRecord

	NewCode []byte // ObjectBinaryCode
}

// WriteHash implements hash.Writer interface.
func (r *ClassAmendRecord) WriteHash(w io.Writer) {
	// hash parent
	r.AmendRecord.WriteHash(w)

	// hash own fields
	var data = []interface{}{
		classAmendRecordID,
		r.NewCode,
	}

	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// MigrationCodes returns a list of data migration procedures for a given code change.
func (r *ClassAmendRecord) MigrationCodes() []*MemoryMigrationCode {
	panic("not implemented")
}

// MemoryMigrationCode is a data migration procedure.
type MemoryMigrationCode struct {
	ClassAmendRecord

	GeneratedByClassRecord Reference
	MigrationCodeRecord    Reference
}

// WriteHash implements hash.Writer interface.
func (r *MemoryMigrationCode) WriteHash(w io.Writer) {
	// hash parent
	r.AmendRecord.WriteHash(w)

	// hash own fields
	r.GeneratedByClassRecord.WriteHash(w)
	r.MigrationCodeRecord.WriteHash(w)
	var data = []interface{}{
		memoryMigrationCodeID,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// DeactivationRecord marks targeted object as disabled.
type DeactivationRecord struct {
	AmendRecord
}

// WriteHash implements hash.Writer interface.
func (r *DeactivationRecord) WriteHash(w io.Writer) {
	// hash parent
	r.AmendRecord.WriteHash(w)

	// hash own fields
	var data = []interface{}{
		deactivationRecordID,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// ObjectAmendRecord is an amendment record for objects.
type ObjectAmendRecord struct {
	AmendRecord

	NewMemory Memory
}

// WriteHash implements hash.Writer interface.
func (r *ObjectAmendRecord) WriteHash(w io.Writer) {
	// hash parent
	r.AmendRecord.WriteHash(w)

	// hash own fields
	var data = []interface{}{
		objectAmendRecordID,
		r.NewMemory,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// StatefulCallResult is a contract call result that produces new state.
type StatefulCallResult struct {
	ObjectAmendRecord

	ResultMemory Memory
}

// WriteHash implements hash.Writer interface.
func (r *StatefulCallResult) WriteHash(w io.Writer) {
	// hash parent
	r.ObjectAmendRecord.WriteHash(w)

	// hash own fields
	var data = []interface{}{
		statefulCallResultID,
		r.ResultMemory,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// StatefulExceptionResult is an exception result that needs to be stored.
type StatefulExceptionResult struct {
	StatefulCallResult

	ExceptionType Reference
}

// WriteHash implements hash.Writer interface.
func (r *StatefulExceptionResult) WriteHash(w io.Writer) {
	// hash parent
	r.StatefulCallResult.WriteHash(w)

	// hash own fields
	r.ExceptionType.WriteHash(w)
	var data = []interface{}{
		statefulExceptionResultID,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// EnforcedObjectAmendRecord is an enforced amendment record for objects.
type EnforcedObjectAmendRecord struct {
	ObjectAmendRecord
}

// WriteHash implements hash.Writer interface.
func (r *EnforcedObjectAmendRecord) WriteHash(w io.Writer) {
	// hash parent
	r.ObjectAmendRecord.WriteHash(w)

	// hash own fields
	var data = []interface{}{
		enforcedObjectAmendRecordID,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// ObjectAppendRecord is an "append state" record for objects. It does not contain full actual state.
type ObjectAppendRecord struct {
	AmendRecord

	AppendMemory Memory
}

// WriteHash implements hash.Writer interface.
func (r *ObjectAppendRecord) WriteHash(w io.Writer) {
	// hash parent
	r.AmendRecord.WriteHash(w)

	// hash own fields
	var data = []interface{}{
		objectAppendRecordID,
		r.AppendMemory,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}
