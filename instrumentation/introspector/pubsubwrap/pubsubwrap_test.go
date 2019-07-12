//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package pubsubwrap

import (
	"context"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/require"
)

func TestWrapper(t *testing.T) {
	psMock := &pubsubMock{}
	pw := NewPubSubWrapper(psMock)

	mi := middleware{}
	pw.Middleware(mi)

	expectAll := 7
	err := pw.Publish("", genMessages(expectAll)...)
	require.NoError(t, err, "should no error on publish messages")

	require.Equal(t, expectAll, mi.counter(), "expect all messages are counted in middleware")
	require.Equal(t, int(expectAll/2), psMock.published, "expect half of messages are passed wrapper")
}

type middleware map[string]int

func (mi middleware) counter() int {
	c, _ := mi["counter"]
	return c
}

func (mi middleware) Filter(m *message.Message) *message.Message {
	counter, _ := mi["counter"]
	mi["counter"] = counter + 1
	if counter%2 == 0 {
		return nil
	}
	return m
}

type pubsubMock struct {
	published int
}

func (pm *pubsubMock) Publish(topic string, messages ...*message.Message) error {
	pm.published += len(messages)
	return nil
}

func (pm *pubsubMock) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	return nil, nil
}

func (pm *pubsubMock) Close() error {
	return nil
}

func genMessages(count int) []*message.Message {
	out := make([]*message.Message, count)
	for i := range out {
		out[i] = &message.Message{}
	}
	return out
}
