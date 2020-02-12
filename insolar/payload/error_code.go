// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
