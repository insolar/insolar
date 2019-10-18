//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package smachine

import (
	"fmt"
	"runtime/debug"
)

var _ error = SlotPanicError{}

type SlotPanicError struct {
	Msg       string
	Recovered interface{}
	Prev      error
	Stack     []byte
	IsAsync   bool
}

func (e SlotPanicError) Error() string {
	sep := ""
	if len(e.Stack) > 0 {
		sep = "\n"
	}
	if e.Prev != nil {
		return fmt.Sprintf("%s: %v%s%s\nCaused by:\n%s", e.Msg, e.Recovered, sep, string(e.Stack), e.Prev.Error())
	}
	return fmt.Sprintf("%s: %v%s%s", e.Msg, e.Recovered, sep, string(e.Stack))
}

func (e SlotPanicError) String() string {
	return fmt.Sprintf("%s: %v", e.Msg, e.Recovered)
}

func RecoverSlotPanic(msg string, recovered interface{}, prev error) error {
	if recovered == nil {
		return prev
	}
	return SlotPanicError{Msg: msg, Recovered: recovered, Prev: prev}
}

func RecoverSlotPanicWithStack(msg string, recovered interface{}, prev error) error {
	if recovered == nil {
		return prev
	}
	return SlotPanicError{Msg: msg, Recovered: recovered, Prev: prev, Stack: debug.Stack()}
}

func RecoverAsyncSlotPanicWithStack(msg string, recovered interface{}, prev error) error {
	if recovered == nil {
		return prev
	}
	return SlotPanicError{Msg: msg, Recovered: recovered, Prev: prev, Stack: debug.Stack(), IsAsync: true}
}
