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
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"sync"
	"time"

	"go.opencensus.io/trace"

	"github.com/insolar/insolar/ledger/storage/pulse"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/hack"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/instrumentation/instracer"
)

const deliverRPCMethodName = "MessageBus.Deliver"

// MessageBus is component that routes application logic requests,
// e.g. glue between network and logic runner
type MessageBus struct {
	Network                    insolar.Network                    `inject:""`
	JetCoordinator             insolar.JetCoordinator             `inject:""`
	NodeNetwork                insolar.NodeNetwork                `inject:""`
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme `inject:""`
	CryptographyService        insolar.CryptographyService        `inject:""`
	DelegationTokenFactory     insolar.DelegationTokenFactory     `inject:""`
	ParcelFactory              message.ParcelFactory              `inject:""`
	PulseAccessor              pulse.Accessor                     `inject:""`

	handlers     map[insolar.MessageType]insolar.MessageHandler
	signmessages bool

	counter uint64
	span    *trace.Span

	globalLock                  sync.RWMutex
	NextPulseMessagePoolChan    chan interface{}
	NextPulseMessagePoolCounter uint32
	NextPulseMessagePoolLock    sync.RWMutex
}

func (mb *MessageBus) Acquire(ctx context.Context) {
	ctx, span := instracer.StartSpan(ctx, "NetworkSwitcher.Acquire")
	defer span.End()
	inslogger.FromContext(ctx).Info("Call Acquire in NetworkSwitcher: ", mb.counter)
	mb.counter = mb.counter + 1
	if mb.counter-1 == 0 {
		inslogger.FromContext(ctx).Info("Lock MB")
		ctx, mb.span = instracer.StartSpan(context.Background(), "GIL Lock (Lock MB)")
		mb.Lock(ctx)
	}
}

func (mb *MessageBus) Release(ctx context.Context) {
	ctx, span := instracer.StartSpan(ctx, "NetworkSwitcher.Release")
	defer span.End()
	inslogger.FromContext(ctx).Info("Call Release in NetworkSwitcher: ", mb.counter)
	if mb.counter == 0 {
		panic("Trying to unlock without locking")
	}
	mb.counter = mb.counter - 1
	if mb.counter == 0 {
		inslogger.FromContext(ctx).Info("Unlock MB")
		mb.Unlock(ctx)
		mb.span.End()
	}
}

// NewMessageBus creates plain MessageBus instance. It can be used to create Player and Recorder instances that
// wrap it, providing additional functionality.
func NewMessageBus(config configuration.Configuration) (*MessageBus, error) {
	mb := &MessageBus{
		handlers:                 map[insolar.MessageType]insolar.MessageHandler{},
		signmessages:             config.Host.SignMessages,
		NextPulseMessagePoolChan: make(chan interface{}),
	}
	mb.Acquire(context.Background())
	return mb, nil
}

// Start initializes message bus.
func (mb *MessageBus) Start(ctx context.Context) error {
	mb.Network.RemoteProcedureRegister(deliverRPCMethodName, mb.deliver)

	return nil
}

// Stop releases resources and stops the bus
func (mb *MessageBus) Stop(ctx context.Context) error { return nil }

func (mb *MessageBus) Lock(ctx context.Context) {
	inslogger.FromContext(ctx).Info("Acquire GIL")
	mb.globalLock.Lock()
}

func (mb *MessageBus) Unlock(ctx context.Context) {
	inslogger.FromContext(ctx).Info("Release GIL")
	mb.globalLock.Unlock()
}

// Register sets a function as a handler for particular message type,
// only one handler per type is allowed
func (mb *MessageBus) Register(p insolar.MessageType, handler insolar.MessageHandler) error {
	_, ok := mb.handlers[p]
	if ok {
		return errors.New("handler for this type already exists")
	}

	mb.handlers[p] = handler
	return nil
}

