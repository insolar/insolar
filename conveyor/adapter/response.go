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

package adapter

import (
	"fmt"

	"github.com/insolar/insolar/conveyor/interfaces/slot"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

// NewResponseSendAdapter creates new instance of adapter for sending response
func NewResponseSendAdapter() PulseConveyorAdapterTaskSink {
	return NewAdapterWithQueue(NewSendResponseProcessor())
}

// SendResponseTask is task for adapter for sending response
type SendResponseTask struct {
	Future core.ConveyorFuture
	Result core.Reply
}

// SendResponseProcessor is worker for adapter for sending response
type SendResponseProcessor struct{}

// NewResponseSender returns new instance of worker which sending response
func NewSendResponseProcessor() Processor {
	return &SendResponseProcessor{}
}

// Process implements Processor interface
func (rs *SendResponseProcessor) Process(task AdapterTask, nestedEventHelper NestedEventHelper, cancelInfo CancelInfo) interface{} {
	payload, ok := task.TaskPayload.(SendResponseTask)
	var msg interface{}

	if !ok {
		msg = errors.Errorf("[ SendResponseProcessor.Process ] Incorrect payload type: %T", task.TaskPayload)
		return msg
	}

	res := payload.Result
	f := payload.Future
	f.SetResult(res)

	msg = fmt.Sprintf("Response was send successfully")
	log.Info("[ SendResponseProcessor.Process ] response message is", msg)
	return msg
}

// ResponseSenderHelper is helper for ResponseSender
type ResponseSenderHelper struct{}

// SendResponse makes correct message and send it to adapter
func (r *ResponseSenderHelper) SendResponse(element slot.SlotElementHelper, result core.Reply, respHandlerID uint32) error {

	pendingMsg, ok := element.GetInputEvent().(core.ConveyorPendingMessage)
	if !ok {
		return errors.Errorf("[ ResponseSenderHelper.SendResponse ] Input event is not core.ConveyorPendingMessage: %T", element.GetInputEvent())
	}

	response := SendResponseTask{
		Future: pendingMsg.Future,
		Result: result,
	}
	err := element.SendTask(uint32(ResponseSenderAdapterID), response, respHandlerID)
	return errors.Wrap(err, "[ ResponseSenderHelper.SendResponse ] Can't SendTask")
}
