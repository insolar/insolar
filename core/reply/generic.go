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

import "github.com/insolar/insolar/core"

// OK is a generic reply for success calls without returned value.
type OK struct {
}

// Type implementation of Reply interface.
func (e *OK) Type() core.ReplyType {
	return TypeOK
}

// Error is common reaction for methods returning id to lifeline states.
type Error struct {
	ErrType ErrType
}

// Type implementation of Reply interface.
func (e *Error) Type() core.ReplyType {
	return TypeError
}

// Error returns concrete error for stored type.
func (e *Error) Error() error {
	switch e.ErrType {
	case ErrDeactivated:
		return core.ErrDeactivated
	}
	return core.ErrUnknown
}
