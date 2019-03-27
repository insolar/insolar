package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
)

const jetMissRetryCount = 10

// followRedirectSender is using for redirecting responses with delegation token
func followRedirectSender(bus insolar.MessageBus) PreSender {
	return func(sender Sender) Sender {
		return func(ctx context.Context, msg insolar.Message, options *insolar.MessageSendOptions) (insolar.Reply, error) {
			rep, err := sender(ctx, msg, options)
			if err != nil {
				return nil, err
			}

			if r, ok := rep.(insolar.RedirectReply); ok {
				stats.Record(ctx, statRedirects.M(1))

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

// retryJetSender is using for refreshing jet-tree, if destination has no idea about a jet from message
func retryJetSender(pulseNumber insolar.PulseNumber, jetModifier jet.Modifier) PreSender {
	return func(sender Sender) Sender {
		return func(ctx context.Context, msg insolar.Message, options *insolar.MessageSendOptions) (insolar.Reply, error) {
			retries := jetMissRetryCount
			for retries > 0 {
				rep, err := sender(ctx, msg, options)
				if err != nil {
					return nil, err
				}

				if r, ok := rep.(*reply.JetMiss); ok {
					jetModifier.Update(ctx, pulseNumber, true, insolar.JetID(r.JetID))
				} else {
					return rep, err
				}

				retries--
			}

			return nil, errors.New("failed to find jet (retry limit exceeded on client)")
		}
	}
}
