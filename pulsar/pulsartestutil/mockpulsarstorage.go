/*
 *    Copyright 2018 Insolar
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

// Package pulsartestutil - test utils for pulsar package
package pulsartestutil

import (
	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/mock"
)

// MockPulsarStorage mocks PulsarStorage interface
type MockPulsarStorage struct {
	mock.Mock
}

func (mock *MockPulsarStorage) GetLastPulse() (*core.Pulse, error) {
	args := mock.Called()
	return args.Get(0).(*core.Pulse), args.Error(1)
}

func (*MockPulsarStorage) SetLastPulse(pulse *core.Pulse) error {
	return nil
}

func (*MockPulsarStorage) SavePulse(pulse *core.Pulse) error {
	return nil
}

func (*MockPulsarStorage) Close() error {
	panic("implement me")
}