// MustRegister is a Register wrapper that panics if an error was returned.
func (mb *MessageBus) MustRegister(p insolar.MessageType, handler insolar.MessageHandler) {
	err := mb.Register(p, handler)
	if err != nil {
		panic(err)
	}
}

// Send an `Message` and get a `Value` or error from remote host.
func (mb *MessageBus) Send(ctx context.Context, msg insolar.Message, ops *insolar.MessageSendOptions) (insolar.Reply, error) {
	ctx, span := instracer.StartSpan(ctx, "MessageBus.Send "+msg.Type().String())
	defer span.End()

	currentPulse, err := mb.PulseAccessor.Latest(ctx)
	if err != nil {
		return nil, err
	}

	parcel, err := mb.CreateParcel(ctx, msg, ops.Safe().Token, currentPulse)
	if err != nil {
		return nil, err
	}

	rep, err := mb.SendParcel(ctx, parcel, currentPulse, ops)
	return rep, err
}

// CreateParcel creates signed message from provided message.
func (mb *MessageBus) CreateParcel(ctx context.Context, msg insolar.Message, token insolar.DelegationToken, currentPulse insolar.Pulse) (insolar.Parcel, error) {
	return mb.ParcelFactory.Create(ctx, msg, mb.NodeNetwork.GetOrigin().ID(), token, currentPulse)
}

// SendParcel sends provided message via network.
func (mb *MessageBus) SendParcel(
	ctx context.Context,
	parcel insolar.Parcel,
	currentPulse insolar.Pulse,
	options *insolar.MessageSendOptions,
) (insolar.Reply, error) {
	parcelType := parcel.Type().String()
	ctx, span := instracer.StartSpan(ctx, "MessageBus.SendParcel "+parcelType)
	ctx = insmetrics.InsertTag(ctx, tagMessageType, parcelType)
	defer span.End()

	readBarrier(ctx, &mb.globalLock)

	var (
		nodes []insolar.Reference
		err   error
	)
	if options != nil && options.Receiver != nil {
		nodes = []insolar.Reference{*options.Receiver}
	} else {
		// TODO: send to all actors of the role if nil Target
		target := parcel.DefaultTarget()
		// FIXME: @andreyromancev. 21.12.18. Temp hack. All messages should have a default target.
		if target == nil {
			target = &insolar.Reference{}
		}
		nodes, err = mb.JetCoordinator.QueryRole(ctx, parcel.DefaultRole(), *target.Record(), currentPulse.PulseNumber)
		if err != nil {
			return nil, err
		}
	}

	start := time.Now()
	defer func() {
		stats.Record(ctx, statParcelsTime.M(float64(time.Since(start).Nanoseconds())/1e6))
	}()

	stats.Record(ctx, statParcelsSentTotal.M(1))

	if len(nodes) > 1 {
		cascade := insolar.Cascade{
			NodeIds:           nodes,
			Entropy:           currentPulse.Entropy,
			ReplicationFactor: 2,
		}
		err := mb.Network.SendCascadeMessage(cascade, deliverRPCMethodName, parcel)
		return nil, err
	}

	// Short path when sending to self node. Skip serialization
	origin := mb.NodeNetwork.GetOrigin()
	if nodes[0].Equal(origin.ID()) {
		stats.Record(ctx, statLocallyDeliveredParcelsTotal.M(1))
		return mb.doDeliver(parcel.Context(context.Background()), parcel)
	}

	res, err := mb.Network.SendMessage(nodes[0], deliverRPCMethodName, parcel)
	if err != nil {
		return nil, err
	}

	return reply.Deserialize(bytes.NewBuffer(res))
}

type serializableError struct {
	S string
}

func (e *serializableError) Error() string {
	return e.S
}

func (mb *MessageBus) OnPulse(context.Context, insolar.Pulse) error {
	close(mb.NextPulseMessagePoolChan)

	mb.NextPulseMessagePoolLock.Lock()
	defer mb.NextPulseMessagePoolLock.Unlock()

	mb.NextPulseMessagePoolChan = make(chan interface{})

	return nil
}

