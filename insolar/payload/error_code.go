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
	CodeUnknown = 0

	CodeDeactivated  = 1001
	CodeFlowCanceled = 1002
	CodeNotFound     = 1003
	CodeNoPendings   = 1004
	CodeNoStartPulse = 1005

	CodeRequestNotFound = 1006
	CodeInvalidRequest  = 1007
	CodeReasonNotFound  = 1008
	CodeReasonIsWrong   = 1009
)

type CodedError struct {
	Text string
	Code uint32
}

func (e *CodedError) GetCode() uint32 {
	return e.Code
}
func (e *CodedError) Error() string {
	return e.Text
}
