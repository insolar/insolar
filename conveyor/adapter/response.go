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

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

// NewResponseSendAdapter creates new instance of adapter for sending response
func NewResponseSendAdapter() PulseConveyorAdapterTaskSink {
	return NewAdapterWithQueue(NewResponseSender())
}

// ResponseSenderTask is task for adapter for sending response
type ResponseSenderTask struct {
	Future core.ConveyorFuture
	Result core.Reply
}

// ResponseSender is worker for adapter for sending response
type ResponseSender struct{}

// NewResponseSender returns new instance of worker which sending response
func NewResponseSender() Processor {
	return &ResponseSender{}
}

// Process implements Processor interface
func (rs *ResponseSender) Process(adapterID uint32, task AdapterTask, cancelInfo CancelInfo) Events {
	payload, ok := task.TaskPayload.(ResponseSenderTask)
	var msg interface{}

	if !ok {
		msg = errors.Errorf("[ ResponseSender.Process ] Incorrect payload type: %T", task.TaskPayload)
		return Events{RespPayload: msg}
	}

	done := make(chan bool, 1)
	go func(payload ResponseSenderTask) {
		res := payload.Result
		f := payload.Future
		f.SetResult(res)
		done <- true
	}(payload)

	select {
	case <-cancelInfo.Cancel():
		log.Info("[ ResponseSender.Process ] Cancel. Return Nil as Response")
		msg = nil
	case <-cancelInfo.Flush():
		log.Info("[ ResponseSender.Process ] Flush. DON'T Return Response")
		return Events{Flushed: true}
	case <-done:
		msg = fmt.Sprintf("Response was send successfully")
	}

	log.Info("[ ResponseSender.Process ] response message is", msg)
	return Events{RespPayload: msg}
}
