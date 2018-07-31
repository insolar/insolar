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

import "time"

// RequestRecord is common type for all requests.
type RequestRecord struct {
	AppDataRecord

	requester Reference
	target    Reference
}

// Requester is a request author's global address.
func (r *RequestRecord) Requester() Reference {
	return r.requester
}

// Target is an address of contract that we want to execute, data that we want to get etc.
func (r *RequestRecord) Target() Reference {
	return r.target
}

// CallRequest is a contract execution request.
// Implements io.ReadWriter interface.
type CallRequest struct {
	RequestRecord

	callInterface       Reference
	callMethodSignature uint
	paramMemory         Memory
}

// CallInterface is a call interface address.
func (r *CallRequest) CallInterface() Reference {
	return r.callInterface
}

// CallMethod is a contract method number to call.
func (r *CallRequest) CallMethod() uint {
	return r.callMethodSignature
}

// Read allows to read Request's paramMemory.
func (r *CallRequest) Read(p []byte) (n int, err error) {
	return copy(p, r.paramMemory), nil
}

// Write allows to write to Request's paramMemory.
func (r *CallRequest) Write(p []byte) (n int, err error) {
	r.paramMemory = make([]byte, len(p))
	return copy(r.paramMemory, p), nil
}

// LockUnlockRequest is a request to temporary lock (or unlock) another record.
type LockUnlockRequest struct {
	RequestRecord

	transaction          Reference
	expectedLockDuration time.Duration
}

// Transaction is a Reference to Transaction record.
func (r *LockUnlockRequest) Transaction() Reference {
	return r.transaction
}

// ExpectedLockDuration is expected time duration that record will be locked.
func (r *LockUnlockRequest) ExpectedLockDuration() time.Duration {
	return r.expectedLockDuration
}

// ReadRequest is a request type to read data.
type ReadRequest struct {
	RequestRecord
}

// ReadRecordRequest is a request type to read another record.
type ReadRecordRequest struct {
	ReadRequest

	expectedRecordType TypeID
}

// ExpectedRecordType is an expected Type of target record.
func (r *ReadRecordRequest) ExpectedRecordType() TypeID {
	return r.expectedRecordType
}

// ReadObject is a request type
type ReadObject struct {
	ReadRequest

	projectionType ProjectionType
}

// ProjectionType is a "view filter" for record.
// E.g. we can read whole object or just it's hash.
func (r *ReadObject) ProjectionType() ProjectionType {
	return r.projectionType
}

// ReadObjectComposite is a request to read object including it's "injected" fields.
type ReadObjectComposite struct {
	ReadObject

	compositeType Reference
}

// CompositeType is reference to a Record describing composition type.
func (r *ReadObjectComposite) CompositeType() Reference {
	return r.compositeType
}
