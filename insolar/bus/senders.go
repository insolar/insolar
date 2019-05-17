package bus

import (
	"context"
	"strings"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

type Middleware func(Sender) Sender

// BuildSender allows us to build a chain of PreSender before calling Sender
// The main idea of it is ability to make a different things before sending message
// For example we can cache some replies. Another example is the sendAndFollow redirect method
func BuildSender(sender Sender, mware ...Middleware) Sender {
	result := sender

	for i := range mware {
		result = mware[len(mware)-1-i](result)
	}

	return result
}

// NewRetryJet is using for refreshing jet-tree, if destination has no idea about a jet from message.
func NewRetryJet(jetModifier jet.Modifier) Middleware {
	const jetMissRetryCount = 10

	return func(sender Sender) Sender {
		f := func(ctx context.Context, msg *message.Message) (<-chan *message.Message, func()) {
			retries := jetMissRetryCount

			var (
				reps <-chan *message.Message
				done func()
			)
			for retries > 0 {
				reps, done = sender.Send(ctx, msg)

				rep, ok := <-reps
				if !ok {
					return reps, done
				}

				if rep.Metadata.Get(MetaType) != payload.TypeJet {
					ch := make(chan *message.Message, 1)
					ch <- rep
					return ch, func() {
						done()
						close(ch)
					}
				}

				pl := payload.Jet{}
				err := pl.Unmarshal(rep.Payload)
				if err != nil {
					inslogger.FromContext(ctx).Error(err, "failed to decode reply")
					return reps, done
				}
				jetModifier.Update(ctx, pl.Pulse, true, pl.JetID)
				done()
				retries--
			}

			return nil, errors.New("failed to find jet (retry limit exceeded on client)")
		}
		return &sendFunc{send: f}
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

				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

type sendFunc struct {
	send func(ctx context.Context, msg *message.Message) (<-chan *message.Message, func())
}

func (f *sendFunc) Send(ctx context.Context, msg *message.Message) (<-chan *message.Message, func()) {
	return f.send(ctx, msg)
}
