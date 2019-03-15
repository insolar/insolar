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
)

type SendResponseTask interface {
	GetFuture() core.Future
	GetResult() core.Reply
}

func NewSendResponseAdapter() PulseConveyorAdapterTaskSink {
	return NewAdapterWithQueue(NewSendResponseProcessing())
}

func NewSendResponseProcessing() TaskProcessing {
	return &SendResponseProcessing{}
}

type SendResponseProcessing struct{}

func (sr *SendResponseProcessing) Process(adapterID uint32, task AdapterTask, cancelInfo *cancelInfoT) {
	done := make(chan bool, 1)
	go func() {
		payload := task.taskPayload.(SendResponseTask)
		res := payload.GetResult()
		f := payload.GetFuture()
		f.SetResult(res)
		done <- true
	}()

	var msg string
	select {
	case <-cancelInfo.cancel:
		msg = "Cancel. Return Response"
	case <-cancelInfo.flush:
		log.Info("[ SimpleWaitAdapter.doWork ] Flush. DON'T Return Response")
		return
	case <-done:
		msg = fmt.Sprintf("Response was send successfully")
	}

	log.Info("[ SimpleWaitAdapter.doWork ] ", msg)

	task.respSink.PushResponse(adapterID, task.elementID, task.handlerID, msg)
}
