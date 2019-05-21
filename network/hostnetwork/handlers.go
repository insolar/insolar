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

package hostnetwork

import (
	"context"
	"io"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/hostnetwork/future"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/pool"
)

// RequestHandler is callback function for request handling
type RequestHandler func(p *packet.PacketBackend)

// StreamHandler parses packets from data stream and calls request handler or response handler
type StreamHandler struct {
	requestHandler  RequestHandler
	responseHandler future.PacketHandler
}

// NewStreamHandler creates new StreamHandler
func NewStreamHandler(requestHandler RequestHandler, responseHandler future.PacketHandler) *StreamHandler {
	return &StreamHandler{
		requestHandler:  requestHandler,
		responseHandler: responseHandler,
	}
}

func (s *StreamHandler) HandleStream(address string, reader io.ReadWriteCloser) {
	for {
		p, err := packet.DeserializePacketBackend(reader)

		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				log.Info("[ HandleStream ] Connection closed by peer")
				return
			}

			log.Error("[ HandleStream ] Failed to deserialize packet: ", err.Error())
		} else {
			ctx, logger := inslogger.WithTraceField(context.Background(), p.TraceID)
			logger.Debug("[ HandleStream ] Handling packet RequestID = ", p.RequestID)

			if p.IsResponse() {
				go s.responseHandler.Handle(ctx, p)
			} else {
				go s.requestHandler(p)
			}
		}
	}
}

// SendPacket sends packet using connection from pool
func SendPacket(ctx context.Context, pool pool.ConnectionPool, p *packet.PacketBackend) error {
	data, err := packet.SerializePacketBackend(p)
	if err != nil {
		return errors.Wrap(err, "Failed to serialize packet")
	}

	conn, err := pool.GetConnection(ctx, p.Receiver)
	if err != nil {
		return errors.Wrap(err, "Failed to get connection")
	}

	n, err := conn.Write(data)
	if err != nil {
		// retry
		pool.CloseConnection(ctx, p.Receiver)
		conn, err = pool.GetConnection(ctx, p.Receiver)

		if err != nil {
			return errors.Wrap(err, "[ SendBuffer ] Failed to get connection")
		}
		n, err = conn.Write(data)
	}
	if err == nil {
		metrics.NetworkSentSize.Add(float64(n))
		return nil
	}
	return errors.Wrap(err, "[ send ] Failed to write data")
}
