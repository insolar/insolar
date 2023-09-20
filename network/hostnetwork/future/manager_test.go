package future

import (
	"testing"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	m := NewManager()

	require.IsType(t, m, &futureManager{})
}

func TestFutureManager_Create(t *testing.T) {
	m := NewManager()

	sender, _ := host.NewHostN("127.0.0.1:31337", gen.Reference())
	receiver, _ := host.NewHostN("127.0.0.2:31338", gen.Reference())

	p := packet.NewPacket(sender, receiver, types.Unknown, 123)
	future := m.Create(p)

	require.EqualValues(t, future.ID(), p.RequestID)
	require.Equal(t, future.Request(), p)
	require.Equal(t, future.Receiver(), receiver)
}

func TestFutureManager_Get(t *testing.T) {
	m := NewManager()

	sender, _ := host.NewHostN("127.0.0.1:31337", gen.Reference())
	receiver, _ := host.NewHostN("127.0.0.2:31338", gen.Reference())

	p := packet.NewPacket(sender, receiver, types.Unknown, 123)

	require.Nil(t, m.Get(p))

	expectedFuture := m.Create(p)
	actualFuture := m.Get(p)

	require.Equal(t, expectedFuture, actualFuture)
}

func TestFutureManager_Canceler(t *testing.T) {
	m := NewManager()

	sender, _ := host.NewHostN("127.0.0.1:31337", gen.Reference())
	receiver, _ := host.NewHostN("127.0.0.2:31338", gen.Reference())

	p := packet.NewPacket(sender, receiver, types.Unknown, 123)

	future := m.Create(p)
	require.NotNil(t, future)

	future.Cancel()

	require.Nil(t, m.Get(p))
}
