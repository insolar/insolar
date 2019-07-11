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
	"github.com/ThreeDotsLabs/watermill/message"
)

// FilterMiddleware an interface for message filtering and modification.
type FilterMiddleware interface {
	Filter(m *message.Message) *message.Message
}

// PubSubWrapper wraps message Publisher and Subscriber.
type PubSubWrapper struct {
	message.Subscriber
	pub message.Publisher

	filters []FilterMiddleware
}

// NewPubSubWrapper creates new message.PubSub wrapper.
func NewPubSubWrapper(pubsub message.PubSub) *PubSubWrapper {
	return &PubSubWrapper{
		Subscriber: pubsub,
		pub:        pubsub,
	}
}

// Middleware adds middleware filters (order matters!).
func (p *PubSubWrapper) Middleware(fm ...FilterMiddleware) {
	p.filters = append(p.filters, fm...)
}

// Publish wraps message.Publish method, i.e. applies all middleware filters for every message.
func (p *PubSubWrapper) Publish(topic string, messages ...*message.Message) error {
	out := make([]*message.Message, 0, len(messages))
	for _, m := range messages {
		for _, f := range p.filters {
			if m = f.Filter(m); m == nil {
				break
			}
		}
		if m != nil {
			out = append(out, m)
		}
	}
	return p.pub.Publish(topic, out...)
}

// Close wraps message.Close method.
func (p *PubSubWrapper) Close() error {
	return p.pub.Close()
}
