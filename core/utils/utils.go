package utils

import (
	"encoding/binary"

	"github.com/satori/go.uuid"
)

// RandTraceID returns random traceID in uuid format
func RandTraceID() string {
	qid, err := uuid.NewV4()
	if err != nil {
		return "createRandomTraceIDFailed:" + err.Error()
	}
	return qid.String()
}

func UInt32ToBytes(n uint32) []byte {
	buff := make([]byte, 4)
	binary.BigEndian.PutUint32(buff, n)
	return buff
}
