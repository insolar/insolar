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

package testutil

import (
	"github.com/insolar/insolar/core"
)

// TestMessageRouter can execute messages on LogicRunner.
type TestMessageRouter struct {
	LogicRunner core.LogicRunner
}

// Start is the dummy mock of Start method.
func (*TestMessageRouter) Start(components core.Components) error { return nil }

// Stop is the dummy mock of Stop method.
func (*TestMessageRouter) Stop() error { return nil }

// Route executes message on LogicRunner.
func (r *TestMessageRouter) Route(msg core.Message) (resp core.Response, err error) {
	return r.LogicRunner.Execute(msg)
}

// NewTestMessageRouter creates TestMessageRouter which mocks the real one.
func NewTestMessageRouter(lr core.LogicRunner) *TestMessageRouter {
	return &TestMessageRouter{LogicRunner: lr}
}
