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

package pool

import (
	"context"
	"net"
	"sync"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/network/utils"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

type entryImpl struct {
	connectionFactory connectionFactory
	address           net.Addr
	onClose           onClose

	mutex *sync.Mutex

	conn net.Conn
}

func newEntryImpl(connectionFactory connectionFactory, address net.Addr, onClose onClose) *entryImpl {
	return &entryImpl{
		connectionFactory: connectionFactory,
		address:           address,
		mutex:             &sync.Mutex{},
		onClose:           onClose,
	}
}

func (e *entryImpl) Open(ctx context.Context) (net.Conn, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.conn != nil {
		return e.conn, nil
	}

	conn, err := e.open(ctx)
	if err != nil {
		return nil, err
	}

	e.conn = conn
	return e.conn, nil
}

func (e *entryImpl) open(ctx context.Context) (net.Conn, error) {
	logger := inslogger.FromContext(ctx)
	ctx, span := instracer.StartSpan(ctx, "connectionPool.open")
	span.AddAttributes(
		trace.StringAttribute("create connect to", e.address.String()),
	)
	defer span.End()

	conn, err := e.connectionFactory.CreateConnection(ctx, e.address)
	if err != nil {
		return nil, errors.Wrap(err, "[ Open ] Failed to create TCP connection")
	}

	go func(e *entryImpl, conn net.Conn) {
		b := make([]byte, 1)
		_, err := conn.Read(b)
		if err != nil {
			logger.Infof("[ Open ] remote host 'closed' connection to %s: %s", e.address, err)
			e.onClose(ctx, e.address)
			return
		}

		logger.Errorf("[ Open ] unexpected data on connection to %s", e.address)
	}(e, conn)

	return conn, nil
}

func (e *entryImpl) Close() {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.conn != nil {
		utils.CloseVerbose(e.conn)
	}
}
