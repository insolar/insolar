/*
 *    Copyright 2018 Insolar
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

package reply

import (
	"fmt"

	"github.com/insolar/insolar/core"
)

// ErrHeavySyncInProgress returned when heavy sync in progress.
const (
	ErrHeavySyncInProgress ErrType = iota + 1
)

// HeavyError carries heavy sync error information.
type HeavyError struct {
	Message  string
	SubType  ErrType
	JetID    core.RecordID
	PulseNum core.PulseNumber
}

// Type implementation of Reply interface.
func (e *HeavyError) Type() core.ReplyType {
	return TypeHeavyError
}

// ConcreteType returns concrete error type.
func (e *HeavyError) ConcreteType() ErrType {
	return e.SubType
}

// Error returns error message for stored type.
func (e *HeavyError) Error() string {
	return fmt.Sprintf("%v (JetID=%v, PulseNum=%v)", e.Message, e.JetID, e.PulseNum)
}

// IsRetryable returns true if retry could be performed.
func (e *HeavyError) IsRetryable() bool {
	return e.SubType == ErrHeavySyncInProgress
}
