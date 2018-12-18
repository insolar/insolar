package logicrunner

import (
	"context"
	"testing"

	"github.com/insolar/insolar/core/reply"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func TestOnPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	pendingFinishedWasSent := false

	mb := testutils.NewMessageBusMock(t)
	mb.SendMock.Set(func(ctx context.Context, msg core.Message, pulse core.Pulse, opts *core.MessageSendOptions) (core.Reply, error) {
		if msg.Type() == core.TypePendingFinished {
			pendingFinishedWasSent = true
		}
		return &reply.ID{}, nil
	})

	lr, _ := NewLogicRunner(&configuration.LogicRunner{})
	lr.MessageBus = mb

	// test empty lr
	pulse := core.Pulse{}

	err := lr.OnPulse(ctx, pulse)
	require.NoError(t, err)

	objectRef := testutils.RandomRef()

	// test empty es
	lr.state[objectRef] = &ObjectState{ExecutionState: &ExecutionState{Behaviour: &ValidationSaver{}}}
	err = lr.OnPulse(ctx, pulse)
	require.NoError(t, err)
	require.Nil(t, lr.state[objectRef].ExecutionState)
	require.False(t, pendingFinishedWasSent)

	// test empty es with query in current
	lr.state[objectRef] = &ObjectState{
		ExecutionState: &ExecutionState{
			Behaviour: &ValidationSaver{},
			Current:   &CurrentExecution{},
		},
	}
	err = lr.OnPulse(ctx, pulse)
	require.NoError(t, err)
	require.True(t, lr.state[objectRef].ExecutionState.pending)
	lr.ProcessExecutionQueue(ctx, lr.state[objectRef].ExecutionState)
	require.False(t, pendingFinishedWasSent)
	// require.False(t, lr.state[objectRef].ExecutionState.pending) // TODO FIXME probably should pass?

	// test empty es with query in current and query in queue - es.pending true, message.ExecutorResults.Pending = true, message.ExecutorResults.Queue one element
	result := make(chan ExecutionQueueResult, 1)

	// TODO maybe need do something more stable and easy to debug
	go func() {
		<-result
	}()

	qe := ExecutionQueueElement{
		result: result,
	}

	queue := append(make([]ExecutionQueueElement, 0), qe)

	lr.state[objectRef] = &ObjectState{
		ExecutionState: &ExecutionState{
			Behaviour: &ValidationSaver{},
			Current:   &CurrentExecution{},
			Queue:     queue,
		},
	}

	err = lr.OnPulse(ctx, pulse)
	require.NoError(t, err)
	require.True(t, lr.state[objectRef].ExecutionState.pending)
	lr.ProcessExecutionQueue(ctx, lr.state[objectRef].ExecutionState)
	//require.True(t, pendingFinishedWasSent) // TODO FIXME
}
