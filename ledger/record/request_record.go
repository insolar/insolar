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
