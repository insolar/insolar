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

package resolver

import (
	"net"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockPacketConn struct {
	mock.Mock
}

func (m *MockPacketConn) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
	args := m.Called(b)
	return args.Int(0), args.Get(1).(net.Addr), args.Error(2)
}

func (m *MockPacketConn) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	args := m.Called(b, addr)
	return args.Int(0), args.Error(1)
}

func (m *MockPacketConn) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPacketConn) LocalAddr() net.Addr {
	args := m.Called()
	return args.Get(0).(net.Addr)
}

func (m *MockPacketConn) SetDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockPacketConn) SetReadDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockPacketConn) SetWriteDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}
