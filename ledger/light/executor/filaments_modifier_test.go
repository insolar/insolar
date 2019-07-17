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

package executor

import (
	"context"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/require"
)

func TestFilamentModifierDefault_CheckOutgoingHasReasonInPendings(t *testing.T) {
	ctx := inslogger.TestContext(t)

	var pn insolar.PulseNumber = insolar.FirstPulseNumber + 100
	uniqID := func() insolar.ID {
		pn++
		return gen.IDWithPulse(pn)
	}
	var (
		requestID          = uniqID()
		lifelineObjectID   = uniqID()
		pendingObjectID    = uniqID()
		notPendingObjectID = uniqID()
	)

	calc := NewFilamentCalculatorMock(t)
	calc.PendingRequestsFunc = func(_ context.Context, _ insolar.PulseNumber, id insolar.ID) ([]insolar.ID, error) {
		return []insolar.ID{pendingObjectID}, nil
	}

	request := record.NewRequestMock(t)

	modifier := &FilamentModifierDefault{calculator: calc}

	request.ReasonRefMock.Return(*insolar.NewReference(notPendingObjectID))
	err := modifier.checkOutgoingHasReasonInPendings(ctx, requestID, lifelineObjectID, request)
	require.Error(t, err, "error if request reason not in pendings")
	require.IsType(t, &errorReasonNotInPendings{}, err, "error has errorReasonNotInPendings type")

	request.ReasonRefMock.Return(*insolar.NewReference(pendingObjectID))
	err = modifier.checkOutgoingHasReasonInPendings(ctx, requestID, lifelineObjectID, request)
	require.NoError(t, err, "no error if request reason in pendings")
}
