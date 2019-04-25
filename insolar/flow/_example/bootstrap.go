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

package _example

import (
	"bytes"
	"context"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/flow/handler"
)

type GetObject struct {
	DBConnection *bytes.Buffer
}

func (s *GetObject) Past(context.Context, flow.Flow) error    { /* ... */ return nil }
func (s *GetObject) Present(context.Context, flow.Flow) error { /* ... */ return nil }
func (s *GetObject) Future(context.Context, flow.Flow) error  { /* ... */ return nil }

func bootstrapExample() { // nolint
	DBConnection := bytes.NewBuffer(nil)

	hand := handler.NewHandler(
		// These functions can provide any variables via closure.
		// IMPORTANT: they must create NEW handle instances on every call.
		func(msg bus.Message) flow.Handle {
			s := GetObject{
				DBConnection: DBConnection,
			}
			return s.Present
		},
	)

	// Use handler to handle incoming messages.
	_ = hand
}
