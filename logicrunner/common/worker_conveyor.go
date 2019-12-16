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

package common

import (
	"sync"
	"time"

	"github.com/insolar/insolar/conveyor"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/log/logcommon"
)

type ConveyorWorker struct {
	conveyor *conveyor.PulseConveyor
	stopped  sync.WaitGroup
}

func (w *ConveyorWorker) Stop() {
	w.conveyor.StopNoWait()
	w.stopped.Wait()
}

func (w *ConveyorWorker) AttachTo(conveyor *conveyor.PulseConveyor) {
	if conveyor == nil {
		panic("illegal value")
	}
	if w.conveyor != nil {
		panic("illegal state")
	}
	w.conveyor = conveyor
	w.stopped.Add(1)
	conveyor.StartWorker(nil, func() {
		w.stopped.Done()
	})
}

func NewConveyorWorker() ConveyorWorker {
	return ConveyorWorker{}
}

type AsyncTimeMessage struct {
	*logcommon.LogObjectTemplate `txt:"async time"`

	AsyncComponent     string `opt:""`
	AsyncExecutionTime int64
}

func LogAsyncTime(log smachine.Logger, timeBefore time.Time, component string) {
	log.Trace(AsyncTimeMessage{
		AsyncComponent:     component,
		AsyncExecutionTime: time.Now().Sub(timeBefore).Nanoseconds(),
	})
}
