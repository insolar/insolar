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

package generator

import (
	"github.com/insolar/insolar/conveyor/interfaces/constant"
)

type state struct {
	handlers				   map[constant.PulseState]map[handlerType]*handler
}

func (s *state) GetTransition() *handler {
	return s.handlers[constant.Present][Transition]
}

func (s *state) GetTransitionFuture() *handler {
	return s.handlers[constant.Future][Transition]
}

func (s *state) GetTransitionPast() *handler {
	if s.handlers[constant.Past][Transition] != nil {
		return s.handlers[constant.Past][Transition]
	}
	return s.handlers[constant.Present][Transition]
}

func (s *state) GetMigration() *handler {
	return s.handlers[constant.Present][Migration]
}

func (s *state) GetMigrationFuturePresent() *handler {
	return s.handlers[constant.Future][Migration]
}


func (s *state) GetAdapterResponse() *handler {
	return s.handlers[constant.Present][AdapterResponse]
}

func (s *state) GetAdapterResponseFuture() *handler {
	return s.handlers[constant.Future][AdapterResponse]
}

func (s *state) GetAdapterResponsePast() *handler {
	if s.handlers[constant.Past][AdapterResponse] != nil {
		return s.handlers[constant.Past][AdapterResponse]
	}
	return s.handlers[constant.Present][AdapterResponse]
}
