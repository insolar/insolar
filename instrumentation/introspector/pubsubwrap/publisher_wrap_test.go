package pubsubwrap

import (
	"context"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/require"
)

func TestWrapper(t *testing.T) {
	pubMock := &pubMock{}
	pw := NewPublisherWrapper(pubMock)

	mi := middleware{}
	pw.Middleware(mi)

	expectAll := 7
	err := pw.Publish("", genMessages(expectAll)...)
	require.NoError(t, err, "should no error on publish messages")

	require.Equal(t, expectAll, mi.counter(), "expect all messages are counted in middleware")
	require.Equal(t, int(expectAll/2), pubMock.published, "expect half of messages are passed wrapper")
}

type middleware map[string]int

func (mi middleware) counter() int {
	c, _ := mi["counter"]
	return c
}

func (mi middleware) Filter(m *message.Message) (*message.Message, error) {
	counter, _ := mi["counter"]
	mi["counter"] = counter + 1
	if counter%2 == 0 {
		return nil, nil
	}
	return m, nil
}

type pubMock struct {
	published int
}

func (pm *pubMock) Publish(topic string, messages ...*message.Message) error {
	pm.published += len(messages)
	return nil
}

func (pm *pubMock) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	return nil, nil
}

func (pm *pubMock) Close() error {
	return nil
}

func genMessages(count int) []*message.Message {
	out := make([]*message.Message, count)
	for i := range out {
		out[i] = &message.Message{}
	}
	return out
}
