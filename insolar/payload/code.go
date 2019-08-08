//
// Copyright 2019 Insolar Technologies GmbH
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
//

package payload

const (
	CodeUnknown            = 1000
	CodeDeactivated        = 1001
	CodeFlowCanceled       = 1002
	CodeNotFound           = 1003
	CodeNoPendings         = 1004
	CodeNoStartPulse       = 1005
	ReasonNotFound         = 1006
	ReasonIsWrong          = 1007
	IncomingRequestIsWrong = 1008
	RequestNotFound        = 1009
)

type ErrorCoder interface {
	Error() string
	GetErrorCode() uint32
}

type CodedError struct {
	ErrorText string
	ErrorCode uint32
}

func (e *CodedError) GetErrorCode() uint32 {
	return e.ErrorCode
}
func (e *CodedError) Error() string {
	return e.ErrorText
}
