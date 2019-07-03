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

package stubs

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/utils"
)

func TestRoutingTable_AddToKnownHosts(t *testing.T) {
	rt := RoutingTable{}

	require.NotPanics(t, func() { rt.AddToKnownHosts(nil) })
}

func TestRoutingTable_Rebalance(t *testing.T) {
	rt := RoutingTable{}

	require.Panics(t, func() { rt.Rebalance(nil) })
}

func TestRoutingTable_Resolve(t *testing.T) {
	rt := RoutingTable{}

	ref := gen.Reference()
	host, err := host.NewHostN("127.0.0.1:8080", ref)
	require.NoError(t, err)
	rt.AddToKnownHosts(host)

	host, err = rt.Resolve(ref)
	require.Error(t, err)
	require.Nil(t, host)
}

func TestRoutingTable_ResolveConsensus(t *testing.T) {
	rt := RoutingTable{}

	ref := gen.Reference()
	short := insolar.ShortNodeID(utils.GenerateUintShortID(ref))
	host, err := host.NewHostNS("127.0.0.1:8080", ref, short)
	require.NoError(t, err)
	rt.AddToKnownHosts(host)

	host, err = rt.ResolveConsensus(short)
	require.Error(t, err)
	require.Nil(t, host)
}

func TestRoutingTable_ResolveConsensusRef(t *testing.T) {
	rt := RoutingTable{}

	ref := gen.Reference()

	host, err := rt.ResolveConsensusRef(ref)
	require.Error(t, err)
	require.Nil(t, host)
}