func (mb *MessageBus) doDeliver(ctx context.Context, msg insolar.Parcel) (insolar.Reply, error) {

	var err error
	ctx, span := instracer.StartSpan(ctx, "MessageBus.doDeliver")
	defer span.End()
	if err = mb.checkPulse(ctx, msg, false); err != nil {
		return nil, errors.Wrap(err, "[ doDeliver ] error in checkPulse")
	}

	// We must check barrier just before exiting function
	// to deliver reply right after pulse switches if it is switching right now.
	defer readBarrier(ctx, &mb.globalLock)
	ctx, _ = inslogger.WithField(ctx, "msg_type", msg.Type().String())
	inslogger.FromContext(ctx).Debug("MessageBus.doDeliver starts ...")
	handler, ok := mb.handlers[msg.Type()]
	if !ok {
		txt := "no handler for received message type"
		inslogger.FromContext(ctx).Error(txt)
		return nil, errors.New(txt)
	}

	origin := mb.NodeNetwork.GetOrigin()
	if msg.GetSender().Equal(origin.ID()) {
		ctx = hack.SetSkipValidation(ctx, true)
	}
	// TODO: sergey.morozov 2018-12-21 there is potential race condition because of readBarrier. We must implement correct locking.

	resp, err := handler(ctx, msg)
	if err != nil {
		return nil, &serializableError{
			S: err.Error(),
		}
	}

	return resp, nil
}

func (mb *MessageBus) checkPulse(ctx context.Context, parcel insolar.Parcel, locked bool) error {
	ctx, span := instracer.StartSpan(ctx, "MessageBus.checkPulse")
	defer span.End()

	pulse, err := mb.PulseAccessor.Latest(ctx)
	if err != nil {
		return errors.Wrap(err, "[ checkPulse ] Couldn't get current pulse number")
	}

	ppn := parcel.Pulse()
	if ppn > pulse.PulseNumber {
		return mb.handleParcelFromTheFuture(ctx, parcel, locked)
	} else if ppn < pulse.PulseNumber {
		if ppn < pulse.PrevPulseNumber {
			inslogger.FromContext(ctx).Errorf(
				"[ checkPulse ] Pulse is TOO OLD: (parcel: %d, current: %d) Parcel is: %#v",
				ppn, pulse.PulseNumber, parcel.Message(),
			)
		}

		// Parcel is from past. Return error for some messages, allow for others.
		switch parcel.Message().(type) {
		case
			*message.GetObject,
			*message.GetDelegate,
			*message.GetChildren,
			*message.SetRecord,
			*message.UpdateObject,
			*message.RegisterChild,
			*message.SetBlob,
			*message.GetObjectIndex,
			*message.GetPendingRequests,
			*message.ValidateRecord,
			*message.CallConstructor,
			*message.HotData,
			*message.CallMethod:
			inslogger.FromContext(ctx).Errorf("[ checkPulse ] Incorrect message pulse (parcel: %d, current: %d)", ppn, pulse.PulseNumber)
			return fmt.Errorf("[ checkPulse ] Incorrect message pulse (parcel: %d, current: %d)", ppn, pulse.PulseNumber)
		}
	}

	return nil
}

