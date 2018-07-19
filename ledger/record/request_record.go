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

type RequestRecord struct {
	AppDataRecord

	requestor RecordReference
	target    RecordReference
}

func (r *RequestRecord) Requestor() RecordReference {
	return r.requestor
}

func (r *RequestRecord) Target() RecordReference {
	return r.target
}

type CallRequest struct {
	RequestRecord

	callInterface       RecordReference
	callMethodSignature uint
	paramMemory         Memory
}

func (r *CallRequest) CallInterface() RecordReference {
	return r.callInterface
}

func (r *CallRequest) CallMethod() uint {
	return r.callMethodSignature
}

func (r *CallRequest) Read(p []byte) (n int, err error) {
	return copy(p, r.paramMemory), nil
}

func (r *CallRequest) Write(p []byte) (n int, err error) {
	r.paramMemory = make([]byte, len(p))
	return copy(r.paramMemory, p), nil
}

type LockUnlockRequest struct {
	RequestRecord

	transaction          RecordReference
	expectedLockDuration time.Duration
}

func (r *LockUnlockRequest) Transaction() RecordReference {
	return r.transaction
}

func (r *LockUnlockRequest) ExpectedLockDuration() time.Duration {
	return r.expectedLockDuration
}

type ReadRequest struct {
	RequestRecord
}

type ReadRecordRequest struct {
	ReadRequest

	expectedRecordType RecordType
}

func (r *ReadRecordRequest) ExpectedRecordType() RecordType {
	return r.expectedRecordType
}

type ReadObject struct {
	ReadRequest

	projectionType ProjectionType
}

func (r *ReadObject) ProjectionType() ProjectionType {
	return r.projectionType
}

type ReadObjectComposite struct {
	ReadObject

	compositeType RecordReference
}

func (r *ReadObjectComposite) CompositeType() RecordReference {
	return r.compositeType
}
