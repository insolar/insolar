package transport

import (
	"context"
	"net"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/instrumentation/inslogger"
)

type tcpConnectionFactory struct{}

func (*tcpConnectionFactory) CreateConnection(ctx context.Context, address string) (net.Conn, error) {
	logger := inslogger.FromContext(ctx)
	tcpAddress, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, errors.New("[ createConnection ] Failed to get tcp address")
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddress)
	if err != nil {
		logger.Errorf("[ createConnection ] Failed to open connection to %s: %s", address, err.Error())
		return nil, errors.Wrap(err, "[ createConnection ] Failed to open connection")
	}

	err = conn.SetKeepAlive(true)
	if err != nil {
		logger.Error("[ createConnection ] Failed to set keep alive")
	}

	err = conn.SetNoDelay(true)
	if err != nil {
		logger.Error("[ createConnection ] Failed to set connection no delay: ", err.Error())
	}

	return conn, nil
}
