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

package insolar

import (
	"context"
)

//go:generate minimock -i github.com/insolar/insolar/insolar.HeavySync -o ../testutils -s _mock.go

// HeavySync provides methods for sync on heavy node.
type HeavySync interface {
	Start(ctx context.Context, jet ID, pn PulseNumber) error
	StoreIndexes(ctx context.Context, jet ID, pn PulseNumber, rawIndexes map[ID][]byte) error
	StoreDrop(ctx context.Context, jetID JetID, rawDrop []byte) error
	StoreBlobs(ctx context.Context, pn PulseNumber, rawBlobs [][]byte) error
	StoreRecords(ctx context.Context, jet ID, pn PulseNumber, rawRecords [][]byte)
	Stop(ctx context.Context, jet ID, pn PulseNumber) error
	Reset(ctx context.Context, jet ID, pn PulseNumber) error
}