func (mb *MessageBus) handleParcelFromTheFuture(ctx context.Context, parcel insolar.Parcel, locked bool) error {
	ctx, span := instracer.StartSpan(ctx, "MessageBus.handleParcelFromTheFuture")
	defer span.End()

	ppn := parcel.Pulse()
	inslogger.FromContext(ctx).Debug(
		"message from the future, msg pulse: ", ppn,
	)
	if locked {
		mb.globalLock.RUnlock()
		defer mb.globalLock.RLock()
	}

	for {
		mb.NextPulseMessagePoolLock.RLock()

		pulse, err := mb.PulseAccessor.Latest(ctx)
		if err != nil {
			mb.NextPulseMessagePoolLock.RUnlock()
			return errors.Wrap(err, "couldn't get current pulse number")
		}
		if ppn > pulse.PulseNumber {
			inslogger.FromContext(ctx).Debug("still in future")

			_, span := instracer.StartSpan(
				ctx, fmt.Sprintf("waiting pulse switch from %d to %d", pulse.PulseNumber, ppn),
			)
			<-mb.NextPulseMessagePoolChan
			span.End()

			pulse, err = mb.PulseAccessor.Latest(ctx)
			if err != nil {
				mb.NextPulseMessagePoolLock.RUnlock()
				return errors.Wrap(err, "couldn't get current pulse number")
			}
		}

		mb.NextPulseMessagePoolLock.RUnlock()

		if ppn <= pulse.PulseNumber {
			inslogger.FromContext(ctx).Debug("releasing message after waiting for pulse")
			return nil
		}
	}
}

// Deliver method calls LogicRunner.Execute on local host
// this method is registered as RPC stub
func (mb *MessageBus) deliver(ctx context.Context, args [][]byte) (result []byte, err error) {
	inslogger.FromContext(ctx).Debug("MessageBus.deliver starts ...")
	if len(args) < 1 {
		return nil, errors.New("need exactly one argument when mb.deliver()")
	}
	parcel, err := message.DeserializeParcel(bytes.NewBuffer(args[0]))
	if err != nil {
		return nil, err
	}

	parcelCtx := parcel.Context(context.Background()) // use ctx when network provide context
	inslogger.FromContext(ctx).Debugf("MessageBus.deliver after deserialize msg. Msg Type: %s", parcel.Type())

	mb.globalLock.RLock()

	if err = mb.checkPulse(parcelCtx, parcel, true); err != nil {
		mb.globalLock.RUnlock()
		return nil, err
	}

	if err = mb.checkParcel(parcelCtx, parcel); err != nil {
		mb.globalLock.RUnlock()
		return nil, err
	}
	mb.globalLock.RUnlock()

	resp, err := mb.doDeliver(parcelCtx, parcel)
	if err != nil {
		return nil, err
	}

	rd, err := reply.Serialize(resp)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(rd)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (mb *MessageBus) checkParcel(ctx context.Context, parcel insolar.Parcel) error {
	sender := parcel.GetSender()

	if mb.signmessages {
		senderKey := mb.NodeNetwork.GetWorkingNode(sender).PublicKey()
		if err := mb.ParcelFactory.Validate(senderKey, parcel); err != nil {
			return errors.Wrap(err, "failed to check a message sign")
		}
	}

	// FIXME: @andreyromancev. 09.01.2019. Implement verify method.
	// if parcel.DelegationToken() != nil {
	// 	valid, err := mb.DelegationTokenFactory.Verify(parcel)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if !valid {
	// 		return errors.New("delegation token is not valid")
	// 	}
	// 	return nil
	// }

	// sendingObject, allowedSenderRole := parcel.AllowedSenderObjectAndRole()
	// if sendingObject == nil {
	// 	return nil
	// }
	//
	// // TODO: temporary solution, this check should be removed after reimplementing token processing on VM side - @nordicdyno 19.Dec.2018
	// if allowedSenderRole.IsVirtualRole() {
	// 	return nil
	// }
	//
	// validSender, err := mb.JetCoordinator.IsAuthorized(
	// 	ctx, allowedSenderRole, *sendingObject.Record(), parcel.Pulse(), sender,
	// )
	// if err != nil {
	// 	return err
	// }
	// if !validSender {
	// 	return errors.New("sender is not allowed to act on behalve of that object")
	// }
	return nil
}

func readBarrier(ctx context.Context, mutex *sync.RWMutex) {
	inslogger.FromContext(ctx).Debug("Locking readBarrier")
	mutex.RLock()
	inslogger.FromContext(ctx).Debug("readBarrier locked")
	mutex.RUnlock()
	inslogger.FromContext(ctx).Debug("readBarrier unlocked")
}

func init() {
	gob.Register(&serializableError{})
}
