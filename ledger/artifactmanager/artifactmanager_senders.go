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

type ledgerArtifactSenders struct {
	pulseStorage core.PulseStorage
	jetStorage   storage.JetStorage
	defaultBus   core.MessageBus

	codeCacheLock sync.Mutex
	codeCache     map[core.RecordRef]*cacheEntry
}

func newLedgerArtifactSenders(
	pulseStorage core.PulseStorage,
	jetStorage storage.JetStorage,
	defaultBus core.MessageBus,
) *ledgerArtifactSenders {
	return &ledgerArtifactSenders{
		pulseStorage: pulseStorage,
		jetStorage:   jetStorage,
		defaultBus:   defaultBus,
		codeCache:    map[core.RecordRef]*cacheEntry{},
	}
}

func (m *ledgerArtifactSenders) bus(ctx context.Context) core.MessageBus {
	return core.MessageBusFromContext(ctx, m.defaultBus)
}

func (m *ledgerArtifactSenders) cachedSender(sender Sender) Sender {
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

		if entry.code != nil {
			return entry.code, nil
		}

		response, err := sender(ctx, msg, options)
		if err != nil {
			return nil, err
		}
		castedResp, ok := response.(*reply.Code)
		if !ok {
			return response, err
		}

		entry.code = castedResp
		return response, err
	}
}

func (m *ledgerArtifactSenders) followRedirectSender(sender Sender) Sender {
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

			rep, err = m.bus(ctx).Send(ctx, redirected, &core.MessageSendOptions{
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

func (m *ledgerArtifactSenders) retryJetSender(sender Sender) Sender {
	return func(ctx context.Context, msg core.Message, options *core.MessageSendOptions) (core.Reply, error) {
		inslog := inslogger.FromContext(ctx)
		inslog.Debug("LedgerArtifactManager.RetryJetSender starts ...")

		retries := jetMissRetryCount

		currentPulse, err := m.pulseStorage.Current(ctx)
		if err != nil {
			return nil, err
		}

		for retries > 0 {
			rep, err := sender(ctx, msg, options)
			if err != nil {
				return nil, err
			}

			if r, ok := rep.(*reply.JetMiss); ok {
				err := m.jetStorage.UpdateJetTree(ctx, currentPulse.PulseNumber, true, r.JetID)
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
