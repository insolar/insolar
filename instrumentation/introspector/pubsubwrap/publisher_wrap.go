// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pubsubwrap

import (
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/pkg/errors"
)

// FilterMiddleware an interface for message filtering and modification.
type FilterMiddleware interface {
	Filter(m *message.Message) (*message.Message, error)
}

// PublisherWrapper wraps message Publisher.
type PublisherWrapper struct {
	pub message.Publisher

	filters []FilterMiddleware
}

// NewPublisherWrapper creates new message.Publisher wrapper.
func NewPublisherWrapper(pb message.Publisher) *PublisherWrapper {
	return &PublisherWrapper{
		pub: pb,
	}
}

// Middleware adds middleware filters (order matters!).
func (p *PublisherWrapper) Middleware(fm ...FilterMiddleware) {
	p.filters = append(p.filters, fm...)
}

// Publish wraps message.Publish method, i.e. applies all middleware filters for every message.
func (p *PublisherWrapper) Publish(topic string, messages ...*message.Message) error {
	if topic == bus.TopicOutgoing {
		return p.pub.Publish(topic, messages...)
	}

	out := make([]*message.Message, 0, len(messages))
	for _, m := range messages {
	FiltersLoop:
		for _, f := range p.filters {
			var err error
			m, err = f.Filter(m)
			if err != nil {
				switch err.(type) {
				case decodeError:
					fmt.Printf("pubsubwrap [middleware %T]: failed to decode message: %v", f, err)
					break FiltersLoop
				default:
					panic(errors.Errorf(
						"pubsubwrap [middleware %T]: unexpected filter error: %v", f, err))
				}
			}
			// message filtered, skip other filters
			if m == nil {
				break FiltersLoop
			}
		}
		if m != nil {
			out = append(out, m)
		}
	}
	return p.pub.Publish(topic, out...)
}

// Close wraps message.Close method.
func (p *PublisherWrapper) Close() error {
	return p.pub.Close()
}
