/*
 *    Copyright 2019 Insolar
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

package conveyour_getobject

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/ledger/conveyor"
)

// Pure handler.
func fetchJet(helper conveyor.Helper, item conveyor.SlotItem, payload interface{}) conveyor.StateID {
	helper.SendTask(&JetAdapterTask{
		Object: *item.Event().DefaultTarget().Record(),
	})
	return StateReceiveJet
}

// Pure handler.
func receiveJet(helper conveyor.Helper, item conveyor.SlotItem, payload interface{}, result JetAdapterResult) conveyor.StateID {
	nodeForJet := func(id core.RecordID) core.RecordID { return core.RecordID{} }
	me := func() core.RecordID { return core.RecordID{} }

	if nodeForJet(result.JetID) != me() {
		return StateError
	}
	item.SetPayload(&JetPayload{
		JetID: result.JetID,
	})

	return StateSelectHandler
}

// Generated raw handler.
func selectHandler(helper conveyor.Helper, item conveyor.SlotItem, payload JetPayload) conveyor.StateID {
	if payload.Err != nil {
		return StateError
	}
	switch item.Event().(type) {
	case *message.GetObject:
		return StateGetObject
		// the rest of the handlers.
	}

	return StateError
}

func getObject(helper conveyor.Helper, item conveyor.SlotItem, payload JetPayload) conveyor.StateID {
	helper.SendTask(&GetObjectTask{
		Object: *item.Event().DefaultTarget().Record(),
		JetID:  payload.JetID,
	})

	return StateReturnObject
}

func returnObject(helper conveyor.Helper, item conveyor.SlotItem, payload GetObjectPayload) core.Reply {
	if payload.Err != nil {
		return &reply.Error{}
	}

	return &reply.Object{Memory: payload.Memory}
}

func returnError(helper conveyor.Helper, item conveyor.SlotItem) core.Reply {
	return &reply.Error{}
}

func init() {
	conveyor.RegisterReply(StateError, returnError)
	conveyor.RegisterActive(StateGetObject, getObject)
	conveyor.RegisterReply(StateReturnObject, returnObject)
	conveyor.RegisterActive(StateFetchJet, fetchJet)
	conveyor.RegisterActive(StateSelectHandler, selectHandler)
	conveyor.RegisterInactive(StateReceiveJet, receiveJet)
}
