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

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

const jetMissRetryCount = 10
const incorrectPulseRetryCount = 1
const flowCancelledRetryCount = 1

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

			msgHash := string(scheme.IntegrityHasher().Hash(message.MustSerialize(msg)))

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

// RetryIncorrectPulse retries messages after small delay when pulses on source and destination are out of sync.
// NOTE: This is not completely correct way to behave: 1) we should wait until pulse switches, not some hardcoded time,
// 2) it should be handled by recipient and get it right with Flow "handles"
func RetryIncorrectPulse(accessor pulse.Accessor) PreSender {
	return retryer(accessor, incorrectPulseRetryCount,
		"Incorrect message pulse",
		"[ RetryIncorrectPulse ] incorrect message pulse, retrying",
		"incorrect message pulse (retry limit exceeded on client)")
}

// RetryFlowCancelled retries message on next pulse when received flow cancelled error.
func RetryFlowCancelled(accessor pulse.Accessor) PreSender {
	return retryer(accessor, flowCancelledRetryCount,
		flow.ErrCancelled.Error(),
		"[ RetryFlowCancelled ] flow cancelled, retrying",
		"flow cancelled (retry limit exceeded on client)")
}

func retryer(accessor pulse.Accessor, retriesCount int, errSubstr string, debugStr string, err string) PreSender {
	return func(sender Sender) Sender {

		return func(ctx context.Context, msg insolar.Message, options *insolar.MessageSendOptions) (insolar.Reply, error) {
			retries := retriesCount
			var lastPulse insolar.PulseNumber
			for retries >= 0 {

				currentPulse, err := accessor.Latest(ctx)
				if err != nil {
					return nil, errors.Wrap(err, "[ retryer ] Can't get latest pulse")
				}

				if currentPulse.PulseNumber == lastPulse {
					inslogger.FromContext(ctx).Debugf("[ retryer ]  wait for pulse change. Current: %d", currentPulse.PulseNumber)
					time.Sleep(100 * time.Millisecond)
					continue
				}
				lastPulse = currentPulse.PulseNumber

				rep, err := sender(ctx, msg, options)
				if err == nil || !strings.Contains(err.Error(), errSubstr) {
					return rep, err
				}

				inslogger.FromContext(ctx).Debug(debugStr)
				retries--
			}
			return nil, errors.New(err)
		}
	}
}
