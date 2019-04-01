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

package slot

import (
	"github.com/insolar/insolar/conveyor/queue"
)

// TaskPusher is interface which permits only safe access to slot
//go:generate minimock -i github.com/insolar/insolar/conveyor/slot.TaskPusher -o ./ -s _mock.go
type TaskPusher interface {
	SinkPush(data interface{}) error
	SinkPushAll(data []interface{}) error
	PushSignal(signalType uint32, callback queue.SyncDone) error
}
