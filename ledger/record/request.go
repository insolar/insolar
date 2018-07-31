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
	"time"
)

// RequestRecord is common type for all requests.
type RequestRecord struct {
	AppDataRecord

	Requester Reference
	Target    Reference
}

// WriteHash implements hash.Writer interface.
func (r *RequestRecord) WriteHash(w io.Writer) {
	// hash own fields
	r.Requester.WriteHash(w)
	r.Target.WriteHash(w)
	err := binary.Write(w, binary.BigEndian, requestRecordID)
	if err != nil {
		panic("binary.Write failed:" + err.Error())
	}
}

// CallRequest is a contract execution request.
// Implements io.ReadWriter interface.
type CallRequest struct {
	RequestRecord

	CallInterface       Reference
	CallMethodSignature uint32
	ParamMemory         Memory
}

// WriteHash implements hash.Writer interface.
func (r *CallRequest) WriteHash(w io.Writer) {
	// hash parent
	r.RequestRecord.WriteHash(w)

	// hash own fields
	r.CallInterface.WriteHash(w)
	var data = []interface{}{
		callRequestID,
		r.CallMethodSignature,
		r.ParamMemory,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
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

// WriteHash implements hash.Writer interface.
func (r *LockUnlockRequest) WriteHash(w io.Writer) {
	// hash parent
	r.RequestRecord.WriteHash(w)

	// hash own fields
	r.Transaction.WriteHash(w)
	var data = []interface{}{
		lockUnlockRequestID,
		r.ExpectedLockDuration,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
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

// WriteHash implements hash.Writer interface.
func (r *ReadRecordRequest) WriteHash(w io.Writer) {
	// hash parent
	r.ReadRequest.WriteHash(w)
	// hash own fields
	var data = []interface{}{
		readRecordRequestID,
		r.ExpectedRecordType,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// ReadObject is a request type
type ReadObject struct {
	ReadRequest

	ProjectionType ProjectionType
}

// WriteHash implements hash.Writer interface.
func (r *ReadObject) WriteHash(w io.Writer) {
	// hash parent
	r.ReadRequest.WriteHash(w)

	// hash own fields
	var data = []interface{}{
		readObjectID,
		r.ProjectionType,
	}
	for _, v := range data {
		err := binary.Write(w, binary.BigEndian, v)
		if err != nil {
			panic("binary.Write failed:" + err.Error())
		}
	}
}

// ReadObjectComposite is a request to read object including it's "injected" fields.
type ReadObjectComposite struct {
	ReadObject

	CompositeType Reference
}

// WriteHash implements hash.Writer interface.
func (r *ReadObjectComposite) WriteHash(w io.Writer) {
	// hash parent
	r.ReadObject.WriteHash(w)
	// hash own fields
	r.CompositeType.WriteHash(w)
}
