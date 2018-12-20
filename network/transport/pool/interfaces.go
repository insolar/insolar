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
	"context"
	"net"
)

type ConnectionPool interface {
	GetConnection(ctx context.Context, address net.Addr) (net.Conn, error)
	Reset()
}

type connectionFactory interface {
	CreateConnection(ctx context.Context, address net.Addr) (net.Conn, error)
}

type iterateFunc func(conn net.Conn)

type unsafeConnectionHolder interface {
	Get(address net.Addr) (net.Conn, bool)
	Delete(address net.Addr)
	Add(address net.Addr, conn net.Conn)
	Size() int
	Clear()
	Iterate(iterateFunc iterateFunc)
}

func newUnsafeConnectionHolder() unsafeConnectionHolder {
	return newUnsafeConnectionHolderImpl()
}

func NewConnectionPool(connectionFactory connectionFactory) ConnectionPool {
	return newConnectionPool(connectionFactory)
}
