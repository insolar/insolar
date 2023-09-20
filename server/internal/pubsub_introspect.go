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
