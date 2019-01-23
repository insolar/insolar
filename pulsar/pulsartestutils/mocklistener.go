/*
 *    Copyright 2019 Insolar
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
package pulsartestutils

import (
	"net"

	"github.com/stretchr/testify/mock"
)

// MockListener mocks net.Listener interface
type MockListener struct {
	mock.Mock
}

func (mock *MockListener) Accept() (net.Conn, error) {
	panic("implement me")
}

func (mock *MockListener) Close() error {
	panic("implement me")
}

func (mock *MockListener) Addr() net.Addr {
	panic("implement me")
}
