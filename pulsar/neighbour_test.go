package pulsar

import (
	"bytes"
	"net/rpc"
	"os"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/pulsar/pulsartestutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func capture(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

func TestNeighbour_CheckAndRefreshConnection_RefreshSuccess(t *testing.T) {
	client := &pulsartestutil.MockRPCClientWrapper{}
	client.On("CreateConnection", configuration.TCP, "expectedAddress").Return(nil)
	client.On("Lock")
	client.On("Unlock")
	neighbour := &Neighbour{
		ConnectionAddress: "expectedAddress",
		ConnectionType:    configuration.TCP,
		OutgoingClient:    client,
	}

	writtenLog := capture(func() { neighbour.CheckAndRefreshConnection(rpc.ErrShutdown) })

	assert.Contains(t, writtenLog, "Restarting RPC Connection to expectedAddress due to error connection is shut down")
	client.AssertCalled(t, "CreateConnection", configuration.TCP, "expectedAddress")
	client.AssertCalled(t, "Lock")
	client.AssertCalled(t, "Unlock")
}

func TestNeighbour_CheckAndRefreshConnection_RefreshFailed(t *testing.T) {
	client := &pulsartestutil.MockRPCClientWrapper{}
	client.On("CreateConnection", configuration.TCP, "expectedAddress").Return(errors.New("oops"))
	client.On("Lock")
	client.On("Unlock")
	neighbour := &Neighbour{
		ConnectionAddress: "expectedAddress",
		ConnectionType:    configuration.TCP,
		OutgoingClient:    client,
	}

	writtenLog := capture(func() { neighbour.CheckAndRefreshConnection(rpc.ErrShutdown) })

	assert.Contains(t, writtenLog, "Restarting RPC Connection to expectedAddress due to error connection is shut down")
	assert.Contains(t, writtenLog, "Refreshing connection to expectedAddress failed due to error oops")
	client.AssertCalled(t, "CreateConnection", configuration.TCP, "expectedAddress")
	client.AssertCalled(t, "Lock")
	client.AssertCalled(t, "Unlock")
}
