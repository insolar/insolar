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

// +build introspection

package internal

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar/bus/pubsubwrap"
	"github.com/insolar/insolar/instrumentation/introspector"
	"golang.org/x/net/context"
)

// PubSubWrapper introspection wrapper of message.PubSub.
func PubSubWrapper(ctx context.Context, cm *component.Manager, cfg configuration.Introspection, pb message.PubSub) message.PubSub {
	pw := pubsubwrap.NewPubSubWrapper(pb)

	mStat := pubsubwrap.NewMessageStatByType()
	mLocker := pubsubwrap.NewMessageLockerByType(ctx)
	pw.Middleware(mStat, mLocker)

	service := pubsubwrap.NewPublisherService(mLocker, mStat)
	iSrv := introspector.NewServer(cfg.Addr, service)

	// TODO: do not use component manager for lifecycle
	cm.Register(iSrv)

	return pw
}
