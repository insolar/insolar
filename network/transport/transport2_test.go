package transport

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransport(t *testing.T) {
	t.Skip("wf")
	tcp := TcpTransport{listenAddress: "127.0.0.1:8080"}
	ctx := context.Background()
	err := tcp.Start(ctx)
	defer tcp.Stop(ctx)
	assert.NoError(t, err)

	_, err = http.Get("http://127.0.0.1:8080")
	assert.NoError(t, err)

	<-time.After(12 * time.Hour)
}
