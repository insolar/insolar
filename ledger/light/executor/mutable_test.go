// Copyright 2020 Insolar Network Ltd.
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

package executor_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/light/executor"
)

func TestOldestMutable(t *testing.T) {
	t.Run("return nothing if receive nothing", func(t *testing.T) {
		res := executor.OldestMutable(nil)
		require.Nil(t, res)
	})

	t.Run("ignore outgoings and immutable", func(t *testing.T) {
		o := record.OutgoingRequest{}
		outCom := record.CompositeFilamentRecord{
			RecordID: *insolar.NewID(insolar.PulseNumber(8), nil),
			Record: record.Material{
				Polymorph: 0,
				Virtual:   record.Wrap(&o),
			},
		}

		inImmutable := record.IncomingRequest{
			Immutable: true,
		}
		inImCom := record.CompositeFilamentRecord{
			RecordID: *insolar.NewID(insolar.PulseNumber(9), nil),
			Record: record.Material{
				Polymorph: 0,
				Virtual:   record.Wrap(&inImmutable),
			},
		}

		inMutable := record.IncomingRequest{
			Immutable: false,
		}
		expectedID := insolar.NewID(insolar.PulseNumber(10), nil)
		expected := record.CompositeFilamentRecord{
			RecordID: *expectedID,
			Record: record.Material{
				Polymorph: 0,
				Virtual:   record.Wrap(&inMutable),
			},
		}

		res := executor.OldestMutable([]record.CompositeFilamentRecord{
			outCom, inImCom, expected,
		})

		require.NotNil(t, res)
		require.Equal(t, res.RecordID, *expectedID)
	})
}
