/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package pool

import (
	"context"
	"net"
)

type ConnectionPool interface {
	GetConnection(ctx context.Context, address net.Addr) (bool, net.Conn, error)
	RegisterConnection(ctx context.Context, address net.Addr, conn net.Conn) bool
	CloseConnection(ctx context.Context, address net.Addr)
	Reset(ctx context.Context)
}

type connectionFactory interface {
	CreateConnection(ctx context.Context, address net.Addr) (net.Conn, error)
}

type entry interface {
	Open(ctx context.Context) (net.Conn, error)
	Close()
}

func newManagerEntry(connectionFactory connectionFactory, address net.Addr) entry {
	return newManagedEntryImpl(connectionFactory, address)
}

func newReadyEntry(conn net.Conn) entry {
	return newReadyEntryImpl(conn)
}

type iterateFunc func(entry entry)

type entryHolder interface {
	Get(address net.Addr) (entry, bool)
	Delete(address net.Addr)
	Add(address net.Addr, entry entry)
	Size() int
	Clear()
	Iterate(iterateFunc iterateFunc)
}

func newEntryHolder() entryHolder {
	return newEntryHolderImpl()
}

func NewConnectionPool(connectionFactory connectionFactory) ConnectionPool {
	return newConnectionPool(connectionFactory)
}
