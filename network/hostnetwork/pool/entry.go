//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package pool

import (
	"context"
	"io"
	"sync"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/network/utils"
)

type entryImpl struct {
	connectionFactory connectionFactory
	address           string
	onClose           onClose

	mutex *sync.Mutex

	conn io.ReadWriteCloser
}

func newEntryImpl(connectionFactory connectionFactory, address string, onClose onClose) *entryImpl {
	return &entryImpl{
		connectionFactory: connectionFactory,
		address:           address,
		mutex:             &sync.Mutex{},
		onClose:           onClose,
	}
}

func (e *entryImpl) Open(ctx context.Context) (io.ReadWriteCloser, error) {
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

func (e *entryImpl) open(ctx context.Context) (io.ReadWriteCloser, error) {
	logger := inslogger.FromContext(ctx)
	ctx, span := instracer.StartSpan(ctx, "connectionPool.open")
	span.AddAttributes(
		trace.StringAttribute("create connect to", e.address),
	)
	defer span.End()

	conn, err := e.connectionFactory.CreateConnection(ctx, e.address)
	if err != nil {
		return nil, errors.Wrap(err, "[ Open ] Failed to create TCP connection")
	}

	go func(e *entryImpl, conn io.ReadWriteCloser) {
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
