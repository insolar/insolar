package utils

import "github.com/satori/go.uuid"

// RandTraceID returns random traceID in uuid format
func RandTraceID() string {
	qid, err := uuid.NewV4()
	if err != nil {
		return "createRandomTraceIDFailed:" + err.Error()
	}
	return qid.String()
}
