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
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MocktConn struct {
	mock.Mock
}

func (m *MocktConn) ReadFrom(b []byte) (n int, addr net.Addr, err error) {
	args := m.Called(b)
	return args.Int(0), args.Get(1).(net.Addr), args.Error(2)
}

func (m *MocktConn) WriteTo(b []byte, addr net.Addr) (n int, err error) {
	//args := m.Called(b, addr)
	//return args.Int(0), args.Error(1)
	return 0, nil
}

func (m *MocktConn) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MocktConn) LocalAddr() net.Addr {
	args := m.Called()
	return args.Get(0).(net.Addr)
}

func (m *MocktConn) SetDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MocktConn) SetReadDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MocktConn) SetWriteDeadline(t time.Time) error {
	args := m.Called(t)
	return args.Error(0)
}

func TestNewStunResolver(t *testing.T) {
	stunAddr := "127.0.0.1:31337"
	resolver := NewStunResolver(stunAddr)

	require.IsType(t, &stunResolver{}, resolver)
}

func TestStunResolver_Resolve(t *testing.T) {
	stunAddr := "127.0.0.1:31337"
	resolver := NewStunResolver(stunAddr)
	require.IsType(t, &stunResolver{}, resolver)

	conn := &MocktConn{}
	conn.On("LocalAddr").Return(net.ResolveUDPAddr("udp", stunAddr))

	resolver.Resolve(conn)
}
