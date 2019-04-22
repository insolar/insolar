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

package messagebus

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

const jetMissRetryCount = 10
const incorrectPulseRetryCount = 3

// PreSender is an alias for a function
// which is working like a `middleware` for messagebus.Send
type PreSender func(Sender) Sender

// Sender is an alias for signature of messagebus.Send
type Sender func(context.Context, insolar.Message, *insolar.MessageSendOptions) (insolar.Reply, error)

// BuildSender allows us to build a chain of PreSender before calling Sender
// The main idea of it is ability to make a different things before sending message
// For example we can cache some replies. Another example is the sendAndFollow redirect method
func BuildSender(sender Sender, preSenders ...PreSender) Sender {
	result := sender

	for i := range preSenders {
		result = preSenders[len(preSenders)-1-i](result)
	}

	return result
}

// Senders is a some kind of a middleware layer
// it contains cache meta-data for calls
type Senders struct {
	cacheLock sync.Mutex
	caches    map[string]*cacheEntry
}

type cacheEntry struct {
	sync.Mutex
	reply insolar.Reply
}

func NewSenders() *Senders {
	return &Senders{
		caches: map[string]*cacheEntry{},
	}
}

// CachedSender is using for caching replies
func (m *Senders) CachedSender(scheme insolar.PlatformCryptographyScheme) PreSender {
	return func(sender Sender) Sender {
		return func(ctx context.Context, msg insolar.Message, options *insolar.MessageSendOptions) (insolar.Reply, error) {

			msgHash := string(scheme.IntegrityHasher().Hash(message.ToBytes(msg)))

			m.cacheLock.Lock()
			entry, ok := m.caches[msgHash]
			if !ok {
				entry = &cacheEntry{}
				m.caches[msgHash] = entry
			}
			m.cacheLock.Unlock()

			entry.Lock()
			defer entry.Unlock()

			if entry.reply != nil {
				return entry.reply, nil
			}

			response, err := sender(ctx, msg, options)
			if err != nil {
				return nil, err
			}

			entry.reply = response
			return response, err
		}
	}
}

// FollowRedirectSender is using for redirecting responses with delegation token
func FollowRedirectSender(bus insolar.MessageBus) PreSender {
	return func(sender Sender) Sender {
		return func(ctx context.Context, msg insolar.Message, options *insolar.MessageSendOptions) (insolar.Reply, error) {
			rep, err := sender(ctx, msg, options)
			if err != nil {
				return nil, err
			}

			if r, ok := rep.(insolar.RedirectReply); ok {
				redirected := r.Redirected(msg)
				inslogger.FromContext(ctx).Debugf("redirect reciever=%v", r.GetReceiver())

				rep, err = bus.Send(ctx, redirected, &insolar.MessageSendOptions{
					Token:    r.GetToken(),
					Receiver: r.GetReceiver(),
				})
				if err != nil {
					return nil, err
				}
				if _, ok := rep.(insolar.RedirectReply); ok {
					return nil, errors.New("double redirects are forbidden")
				}
				return rep, nil
			}

			return rep, err
		}
	}
}

// RetryJetSender is using for refreshing jet-tree, if destination has no idea about a jet from message
func RetryJetSender(jetModifier jet.Modifier) PreSender {
	return func(sender Sender) Sender {
		return func(ctx context.Context, msg insolar.Message, options *insolar.MessageSendOptions) (insolar.Reply, error) {
			retries := jetMissRetryCount
			for retries > 0 {
				rep, err := sender(ctx, msg, options)
				if err != nil {
					return nil, err
				}

				if r, ok := rep.(*reply.JetMiss); ok {
					jetModifier.Update(ctx, r.Pulse, true, insolar.JetID(r.JetID))
				} else {
					return rep, err
				}

				retries--
			}

			return nil, errors.New("failed to find jet (retry limit exceeded on client)")
		}
	}
}

// RetryIncorrectPulse retries messages after small delay when pulses on source and destination are out of sync.
// NOTE: This is not completely correct way to behave: 1) we should wait until pulse switches, not some hardcoded time,
// 2) it should be handled by recipient and get it right with Flow "handles"
func RetryIncorrectPulse() PreSender {
	return func(sender Sender) Sender {
		return func(
			ctx context.Context, msg insolar.Message, options *insolar.MessageSendOptions,
		) (insolar.Reply, error) {
			retries := incorrectPulseRetryCount
			for {
				rep, err := sender(ctx, msg, options)
				if err == nil || !strings.Contains(err.Error(), "Incorrect message pulse") {
					return rep, err
				}

				if retries <= 0 {
					inslogger.FromContext(ctx).Warn("got incorrect message pulse too many times")
					return rep, err

				}
				retries--

				time.Sleep(100*time.Millisecond)
			}
		}
	}
}
