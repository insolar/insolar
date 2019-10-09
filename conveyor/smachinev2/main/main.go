///
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
///

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/insolar/insolar/conveyor/smachine/tools"
	smachine "github.com/insolar/insolar/conveyor/smachinev2"
	"github.com/insolar/insolar/conveyor/smachinev2/main/example"
)

func main() {
	//var fn1 func(example.ServiceA, string) string
	//fn1 = example.ServiceA.DoSomething
	//fn1(&implA{}, "test")

	sm := smachine.NewSlotMachine(smachine.SlotMachineConfig{
		SyncStrategy:    &syncStrategy{},
		SlotPageSize:    1000,
		PollingPeriod:   1000 * time.Millisecond,
		PollingTruncate: 100 * time.Millisecond,
		ScanCountLimit:  1000,
	}, nil)

	//example.SetInjectServiceAdapterA(&implA{}, &sm)

	sm.AddNew(context.Background(), smachine.NoLink(), &example.StateMachine1{})

	signal := tools.NewVersionedSignal()
	worker := example.NewSimpleSlotWorker(signal.Mark())

	prev := 0
	for i := 0; ; i++ {
		repeatNow, nextPollTime := sm.ScanOnce(worker)
		//fmt.Printf("%03d %v ================================== slots=%v of %v\n", i, time.Now(), sm.OccupiedSlotCount(), sm.AllocatedSlotCount())
		//fmt.Printf("%03d %v =============== repeatNow=%v nextPollTime=%v\n", i, time.Now(), repeatNow, nextPollTime)

		if i >= prev+10 {
			prev = i
			fmt.Printf("%03d %v ================================== slots=%v of %v\n", i, time.Now(), sm.OccupiedSlotCount(), sm.AllocatedSlotCount())
			sm.Cleanup(worker)
			fmt.Printf("%03d %v ================================== slots=%v of %v\n", i, time.Now(), sm.OccupiedSlotCount(), sm.AllocatedSlotCount())
		}

		if repeatNow {
			continue
		}
		if !nextPollTime.IsZero() {
			time.Sleep(time.Until(nextPollTime))
		} else {
			time.Sleep(3 * time.Second)
		}
	}
}

type implA struct {
}

func (*implA) DoSomething(param string) string {
	return param
}

func (*implA) DoSomethingElse(param0 string, param1 int) (bool, string) {
	return param1 != 0, param0
}
