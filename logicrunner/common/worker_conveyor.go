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
	"github.com/insolar/insolar/conveyor"
	"github.com/insolar/insolar/conveyor/sworker"
)

type ConveyorWorker struct {
	stoppedChan chan struct{}
}

func (w *ConveyorWorker) Stop() {
	close(w.stoppedChan)
}

func (w *ConveyorWorker) AttachTo(
	conveyor *conveyor.PulseConveyor,
) {
	workerFactory := sworker.NewAttachableSimpleSlotWorker()
	go conveyor.RunOnWorker(workerFactory, w.stoppedChan)
}

func NewConveyorWorker() *ConveyorWorker {
	return &ConveyorWorker{
		stoppedChan: make(chan struct{}),
	}
}
