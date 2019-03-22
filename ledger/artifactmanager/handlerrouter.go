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

package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/insolar"
)

// Handler is an alias for wrapping func-s
type Handler func(insolar.MessageHandler) insolar.MessageHandler

// BuildMiddleware return wrapping insolar.MessageHandler
// If want call wrapHandler1(wrapHandler2(wrapHandler3(handler))),
// we should use Build(handler, wrapHandler1, wrapHandler2, wrapHandler3).
func BuildMiddleware(handler insolar.MessageHandler, wrapHandlers ...Handler) insolar.MessageHandler {
	result := handler

	for i := range wrapHandlers {
		result = wrapHandlers[len(wrapHandlers)-1-i](result)
	}

	return result
}

// PreSender is an alias for a function
// which is working like a `middleware` for messagebus.Send
type PreSender func(Sender) Sender

// Sender is an alias for signature of messagebus.Send
type Sender func(context.Context, insolar.Message, *insolar.MessageSendOptions) (insolar.Reply, error)

// BuildSender allows us to build a chain of PreSender before calling Sender
// The main idea of it is ability to make a different things before sending message
// For example we can cache some replies. Another example is the sendAndFollow redirect method
func BuildSender(sender Sender, preSenders ...PreSender) Sender {
	result := sender

	for i := range preSenders {
		result = preSenders[len(preSenders)-1-i](result)
	}

	return result
}
