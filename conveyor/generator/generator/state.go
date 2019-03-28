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

type PulseState uint

const (
	Future = PulseState(iota)
	Present
	Past
)

type state struct {
	handlers [3]map[handlerType]*handler
}

func (s *state) GetTransition() *handler {
	return s.handlers[Present][Transition]
}

func (s *state) GetTransitionFuture() *handler {
	return s.handlers[Future][Transition]
}

func (s *state) GetTransitionPast() *handler {
	if s.handlers[Past][Transition] != nil {
		return s.handlers[Past][Transition]
	}
	return s.handlers[Present][Transition]
}

func (s *state) GetMigration() *handler {
	return s.handlers[Present][Migration]
}

func (s *state) GetMigrationFuturePresent() *handler {
	return s.handlers[Future][Migration]
}


func (s *state) GetAdapterResponse() *handler {
	return s.handlers[Present][AdapterResponse]
}

func (s *state) GetAdapterResponseFuture() *handler {
	return s.handlers[Future][AdapterResponse]
}

func (s *state) GetAdapterResponsePast() *handler {
	if s.handlers[Past][AdapterResponse] != nil {
		return s.handlers[Past][AdapterResponse]
	}
	return s.handlers[Present][AdapterResponse]
}
