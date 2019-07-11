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

// func TestReplicatorRoot_Start(t *testing.T) {
// 	var (
// 		ctx   = inslogger.TestContext(t)
// 		pulse = insolar.GenesisPulse.PulseNumber
// 	)
// 	net := network.NewHostNetworkMock(t)
// 	net.RegisterRequestHandlerMock.Return()
// 	jetKeeper := NewJetKeeperMock(t)
// 	jetKeeper.TopSyncPulseMock.Return(pulse)
// 	db := store.NewMemoryMockDB()
// 	cs := testutils.NewCryptographyServiceMock(t)
// 	config := configuration.Replica{
// 		Role:         "root",
// 		In:           "in",
// 		Out:          []string{"in"},
// 		ParentCertPath: "",
// 	}
// 	replicator := NewReplicator(config, jetKeeper, cs)
// 	replicator.DB = db
// 	replicator.ServiceNetwork = net
//
// 	err := replicator.Init(ctx)
// 	require.NoError(t, err)
// 	err = replicator.Start(ctx)
// 	require.NoError(t, err)
// }

// func TestReplicatorReplica_Start(t *testing.T) {
// 	var (
// 		ctx   = inslogger.TestContext(t)
// 		pulse = insolar.GenesisPulse.PulseNumber
// 	)
// 	net := network.NewHostNetworkMock(t)
// 	net.RegisterRequestHandlerMock.Return()
// 	net.SendRequestToHostMock.Return(makeFuture([]byte{}), nil)
// 	jetKeeper := NewJetKeeperMock(t)
// 	jetKeeper.TopSyncPulseMock.Return(pulse)
// 	db := store.NewMemoryMockDB()
// 	cs := testutils.NewCryptographyServiceMock(t)
// 	config := configuration.Replica{
// 		Role:         "replica",
// 		In:           "inside",
// 		Out:         []string{"inside"},
// 		ParentCertPath: "",
// 	}
// 	replicator := NewReplicator(config, jetKeeper, )
//
// 	err := replicator.Init(ctx)
// 	require.NoError(t, err)
// 	err = replicator.Start(ctx)
// 	require.NoError(t, err)
// }
//
// func TestReplicatorObserver_Start(t *testing.T) {
// 	var (
// 		ctx   = inslogger.TestContext(t)
// 		pulse = insolar.GenesisPulse.PulseNumber
// 	)
// 	net := network.NewHostNetworkMock(t)
// 	net.RegisterRequestHandlerMock.Return()
// 	net.SendRequestToHostMock.Return(makeFuture([]byte{}), nil)
// 	jetKeeper := NewJetKeeperMock(t)
// 	jetKeeper.TopSyncPulseMock.Return(pulse)
// 	db := store.NewMemoryMockDB()
// 	cs := testutils.NewCryptographyServiceMock(t)
// 	config := configuration.Replica{
// 		Role:         "observer",
// 		In:           "inside",
// 		Out:         []string{"inside"},
// 		ParentCertPath: "",
// 	}
// 	replicator := NewReplicator(config, jetKeeper, db, cs)
// 	replicator.Network = net
//
// 	err := replicator.Init(ctx)
// 	require.NoError(t, err)
// 	err = replicator.Start(ctx)
// 	require.NoError(t, err)
// }
