// Copyright 2020 Insolar Network Ltd.
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

// +build introspection

package internal

import (
	"github.com/ThreeDotsLabs/watermill/message"
	component "github.com/insolar/component-manager"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/instrumentation/introspector"
	"github.com/insolar/insolar/instrumentation/introspector/pubsubwrap"
	"golang.org/x/net/context"
)

// PublisherWrapper setups and returns introspection wrapper for message.Publisher.
func PublisherWrapper(
	ctx context.Context,
	cm *component.Manager,
	cfg configuration.Introspection,
	pb message.Publisher,
) message.Publisher {
	pw := pubsubwrap.NewPublisherWrapper(pb)

	// init pubsub middlewares and add them to wrapper
	mStat := pubsubwrap.NewMessageStatByType()
	mLocker := pubsubwrap.NewMessageLockerByType(ctx)
	pw.Middleware(mStat)
	pw.Middleware(mLocker)

	// create introspection server with service which implements introproto.PublisherServer
	service := pubsubwrap.NewPublisherService(mLocker, mStat)
	iSrv := introspector.NewServer(cfg.Addr, service)

	// use component manager for lifecycle (component.Manager calls Start/Stop on server instance)
	cm.Register(iSrv)

	return pw
}
