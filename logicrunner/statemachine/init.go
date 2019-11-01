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

package statemachine

import (
	"fmt"

	"github.com/insolar/insolar/conveyor"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/logicrunner/sm_execute_request"
	"github.com/insolar/insolar/logicrunner/sm_request"
	"github.com/insolar/insolar/pulse"
)

func DefaultHandlersFactory(_ pulse.Number, input conveyor.InputEvent) smachine.CreateFunc {
	switch inputConverted := input.(type) {
	case *payload.Meta:
		return sm_request.HandlerFactoryMeta(inputConverted)
	case *sm_execute_request.SMEventSendOutgoing:
		return sm_execute_request.HandlerFactoryOutgoingSender(inputConverted)
	default:
		panic(fmt.Sprintf("unknown event type, got %T", input))
	}
}
