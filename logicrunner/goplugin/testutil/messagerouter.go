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

// TestMessageBus can execute messages on LogicRunner.
type TestMessageBus struct {
	LogicRunner core.LogicRunner
}

func (eb *TestMessageBus) Register(p core.MessageType, handler core.MessageHandler) error {
	return nil
}

// Start is the dummy mock of Start method.
func (*TestMessageBus) Start(components core.Components) error { return nil }

// Stop is the dummy mock of Stop method.
func (*TestMessageBus) Stop() error { return nil }

// Send executes message on LogicRunner.
func (eb *TestMessageBus) Send(msg core.Message) (resp core.Reply, err error) {
	return eb.LogicRunner.Execute(msg.(core.LogicRunnerEvent))
}

func (*TestMessageBus) SendAsync(msg core.Message) {}

// NewTestMessageBus creates TestMessageBus which mocks the real one.
func NewTestMessageBus(lr core.LogicRunner) *TestMessageBus {
	return &TestMessageBus{LogicRunner: lr}
}
