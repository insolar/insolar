package object

import (
	"context"

	"github.com/insolar/insolar/insolar"
)

type IndexMetaAccessor interface {
	ForPnAndID(ctx context.Context, pn insolar.PulseNumber, id insolar.ID)
}

type IndexMetaModifier interface {
	Set(ctx context.Context, pn insolar.PulseNumber, id insolar.ID, index Lifeline) error
	SetRequest(ctx context.Context, pn insolar.PulseNumber)
}
