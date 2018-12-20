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

package pool

import (
	"net"
)

type unsafeConnectionsHolderImpl struct {
	connections map[string]net.Conn
}

func newUnsafeConnectionHolderImpl() unsafeConnectionHolder {
	return &unsafeConnectionsHolderImpl{
		connections: make(map[string]net.Conn),
	}
}

func (uch *unsafeConnectionsHolderImpl) key(address net.Addr) string {
	return address.String()
}

func (uch *unsafeConnectionsHolderImpl) Get(address net.Addr) (net.Conn, bool) {
	conn, ok := uch.connections[uch.key(address)]

	return conn, ok
}

func (uch *unsafeConnectionsHolderImpl) Delete(address net.Addr) {
	delete(uch.connections, uch.key(address))
}

func (uch *unsafeConnectionsHolderImpl) Add(address net.Addr, conn net.Conn) {
	uch.connections[uch.key(address)] = conn
}

func (uch *unsafeConnectionsHolderImpl) Clear() {
	for key := range uch.connections {
		delete(uch.connections, key)
	}
}

func (uch *unsafeConnectionsHolderImpl) Iterate(iterateFunc iterateFunc) {
	for _, conn := range uch.connections {
		iterateFunc(conn)
	}
}

func (uch *unsafeConnectionsHolderImpl) Size() int {
	return len(uch.connections)
}
