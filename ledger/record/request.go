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

	Requester Reference
	Target    Reference
}

// CallRequest is a contract execution request.
// Implements io.ReadWriter interface.
type CallRequest struct {
	RequestRecord

	CallInterface       Reference
	CallMethodSignature uint32
	ParamMemory         Memory
}

// CallMethod is a contract method number to call.
func (r *CallRequest) CallMethod() uint32 {
	return r.CallMethodSignature
}

// Read allows to read Request's paramMemory.
func (r *CallRequest) Read(p []byte) (n int, err error) {
	return copy(p, r.ParamMemory), nil
}

// Write allows to write to Request's paramMemory.
func (r *CallRequest) Write(p []byte) (n int, err error) {
	r.ParamMemory = make([]byte, len(p))
	return copy(r.ParamMemory, p), nil
}

// LockUnlockRequest is a request to temporary lock (or unlock) another record.
type LockUnlockRequest struct {
	RequestRecord

	Transaction          Reference
	ExpectedLockDuration time.Duration
}

// ReadRequest is a request type to read data.
type ReadRequest struct {
	RequestRecord
}

// ReadRecordRequest is a request type to read another record.
type ReadRecordRequest struct {
	ReadRequest

	ExpectedRecordType TypeID
}

// ReadObject is a request type
type ReadObject struct {
	ReadRequest

	ProjectionType ProjectionType
}

// ReadObjectComposite is a request to read object including it's "injected" fields.
type ReadObjectComposite struct {
	ReadObject

	CompositeType Reference
}
