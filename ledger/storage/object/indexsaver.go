package object

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/pkg/errors"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/storage/object.IndexSaver -o ./ -s _mock.go

type IndexSaver interface {
	SaveIndexFromHeavy(
		ctx context.Context, jetID insolar.ID, obj insolar.Reference, heavy *insolar.Reference,
	) (Lifeline, error)
}

type indexSaverImpl struct {
	Bus           insolar.MessageBus `inject:""`
	IndexModifier IndexModifier      `inject:""`
}

func NewIndexSaver(bus insolar.MessageBus, indexModifier IndexModifier) IndexSaver {
	return &indexSaverImpl{
		Bus:           bus,
		IndexModifier: indexModifier,
	}
}

func (h *indexSaverImpl) SaveIndexFromHeavy(
	ctx context.Context, jetID insolar.ID, obj insolar.Reference, heavy *insolar.Reference,
) (Lifeline, error) {
	genericReply, err := h.Bus.Send(ctx, &message.GetObjectIndex{
		Object: obj,
	}, &insolar.MessageSendOptions{
		Receiver: heavy,
	})
	if err != nil {
		return Lifeline{}, errors.Wrap(err, "failed to send")
	}
	rep, ok := genericReply.(*reply.ObjectIndex)
	if !ok {
		return Lifeline{}, fmt.Errorf("failed to fetch object index: unexpected reply type %T (reply=%+v)", genericReply, genericReply)
	}
	idx, err := DecodeIndex(rep.Index)
	if err != nil {
		return Lifeline{}, errors.Wrap(err, "failed to decode")
	}

	idx.JetID = insolar.JetID(jetID)
	err = h.IndexModifier.Set(ctx, *obj.Record(), idx)
	if err != nil {
		return Lifeline{}, errors.Wrap(err, "failed to save")
	}
	return idx, nil

}
