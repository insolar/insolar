package record

import (
	"crypto/sha256"
)

const (
	HashSize = sha256.Size

	RecTypeCodeBase = RecordType(iota + 1)
	RecTypeCodeAmendment
	RecTypeObjectInstance
	RecTypeObjectAmendment
	RecTypeObjectData
)

type RecordHash [HashSize]byte

type RecordType uint

type ProjectionType uint

type Memory []byte

type Record interface {
	Hash() RecordHash
	TimeSlot() uint64
	Type() RecordType
}

// TODO: Should implement normal RecordReference type (not interface)
// TODO: Globally unique record identifier must be found
type RecordReference interface {
	Record
}

type AppDataRecord struct {
	timeSlotNo     uint64
	recType        RecordType
	collisionNonce uint // TODO: combine nonce with type?
}

func (r *AppDataRecord) Hash() RecordHash {
	panic("implement me")
}

func (r *AppDataRecord) TimeSlot() uint64 {
	return r.timeSlotNo
}

func (r *AppDataRecord) Type() RecordType {
	return r.recType
}
