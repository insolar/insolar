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
	"time"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/heavy/sequence"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/testutils"
)

func TestReplicatorRoot_InitStart(t *testing.T) {
	var (
		ctx = inslogger.TestContext(t)
		pn  = insolar.GenesisPulse.PulseNumber
	)
	JetKeeper := NewJetKeeperMock(t)
	JetKeeper.TopSyncPulseMock.Return(pn)
	db := store.NewMemoryMockDB()
	sequencer := sequence.NewSequencer(db)
	cs := testutils.NewCryptographyServiceMock(t)
	config := configuration.Configuration{
		Ledger: configuration.Ledger{
			Replica: configuration.Replica{
				Role:              "root",
				ParentAddress:     "127.0.0.1:13831",
				ParentPubKey:      "",
				ScopesToReplicate: []byte{2},
				Attempts:          60,
				DelayForAttempt:   1 * time.Second,
				DefaultBatchSize:  uint32(1000),
			},
		},
	}
	serviceNetwork, _ := servicenetwork.NewServiceNetwork(config, nil, false)
	transport := NewTransport(serviceNetwork)
	replicator := NewReplicator(config, JetKeeper)
	replicator.Sequencer = sequencer
	replicator.CryptoService = cs
	replicator.Transport = transport

	err := replicator.Init(ctx)
	require.NoError(t, err)
	err = replicator.Start(ctx)
	require.NoError(t, err)
}

func TestReplicatorReplica_InitStart(t *testing.T) {
	var (
		ctx     = inslogger.TestContext(t)
		pn      = insolar.GenesisPulse.PulseNumber
		address = "127.0.0.1:13831"
	)
	JetKeeper := NewJetKeeperMock(t)
	JetKeeper.TopSyncPulseMock.Return(pn)
	db := store.NewMemoryMockDB()
	sequencer := sequence.NewSequencer(db)
	cs := testutils.NewCryptographyServiceMock(t)
	config := configuration.Configuration{
		Ledger: configuration.Ledger{
			Replica: configuration.Replica{
				Role:              "replica",
				ParentAddress:     address,
				ParentPubKey:      "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE3IbYswQbBlxMg0dFNzxN+3hY3VWP\nyCL8T/XEx3trWuFn74M5vGvbabvJSkF6U8Qlq6mwQ8zx7teIFML7IlQwwg==\n-----END PUBLIC KEY-----\n",
				ScopesToReplicate: []byte{2},
				Attempts:          60,
				DelayForAttempt:   1 * time.Second,
				DefaultBatchSize:  uint32(1000),
			},
		},
	}
	transport := NewTransportMock(t)
	transport.RegisterMock.Return()
	transport.MeMock.Return(address)
	reply, _ := insolar.Serialize(GenericReply{Data: []byte{}, Error: nil})
	transport.SendMock.Return(reply, nil)
	replicator := NewReplicator(config, JetKeeper)
	replicator.Sequencer = sequencer
	replicator.CryptoService = cs
	replicator.Transport = transport

	err := replicator.Init(ctx)
	require.NoError(t, err)
	err = replicator.Start(ctx)
	require.NoError(t, err)
}
