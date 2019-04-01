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

	"github.com/insolar/insolar/conveyor/queue"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
)

// NodeStateTask is task for adapter for getting node state
type NodeStateTask struct {
	Callback queue.SyncDone
	Pulse    insolar.Pulse
}

// NodeStateProcessor is worker for adapter for getting node state
type NodeStateProcessor struct {
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme `inject:""`
}

// NewNodeStateProcessor returns new instance of worker which getting node state
func NewNodeStateProcessor() Processor {
	return &NodeStateProcessor{}
}

// Process implements Processor interface
func (rs *NodeStateProcessor) Process(task AdapterTask, nestedEventHelper NestedEventHelper, ci CancelInfo) interface{} {
	payload, ok := task.TaskPayload.(NodeStateTask)

	if !ok {
		log.Errorf("[ NodeStateProcessor.Process ] Incorrect payload type: %T", task.TaskPayload)
		return nil
	}

	log.Errorf(">>>>>>>>>>>>>: %+v", rs.PlatformCryptographyScheme)
	log.Errorf(">>>>>>>>>>>>>++++: %+v", rs.PlatformCryptographyScheme.IntegrityHasher())

	// TODO: calculate node state hash with info about pulse, for now - just return 1,2,3
	res := rs.PlatformCryptographyScheme.IntegrityHasher().Hash([]byte{1, 2, 3})

	c := payload.Callback
	c.SetResult(res)

	log.Info("[ NodeStateProcessor.Process ] NodeState is", res)
	return fmt.Sprintf("NodeState was calculated successfully")
}
