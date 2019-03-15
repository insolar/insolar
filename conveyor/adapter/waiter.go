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
	"time"

	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

// NewWaitAdapter creates new instance of SimpleWaitAdapter with Waiter as worker
func NewWaitAdapter() PulseConveyorAdapterTaskSink {
	return NewAdapterWithQueue(NewWaiter())
}

type WaiterTask struct {
	waitPeriodMilliseconds int
}

type Waiter struct{}

func NewWaiter() Worker {
	return &Waiter{}
}

func (w *Waiter) Process(adapterID uint32, task AdapterTask, cancelInfo *cancelInfoT) {
	log.Info("[ doWork ] Start. cancelInfo.id: ", cancelInfo.id)

	payload, ok := task.taskPayload.(WaiterTask)
	var msg interface{}

	if !ok {
		msg = errors.Errorf("[ PushTask ] Incorrect payload type: %T", task.taskPayload)
		task.respSink.PushResponse(adapterID, task.elementID, task.handlerID, msg)
		return
	}

	select {
	case <-cancelInfo.cancel:
		log.Info("[ SimpleWaitAdapter.doWork ] Cancel. Return Nil as Response")
		msg = nil
	case <-cancelInfo.flush:
		log.Info("[ SimpleWaitAdapter.doWork ] Flush. DON'T Return Response")
		return
	case <-time.After(time.Duration(payload.waitPeriodMilliseconds) * time.Millisecond):
		msg = fmt.Sprintf("Work completed successfully. Waited %d millisecond", payload.waitPeriodMilliseconds)
	}

	log.Info("[ SimpleWaitAdapter.doWork ] ", msg)

	task.respSink.PushResponse(adapterID,
		task.elementID,
		task.handlerID,
		msg)

	// TODO: remove cancelInfo from swa.taskHolder
}
