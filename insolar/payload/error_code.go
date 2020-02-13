// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package payload

type ErrorCode uint32

//go:generate stringer -type=ErrorCode

const (
	CodeUnknown ErrorCode = iota
	CodeDeactivated
	CodeFlowCanceled
	CodeNotFound
	CodeNoPendings
	CodeNoStartPulse
	CodeRequestNotFound
	CodeRequestInvalid
	CodeRequestNonClosedOutgoing
	CodeRequestNonOldestMutable
	CodeReasonIsWrong
	CodeNonActivated
	CodeLoopDetected
)

type CodedError struct {
	Text string
	Code ErrorCode
}

func (e *CodedError) GetCode() ErrorCode {
	return e.Code
}

func (e *CodedError) Error() string {
	return e.Text
}

func (i *ErrorCode) Equal(code ErrorCode) bool {
	return *i == code
}
