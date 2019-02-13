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

	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/network/utils"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

type readyEntryImpl struct {
	mutex *sync.Mutex
	conn  net.Conn
}

func newReadyEntryImpl(conn net.Conn) *readyEntryImpl {
	return &readyEntryImpl{
		conn:  conn,
		mutex: &sync.Mutex{},
	}
}

func (e *readyEntryImpl) Open(ctx context.Context) (net.Conn, error) {
	return e.conn, nil
}

func (e *readyEntryImpl) Close() {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	utils.CloseVerbose(e.conn)
}

type managedEntryImpl struct {
	connectionFactory connectionFactory
	address           net.Addr

	mutex *sync.Mutex
	conn  net.Conn
}

func newManagedEntryImpl(connectionFactory connectionFactory, address net.Addr) *managedEntryImpl {
	return &managedEntryImpl{
		connectionFactory: connectionFactory,
		address:           address,
		mutex:             &sync.Mutex{},
	}
}

func (e *managedEntryImpl) Open(ctx context.Context) (net.Conn, error) {
	if e.conn != nil {
		return e.conn, nil
	}

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

func (e *managedEntryImpl) open(ctx context.Context) (net.Conn, error) {
	ctx, span := instracer.StartSpan(ctx, "connectionPool.open")
	span.AddAttributes(
		trace.StringAttribute("create connect to", e.address.String()),
	)
	defer span.End()

	conn, err := e.connectionFactory.CreateConnection(ctx, e.address)
	if err != nil {
		return nil, errors.Wrap(err, "[ Open ] Failed to create TCP connection")
	}

	return conn, nil
}

func (e *managedEntryImpl) Close() {
	if e.conn == nil {
		return
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.conn == nil {
		return
	}

	utils.CloseVerbose(e.conn)
}
