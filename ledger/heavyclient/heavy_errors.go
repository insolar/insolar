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

package heavyclient

import (
	"github.com/insolar/insolar/core/reply"
)

// HeavyErr holds core.Reply and implements core.Retryable and error interfaces.
type HeavyErr struct {
	reply *reply.HeavyError
	err   error
}

// Error implements error interface.
func (he HeavyErr) Error() string {
	if he.err != nil {
		return he.err.Error()
	}
	if he.reply != nil {
		return he.reply.Error()
	}
	panic("neither reply or error defined in HeavyErr")
}

// IsRetryable checks retryability of message.
func (he HeavyErr) IsRetryable() bool {
	if he.reply == nil {
		return false
	}
	return he.reply.IsRetryable()
}
