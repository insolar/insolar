package artifactmanager

import (
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
)

// ObjectDescriptor represents meta info required to fetch all object data
type ObjectDescriptor struct {
	StateRef record.Reference

	activateRecord    *record.ObjectActivateRecord
	latestAmendRecord *record.ObjectAmendRecord
	lifelineIndex     *index.ObjectLifeline
}
