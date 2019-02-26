/*
 *    Copyright 2019 Insolar Technologies
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

package core

import (
	"fmt"
)

// TerminationHandler handles such node events as graceful stop, abort, etc.
//go:generate minimock -i github.com/insolar/insolar/core.TerminationHandler -o ../testutils -s _mock.go
type TerminationHandler interface {
	// Abort forces to stop all node components
	Abort(reason string)
}

type terminationHandler struct{}

func (terminationHandler) Abort(reason string) {
	panic(fmt.Sprintf("Node leave acknowledged by network. Goodbye! Reason: %s", reason))
}

func NewTerminationHandler() TerminationHandler {
	return &terminationHandler{}
}
