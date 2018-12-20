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
	"sync"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/utils"
	"github.com/pkg/errors"
)

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

func (cp *connectionPool) GetConnection(ctx context.Context, address net.Addr) (net.Conn, error) {
	logger := inslogger.FromContext(ctx)

	conn, ok := cp.getConnection(address)

	logger.Debugf("[ GetConnection ] Finding connection to %s in pool: %s", address, ok)

	if ok && !connectionClosedByPeer(conn) {
		return conn, nil
	}

	logger.Debugf("[ GetConnection ] Missing open connection to %s in pool ", address)

	return cp.getOrCreateConnection(ctx, address)
}

func (cp *connectionPool) getConnection(address net.Addr) (net.Conn, bool) {
	cp.mutex.RLock()
	defer cp.mutex.RUnlock()

	return cp.unsafeConnectionsHolder.Get(address)
}

func (cp *connectionPool) getOrCreateConnection(ctx context.Context, address net.Addr) (net.Conn, error) {
	logger := inslogger.FromContext(ctx)

	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	conn, ok := cp.unsafeConnectionsHolder.Get(address)
	logger.Debugf("[ getOrCreateConnection ] Finding connection to %s in pool: %s", address, ok)

	if ok {
		if !connectionClosedByPeer(conn) {
			return conn, nil
		}

		logger.Debugf("[ getOrCreateConnection ] Connection to %s closed by peer, closing it on our side", address)
		utils.CloseVerbose(conn)

		logger.Debugf("[ getOrCreateConnection ] Delete connection to %s from pool: %s", address)
		cp.unsafeConnectionsHolder.Delete(address)
	}

	logger.Debugf("[ getOrCreateConnection ] Failed to retrieve connection to %s, creating it", address)

	conn, err := cp.connectionFactory.CreateConnection(ctx, address)
	if err != nil {
		return nil, errors.Wrap(err, "[ send ] Failed to create TCP connection")
	}

	cp.unsafeConnectionsHolder.Add(address, conn)
	logger.Debugf(
		"[ getOrCreateConnection ] Added connection to %s. Current pool size: %d",
		conn.RemoteAddr(),
		cp.unsafeConnectionsHolder.Size(),
	)

	return conn, nil
}

func (cp *connectionPool) Reset() {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	cp.unsafeConnectionsHolder.Iterate(func(conn net.Conn) {
		utils.CloseVerbose(conn)
	})
	cp.unsafeConnectionsHolder.Clear()
}
