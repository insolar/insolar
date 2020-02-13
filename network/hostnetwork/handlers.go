// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package hostnetwork

import (
	"context"
	"io"

	"github.com/insolar/insolar/network"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/hostnetwork/future"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/pool"
)

// RequestHandler is callback function for request handling
type RequestHandler func(ctx context.Context, p *packet.ReceivedPacket)

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

func (s *StreamHandler) HandleStream(ctx context.Context, address string, reader io.ReadWriteCloser) {
	mainLogger := inslogger.FromContext(ctx)

	logLevel := inslogger.GetLoggerLevel(ctx)
	// get only log level from context, discard TraceID in favor of packet TraceID
	packetCtx := inslogger.WithLoggerLevel(context.Background(), logLevel)

	closer := make(chan struct{})
	go func() {
		select {
		// transport is stopping
		case <-ctx.Done():
		// stream end by remote end
		case <-closer:
		}

		network.CloseVerbose(reader)
	}()

	for {
		p, length, err := packet.DeserializePacket(mainLogger, reader)

		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				mainLogger.Debug("[ HandleStream ] Connection closed by peer")
				close(closer)
				return
			}

			if network.IsConnectionClosed(err) || network.IsClosedPipe(err) {
				mainLogger.Info("[ HandleStream ] Connection closed.")
				return
			}

			mainLogger.Warnf("[ HandleStream ] Failed to deserialize packet: ", err.Error())
			return
		}

		packetCtx, logger := inslogger.WithTraceField(packetCtx, p.TraceID)
		span, err := instracer.Deserialize(p.TraceSpanData)
		if err == nil {
			packetCtx = instracer.WithParentSpan(packetCtx, span)
		} else {
			inslogger.FromContext(packetCtx).Warn("Incoming packet without span")
		}
		logger.Debugf("[ HandleStream ] Handling packet RequestID = %d, size = %d", p.RequestID, length)
		metrics.NetworkRecvSize.Observe(float64(length))
		if p.IsResponse() {
			go s.responseHandler.Handle(packetCtx, p)
		} else {
			go s.requestHandler(packetCtx, p)
		}

	}
}

// SendPacket sends packet using connection from pool
func SendPacket(ctx context.Context, pool pool.ConnectionPool, p *packet.Packet) error {
	data, err := packet.SerializePacket(p)
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
		inslogger.FromContext(ctx).Warn("[ SendPacket ] retry conn.Write")
		pool.CloseConnection(ctx, p.Receiver)
		conn, err = pool.GetConnection(ctx, p.Receiver)

		if err != nil {
			return errors.Wrap(err, "[ SendPacket ] Failed to get connection")
		}
		n, err = conn.Write(data)
	}
	if err == nil {
		metrics.NetworkSentSize.Observe(float64(n))
		return nil
	}
	return errors.Wrap(err, "[ SendPacket ] Failed to write data")
}
