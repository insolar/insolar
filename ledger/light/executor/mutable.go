package executor

import (
	"github.com/insolar/insolar/insolar/record"
)

// OldestMutable searches for a oldest mutable request through a provided list of open requests
// openRequests MUST be time-ascending order
func OldestMutable(openRequests []record.CompositeFilamentRecord) *record.CompositeFilamentRecord {
	isMutableIncoming := func(rec record.CompositeFilamentRecord) bool {
		req := record.Unwrap(&rec.Record.Virtual).(record.Request)
		inReq, isIn := req.(*record.IncomingRequest)
		return isIn && !inReq.Immutable
	}

	if len(openRequests) == 0 {
		return nil
	}

	for i := 0; i < len(openRequests); i++ {
		if isMutableIncoming(openRequests[i]) {
			return &openRequests[i]
		}
	}

	return nil
}
