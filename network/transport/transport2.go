package transport

import (
	"context"
	"net"
	"strings"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/transport/pool"
)

type TcpTransport struct {
	listenAddress string
	processor     StreamProcessor
	pool          pool.ConnectionPool
}

func NewTcpTransport(cfg configuration.Configuration) *TcpTransport {
	return &TcpTransport{
		listenAddress: cfg.Host.Transport.Address,
		pool:          pool.NewConnectionPool(&tcpConnectionFactory{}),
	}
}

func (t *TcpTransport) SendBuffer(ctx context.Context, address string, buff []byte) error {
	conn, err := t.pool.GetConnection(ctx, address)
	if err != nil {
		return err
	}
	_, err = conn.Write(buff)
	return err
}

func (t *TcpTransport) SendDgram(ctx context.Context, address string, buff []byte) error {
	panic("tcp can't send dgram")
}

func (t *TcpTransport) Start(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	logger.Info("[ Start ] Start TCP transport")

	listener, err := net.Listen("tcp", t.listenAddress)
	if err != nil {
		logger.Info("[ Start ] Failed to prepare TCP transport")
		return err
	}

	go func() {
		for {
			// TODO handle Stop
			conn, err := listener.Accept()
			if err != nil {
				//<-t.disconnectFinished
				if strings.Contains(strings.ToLower(err.Error()), "use of closed network connection") {
					logger.Info("Connection closed, quiting accept loop")
					return
				}

				logger.Error("[ Start ] Failed to accept connection: ", err.Error())
				return
			}

			logger.Debugf("[ Start ] Accepted new connection from %s", conn.RemoteAddr())

			go t.handleAcceptedConnection(conn)
		}

	}()

	return nil
}

func (t *TcpTransport) Stop(ctx context.Context) error {
	return nil
}

func (t *TcpTransport) handleAcceptedConnection(conn net.Conn) {
	//defer utils.CloseVerbose(conn)

	err := t.processor.ProcessStream(conn.RemoteAddr().String(), conn)
	if err != nil {
		inslogger.FromContext(context.Background()).Errorf("Failed to process stream from %s: %s", conn.RemoteAddr().String(), err.Error())
	}
	// for {
	// //lengthBytes := make([]byte, 8)
	//
	// r := bufio.NewReader(conn)
	//
	// lengthBytes, err := r.Peek(16)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(lengthBytes)
	// }
	// bn := r.Buffered()
	// fmt.Println(bn)
	// r.Discard(16)
	//
	// lengthBytes, err = r.Peek(10)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	fmt.Println(lengthBytes)
	// }
	//msg, err := t.serializer.DeserializePacket(conn)
	// if err != nil {
	// 	if err == io.EOF || err == io.ErrUnexpectedEOF {
	// 		log.Warn("[ handleAcceptedConnection ] Connection closed by peer")
	// 		return
	// 	}
	//
	// 	log.Error("[ handleAcceptedConnection ] Failed to deserialize packet: ", err.Error())
	// } else {
	// 	ctx, logger := inslogger.WithTraceField(context.Background(), msg.TraceID)
	// 	logger.Debug("[ handleAcceptedConnection ] Handling packet: ", msg.RequestID)
	//
	// 	go t.packetHandler.Handle(ctx, msg)
	//}

}
