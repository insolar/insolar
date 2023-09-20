package hostnetwork

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/testutils"
)

func TestNewStreamHandler(t *testing.T) {
	defer testutils.LeakTester(t)

	requestHandler := func(ctx context.Context, p *packet.ReceivedPacket) {
		inslogger.FromContext(ctx).Info("requestHandler")
	}

	h := NewStreamHandler(requestHandler, nil)

	con1, _ := net.Pipe()

	done := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		h.HandleStream(ctx, "127.0.0.1:8080", con1)
		done <- struct{}{}
	}()

	cancel()
	// con2.Close()

	select {
	case <-done:
		return
	case <-time.After(time.Second * 5):
		t.Fail()
	}
}
