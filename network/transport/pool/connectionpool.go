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
	"github.com/insolar/insolar/network/utils"
	"github.com/pkg/errors"
)

type lockableConnection struct {
	net.Conn
	sync.Locker
}

func (lc *lockableConnection) Write(data []byte) (int, error) {
	lc.Lock()
	defer lc.Unlock()

	// TODO: sergey.morozov 16.01.19: possible malformed packet fix; uncomment this when you meet errors ;)
	// var written int
	// for written < len(data) {
	// 	n, err := lc.Conn.Write(data[written:])
	// 	written += n
	// 	if err != nil {
	// 		return written, err
	// 	}
	// }
	// return written, nil

	return lc.Conn.Write(data)
}

type connectionPool struct {
	connectionFactory connectionFactory

	unsafeConnectionsHolder unsafeConnectionHolder
	mutex                   sync.RWMutex
}

func newConnectionPool(connectionFactory connectionFactory) *connectionPool {
	return &connectionPool{
		connectionFactory: connectionFactory,

		unsafeConnectionsHolder: newUnsafeConnectionHolder(),
	}
}

func (cp *connectionPool) GetConnection(ctx context.Context, address net.Addr) (bool, net.Conn, error) {
	logger := inslogger.FromContext(ctx)

	conn, ok := cp.getConnection(address)

	logger.Debugf("[ GetConnection ] Finding connection to %s in pool: %t", address, ok)

	if ok {
		return false, conn, nil
	}

	logger.Debugf("[ GetConnection ] Missing open connection to %s in pool ", address)

	return cp.getOrCreateConnection(ctx, address)
}

func (cp *connectionPool) RegisterConnection(ctx context.Context, address net.Addr, conn net.Conn) bool {
	logger := inslogger.FromContext(ctx)

	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	_, ok := cp.unsafeConnectionsHolder.Get(address)

	logger.Debugf("[ RegisterConnection ] Finding connection to %s in pool: %s", address, ok)

	if ok {
		return false
	}

	logger.Debugf("[ RegisterConnection ] Missing open connection to %s in pool ", address)

	cp.unsafeConnectionsHolder.Add(address, conn)
	logger.Debugf(
		"[ RegisterConnection ] Added connection to %s. Current pool size: %d",
		conn.RemoteAddr(),
		cp.unsafeConnectionsHolder.Size(),
	)
	return true
}

func (cp *connectionPool) CloseConnection(ctx context.Context, address net.Addr) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	logger := inslogger.FromContext(ctx)

	conn, ok := cp.unsafeConnectionsHolder.Get(address)
	logger.Debugf("[ CloseConnection ] Finding connection to %s in pool: %s", address, ok)

	if ok {
		utils.CloseVerbose(conn)

		logger.Debugf(
			"[ CloseConnection ] Delete connection to %s. Current pool size: %d",
			address,
			cp.unsafeConnectionsHolder.Size(),
		)
		cp.unsafeConnectionsHolder.Delete(address)
		metrics.NetworkConnections.Set(float64(cp.unsafeConnectionsHolder.Size()))
	}
}

func (cp *connectionPool) getConnection(address net.Addr) (net.Conn, bool) {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()

	return cp.unsafeConnectionsHolder.Get(address)
}

func (cp *connectionPool) getOrCreateConnection(ctx context.Context, address net.Addr) (bool, net.Conn, error) {
	logger := inslogger.FromContext(ctx)

	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	conn, ok := cp.unsafeConnectionsHolder.Get(address)
	logger.Debugf("[ getOrCreateConnection ] Finding connection to %s in pool: %s", address, ok)

	if ok {
		return false, conn, nil
	}

	logger.Debugf("[ getOrCreateConnection ] Failed to retrieve connection to %s, creating it", address)

	conn, err := cp.connectionFactory.CreateConnection(ctx, address)
	if err != nil {
		return false, nil, errors.Wrap(err, "[ send ] Failed to create TCP connection")
	}


	lc := &lockableConnection{
		Conn:   conn,
		Locker: &sync.Mutex{},
	}

	cp.unsafeConnectionsHolder.Add(address, lc)
	size := cp.unsafeConnectionsHolder.Size()
	logger.Debugf(
		"[ getOrCreateConnection ] Added connection to %s. Current pool size: %d",
		conn.RemoteAddr(),
		size,
	)
	metrics.NetworkConnections.Set(float64(size))

	return true, conn, nil
}

func (cp *connectionPool) Reset(ctx context.Context) {
	logger := inslogger.FromContext(ctx)

	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	logger.Debugf("[ Reset ] Reset pool of size: %d", cp.unsafeConnectionsHolder.Size())

	cp.unsafeConnectionsHolder.Iterate(func(conn net.Conn) {
		utils.CloseVerbose(conn)
	})
	cp.unsafeConnectionsHolder.Clear()
	metrics.NetworkConnections.Set(float64(cp.unsafeConnectionsHolder.Size()))
}
