package logicrunner

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewQueueIsEmtpy(t *testing.T) {
	t.Parallel()
	q := NewGenericQueue()
	require.True(t, q.Empty())
}

func TestQueueEnqueueDequeue(t *testing.T) {
	type Payload struct {
		payload string
	}
	t.Parallel()
	q := NewGenericQueue()
	inMsg := Payload{payload: "hello"}
	require.True(t, q.Empty())
	q.Enqueue(inMsg)
	require.False(t, q.Empty())
	outMsg := q.Dequeue()

	require.True(t, q.Empty())
	require.Equal(t, inMsg, outMsg.(Payload))
}

func TestQueueOrdering(t *testing.T) {
	type Payload struct {
		payload string
	}
	t.Parallel()
	q := NewGenericQueue()
	msg1 := Payload{payload: "msg1"}
	msg2 := Payload{payload: "msg2"}
	msg3 := Payload{payload: "msg3"}
	q.Enqueue(msg1)
	q.Enqueue(msg2)
	q.Enqueue(msg3)
	require.Equal(t, msg1, q.Dequeue())
	require.Equal(t, msg2, q.Dequeue())
	require.Equal(t, msg3, q.Dequeue())
	require.True(t, q.Empty())
}
