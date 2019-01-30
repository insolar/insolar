/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package artifactmanager

import (
	"context"
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
)

// ledgerArtifactSenders is a some kind of a middleware layer
// it contains cache meta-data for calls
type ledgerArtifactSenders struct {
	codeCacheLock sync.Mutex
	codeCache     map[core.RecordRef]*cacheEntry
}

type cacheEntry struct {
	sync.Mutex
	reply core.Reply
}

func newLedgerArtifactSenders() *ledgerArtifactSenders {
	return &ledgerArtifactSenders{
		codeCache: map[core.RecordRef]*cacheEntry{},
	}
}

// cachedSender is using for caching replies
func (m *ledgerArtifactSenders) cachedSender() PreSender {
	return func(sender Sender) Sender {
		return func(ctx context.Context, msg core.Message, options *core.MessageSendOptions) (core.Reply, error) {
			codeMsg := msg.(*message.GetCode)

			m.codeCacheLock.Lock()
			entry, ok := m.codeCache[codeMsg.Code]
			if !ok {
				entry = &cacheEntry{}
				m.codeCache[codeMsg.Code] = entry
			}
			m.codeCacheLock.Unlock()

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

// followRedirectSender is using for redirecting responses with delegation token
func followRedirectSender(bus core.MessageBus) PreSender {
	return func(sender Sender) Sender {
		return func(ctx context.Context, msg core.Message, options *core.MessageSendOptions) (core.Reply, error) {
			inslog := inslogger.FromContext(ctx)
			inslog.Debug("LedgerArtifactManager.SendAndFollowRedirectSender starts ...")

			rep, err := sender(ctx, msg, options)
			if err != nil {
				return nil, err
			}

			if r, ok := rep.(core.RedirectReply); ok {
				stats.Record(ctx, statRedirects.M(1))

				redirected := r.Redirected(msg)
				inslog.Debugf("redirect reciever=%v", r.GetReceiver())

				rep, err = bus.Send(ctx, redirected, &core.MessageSendOptions{
					Token:    r.GetToken(),
					Receiver: r.GetReceiver(),
				})
				if err != nil {
					return nil, err
				}
				if _, ok := rep.(core.RedirectReply); ok {
					return nil, errors.New("double redirects are forbidden")
				}
				return rep, nil
			}

			return rep, err
		}
	}
}

// retryJetSender is using for refreshing jet-tree, if destination has no idea about a jet from message
func retryJetSender(pulseNumber core.PulseNumber, jetStorage storage.JetStorage) PreSender {
	return func(sender Sender) Sender {
		return func(ctx context.Context, msg core.Message, options *core.MessageSendOptions) (core.Reply, error) {
			inslog := inslogger.FromContext(ctx)
			inslog.Debug("LedgerArtifactManager.RetryJetSender starts ...")

			retries := jetMissRetryCount

			for retries > 0 {
				rep, err := sender(ctx, msg, options)
				if err != nil {
					return nil, err
				}

				if r, ok := rep.(*reply.JetMiss); ok {
					err := jetStorage.UpdateJetTree(ctx, pulseNumber, true, r.JetID)
					if err != nil {
						return nil, err
					}
				} else {
					return rep, err
				}

				retries--
			}

			return nil, errors.New("failed to find jet (retry limit exceeded on client)")
		}
	}
}
