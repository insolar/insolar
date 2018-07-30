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

// Package record contains various record type definitions.
package record

// ProjectionType is a "view filter" for record.
// E.g. we can read whole object or just it's hash.
type ProjectionType uint

// Memory is actual contracts' state, variables etc.
type Memory []byte

// Record is base interface for all records.
type Record interface {
	Hash() Hash
	TimeSlot() uint64
	Type() TypeID
}

// Reference is a pointer that allows to address any record across whole network.
// TODO: Should implement normal Reference type (not interface)
// TODO: Globally unique record identifier must be found
type Reference interface {
	Record
}

// AppDataRecord is persistent data record stored in ledger.
type AppDataRecord struct {
	timeSlotNo uint64
	recType    TypeID
}

// Hash returns SHA-3 hash sum of Record
func (r *AppDataRecord) Hash() Hash {
	panic("implement me")
}

// TimeSlot returns time slot number that Record belongs to.
func (r *AppDataRecord) TimeSlot() uint64 {
	return r.timeSlotNo
}

// Type returns Record type.
func (r *AppDataRecord) Type() TypeID {
	return r.recType
}
