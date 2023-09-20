// +build !introspection

package internal

import (
	"github.com/ThreeDotsLabs/watermill/message"
	component "github.com/insolar/component-manager"
	"github.com/insolar/insolar/configuration"
	"golang.org/x/net/context"
)

// PublisherWrapper stub for message.Publisher introspection wrapper for binaries without introspection API.
func PublisherWrapper(
	ctx context.Context, cm *component.Manager, cfg configuration.Introspection, pb message.Publisher,
) message.Publisher {
	return pb
}
