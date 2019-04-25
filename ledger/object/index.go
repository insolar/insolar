package object

import (
	"context"

	"github.com/insolar/insolar/insolar"
)

type IndexAccessor interface {
	ForPnAndID(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID)
}

type IndexModifier interface {
	Set(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, index Lifeline) error
	SetRequest(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, req Request)
	SetResultRecord(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, res ResultRecord)
}

type InMemoryIndex struct {
}

func (*InMemoryIndex) Set(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, index Lifeline) error {
	panic("implement me")
}

func (*InMemoryIndex) SetRequest(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, req Request) {
	panic("implement me")
}

func (*InMemoryIndex) SetResultRecord(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, res ResultRecord) {
	panic("implement me")
}
