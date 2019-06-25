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

package replica

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/secrets"
	"github.com/insolar/insolar/ledger/heavy/sequence"
)

func TestReplica_Subscribe(t *testing.T) {
	var (
		pulse    = insolar.PulseNumber(10)
		at       = Position{Index: 10, Pulse: pulse}
		expected []sequence.Item
	)
	f := fuzz.New().Funcs(func(item *sequence.Item, c fuzz.Continue) {
		c.Fuzz(&item.Key)
		c.Fuzz(&item.Value)
	})
	f.NumElements(3, 10).Fuzz(&expected)
	cursor := NewJetKeeperMock(t)
	cursor.TopSyncPulseMock.Return(pulse)
	records := sequence.NewSequencerMock(t)
	records.UpdateMock.Return()
	records.FromMock.Return(expected, nil)
	keys, err := secrets.GenerateKeyPair()
	require.NoError(t, err)
	cs := cryptography.NewKeyBoundCryptographyService(keys.Private)
	parentPubKey, err := cs.GetPublicKey()
	require.NoError(t, err)
	parent := NewReplica(cursor, records, nil, NewIntegrity(cs, nil))
	child := NewReplica(cursor, records, parent, NewIntegrity(cs, parentPubKey))

	parent.Subscribe(at)
}

func TestReplica_Pull(t *testing.T) {
	var (
		pulse    = insolar.PulseNumber(10)
		from     = Position{Index: 10, Pulse: pulse}
		limit    = uint32(10)
		expected []sequence.Item
	)
	f := fuzz.New().Funcs(func(item *sequence.Item, c fuzz.Continue) {
		c.Fuzz(&item.Key)
		c.Fuzz(&item.Value)
	})
	f.NumElements(3, 10).Fuzz(&expected)
	cursor := NewJetKeeperMock(t)
	cursor.TopSyncPulseMock.Return(pulse)
	records := sequence.NewSequencerMock(t)
	records.UpdateMock.Return()
	records.FromMock.Return(expected, nil)
	keys, err := secrets.GenerateKeyPair()
	require.NoError(t, err)
	cs := cryptography.NewKeyBoundCryptographyService(keys.Private)
	parentPubKey, err := cs.GetPublicKey()
	require.NoError(t, err)
	parent := NewReplica(cursor, records, nil, NewIntegrity(cs, nil))

	packet, err := parent.Pull(0, from, limit)

	require.NoError(t, err)
	integrity := NewIntegrity(cs, parentPubKey)
	actual, err := integrity.UnwrapAndValidate(packet)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestReplica_Notify(t *testing.T) {
	var (
		expected []sequence.Item
	)
	f := fuzz.New().Funcs(func(item *sequence.Item, c fuzz.Continue) {
		c.Fuzz(&item.Key)
		c.Fuzz(&item.Value)
	})
	f.NumElements(3, 10).Fuzz(&expected)
	records := sequence.NewSequencerMock(t)
	records.UpdateMock.Return()
	records.FromMock.Return(expected, nil)
	keys, err := secrets.GenerateKeyPair()
	require.NoError(t, err)
	cs := cryptography.NewKeyBoundCryptographyService(keys.Private)
	parentPubKey, err := cs.GetPublicKey()
	require.NoError(t, err)
	parent := NewReplica(nil, records, nil, NewIntegrity(cs, nil))
	child := NewReplica(nil, records, parent, NewIntegrity(cs, parentPubKey))

	child.Notify()
}
