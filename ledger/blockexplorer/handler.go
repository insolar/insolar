package blockexplorer

import (
	"bytes"
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/pkg/errors"
	"strconv"
)

type internalHandler func(ctx context.Context, pulseNumber core.PulseNumber, parcel core.Parcel) (core.Reply, error)

// MessageHandler processes messages for local storage interaction.
type MessageHandler struct {
	db                         *storage.DB
	Bus                        core.MessageBus                 `inject:""`
	PlatformCryptographyScheme core.PlatformCryptographyScheme `inject:""`
}

// NewMessageHandler creates new handler.
func NewMessageHandler(db *storage.DB) *MessageHandler {
	return &MessageHandler{
		db: db,
	}
}

// Init initializes handlers.
func (h *MessageHandler) Init(ctx context.Context) error {

	h.Bus.MustRegister(core.TypeGetHistory, h.messagePersistingWrapper(h.handleGetHistory))
	return nil
}

func (h *MessageHandler) messagePersistingWrapper(handler internalHandler) core.MessageHandler {
	return func(ctx context.Context, genericMsg core.Parcel) (core.Reply, error) {
		err := persistMessageToDb(ctx, h.db, genericMsg.Message())
		if err != nil {
			return nil, err
		}

		lastPulseNumber, err := h.db.GetLatestPulseNumber(ctx)
		if err != nil {
			return nil, err
		}

		return handler(ctx, lastPulseNumber, genericMsg)
	}
}

func persistMessageToDb(ctx context.Context, db *storage.DB, genericMsg core.Message) error {
	lastPulse, err := db.GetLatestPulseNumber(ctx)
	if err != nil {
		return err
	}
	err = db.SetMessage(ctx, lastPulse, genericMsg)
	if err != nil {
		return err
	}

	return nil
}

func (h *MessageHandler) handleGetHistory(ctx context.Context, pulseNumber core.PulseNumber, inmsg core.Parcel) (core.Reply, error) {
	msg := inmsg.Message().(*message.GetHistory)
	idx, _, _, err := getObject(ctx, h.db, msg.Object.Record(), nil, false)
	if err != nil {
		return nil, err
	}
	var history []reply.ExplorerObject
	var current *core.RecordID

	if msg.From != nil {
		current = msg.From
	} else {
		current = idx.LatestState
	}

	counter := 0
	for current != nil {
		// We have enough results.
		if counter >= msg.Amount {
			return &reply.ExplorerList{States: history, NextState: current}, nil
		}
		counter++

		rec, err := h.db.GetRecord(ctx, current)
		if err != nil {
			return nil, errors.New("failed to retrieve object state")
		}

		switch rec.(type) {
		case record.ObjectState:
			{
				currentState, ok := rec.(record.ObjectState)
				if !ok {
					return nil, errors.New("Cannot cast to object state: " + strconv.FormatUint(uint64(rec.Type()), 10))
				}
				current = currentState.PrevStateID()

				var memory []byte
				if currentState.GetMemory() != nil {
					memory, err = h.db.GetBlob(ctx, currentState.GetMemory())
					if err != nil {
						return nil, err
					}
				}

				parcel, err := h.getParcel(ctx, currentState.GetRequest())
				if err != nil && err != errors.New("storage object not found") {
					return nil, err
				}
				history = append(history, reply.ExplorerObject{
					Parcel:    parcel,
					Memory:    memory,
					NextState: currentState.PrevStateID(),
				})
			}
		case record.Request:
			{
				currentState, ok := rec.(record.Request)
				if !ok {
					return nil, errors.New("Cannot cast to object state: " + strconv.FormatUint(uint64(rec.Type()), 10))
				}
				parcel, err := extractParcelFromRecord(currentState)
				if err != nil {
					return nil, err
				}

				history = append(history, reply.ExplorerObject{
					Memory:    nil,
					Parcel:    parcel,
					NextState: nil,
				})
				current = nil
			}
		}
	}
	return &reply.ExplorerList{States: history, NextState: nil}, nil
}

func (h *MessageHandler) getParcel(ctx context.Context, request *core.RecordID) (core.Parcel, error) {
	if request == nil {
		return nil, errors.New("Сan not get the history of the incoming request")
	}
	req, err := h.db.GetRecord(ctx, request)
	if err != nil {
		return nil, err
	}
	return extractParcelFromRecord(req)
}

func extractParcelFromRecord(rec record.Record) (core.Parcel, error) {
	parcel, err := message.Deserialize(bytes.NewBuffer(rec.(record.Request).GetPayload()))
	if err != nil {
		return nil, err
	}
	return parcel, nil
}

func getObject(
	ctx context.Context,
	s storage.Store,
	head *core.RecordID,
	state *core.RecordID,
	approved bool,
) (*index.ObjectLifeline, *core.RecordID, record.ObjectState, error) {
	idx, err := s.GetObjectIndex(ctx, head, false)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to fetch object index")
	}

	var stateID *core.RecordID
	if state != nil {
		stateID = state
	} else {
		if approved {
			stateID = idx.LatestStateApproved
		} else {
			stateID = idx.LatestState
		}
	}

	if stateID == nil {
		return nil, nil, nil, ErrStateNotAvailable
	}

	rec, err := s.GetRecord(ctx, stateID)
	if err != nil {
		return nil, nil, nil, err
	}
	stateRec, ok := rec.(record.ObjectState)
	if !ok {
		return nil, nil, nil, errors.New("invalid object record")
	}
	if stateRec.State() == record.StateDeactivation {
		return nil, nil, nil, ErrObjectDeactivated
	}

	return idx, stateID, stateRec, nil
}
