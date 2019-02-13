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
	"sync"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/metrics"
)

type connectionPool struct {
	connectionFactory connectionFactory

	entryHolder entryHolder
	mutex       sync.RWMutex
}

func newConnectionPool(connectionFactory connectionFactory) *connectionPool {
	return &connectionPool{
		connectionFactory: connectionFactory,

		entryHolder: newEntryHolder(),
	}
}

func (cp *connectionPool) GetConnection(ctx context.Context, address net.Addr) (bool, net.Conn, error) {
	logger := inslogger.FromContext(ctx)

	entry, ok := cp.getEntry(address)

	logger.Debugf("[ GetConnection ] Finding entry for connection to %s in pool: %t", address, ok)

	if ok {
		conn, err := entry.Open(ctx)
		return false, conn, err
	}

	logger.Debugf("[ GetConnection ] Missing entry for connection to %s in pool ", address)
	created, entry := cp.getOrCreateEntry(ctx, address)
	conn, err := entry.Open(ctx)

	return created, conn, err
}

func (cp *connectionPool) RegisterConnection(ctx context.Context, address net.Addr, conn net.Conn) bool {
	logger := inslogger.FromContext(ctx)

	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	_, ok := cp.entryHolder.Get(address)

	logger.Debugf("[ RegisterConnection ] Finding entry for connection to %s in pool: %s", address, ok)

	if ok {
		return false
	}

	logger.Debugf("[ RegisterConnection ] Missing entry for connection to %s in pool ", address)

	cp.entryHolder.Add(address, newReadyEntry(conn))
	logger.Debugf(
		"[ RegisterConnection ] Added entry for connection to %s. Current pool size: %d",
		conn.RemoteAddr(),
		cp.entryHolder.Size(),
	)
	return true
}

func (cp *connectionPool) CloseConnection(ctx context.Context, address net.Addr) {
	logger := inslogger.FromContext(ctx)

	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	entry, ok := cp.entryHolder.Get(address)
	logger.Debugf("[ CloseConnection ] Finding entry for connection to %s in pool: %s", address, ok)

	if ok {
		entry.Close()

		logger.Debugf(
			"[ CloseConnection ] Delete connection to %s. Current pool size: %d",
			address,
			cp.entryHolder.Size(),
		)
		cp.entryHolder.Delete(address)
		metrics.NetworkConnections.Dec()
	}
}

func (cp *connectionPool) getEntry(address net.Addr) (entry, bool) {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()

	return cp.entryHolder.Get(address)
}

func (cp *connectionPool) getOrCreateEntry(ctx context.Context, address net.Addr) (bool, entry) {
	logger := inslogger.FromContext(ctx)

	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	entry, ok := cp.entryHolder.Get(address)
	logger.Debugf("[ getOrCreateEntry ] Finding entry for connection to %s in pool: %s", address, ok)

	if ok {
		return false, entry
	}

	logger.Debugf("[ getOrCreateEntry ] Failed to retrieve entry for connection to %s, creating it", address)

	entry = newManagerEntry(cp.connectionFactory, address)

	cp.entryHolder.Add(address, entry)
	size := cp.entryHolder.Size()
	logger.Debugf(
		"[ getOrCreateEntry ] Added entry for connection to %s. Current pool size: %d",
		address.String(),
		size,
	)
	metrics.NetworkConnections.Inc()

	return true, entry
}

func (cp *connectionPool) Reset(ctx context.Context) {
	logger := inslogger.FromContext(ctx)

	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	logger.Debugf("[ Reset ] Reset pool of size: %d", cp.entryHolder.Size())

	cp.entryHolder.Iterate(func(entry entry) {
		entry.Close()
	})
	cp.entryHolder.Clear()
	metrics.NetworkConnections.Set(float64(cp.entryHolder.Size()))
}
