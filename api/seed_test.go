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

package api

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/insolar/insolar/api/seedmanager"
	"github.com/insolar/rpc/v2"
	"github.com/insolar/rpc/v2/json2"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/testutils"
)

func TestNodeService_GetSeed(t *testing.T) {
	defer testutils.LeakTester(t)

	availableFlag := false
	mc := minimock.NewController(t)
	checker := testutils.NewAvailabilityCheckerMock(mc)
	checker = checker.IsAvailableMock.Set(func(ctx context.Context) (b1 bool) {
		return availableFlag
	})

	// 0 = false, 1 = pulse.ErrNotFound, 2 = another error
	pulseError := 0
	accessor := mockPulseAccessor(t)
	accessor = accessor.LatestMock.Set(func(ctx context.Context) (p1 insolar.Pulse, err error) {
		switch pulseError {
		case 1:
			return insolar.Pulse{}, pulse.ErrNotFound
		case 2:
			return insolar.Pulse{}, errors.New("fake error")
		default:
			return *insolar.GenesisPulse, nil
		}
	})

	runner := Runner{
		AvailabilityChecker: checker,
		PulseAccessor:       accessor,
		SeedManager:         seedmanager.New(),
		SeedGenerator:       seedmanager.SeedGenerator{},
	}
	s := NewNodeService(&runner)
	defer runner.SeedManager.Stop()

	body := rpc.RequestBody{Raw: []byte(`{}`)}

	t.Run("success", func(t *testing.T) {
		availableFlag = true
		reply := requester.SeedReply{}

		err := s.GetSeed(&http.Request{}, &SeedArgs{}, &body, &reply)
		require.Nil(t, err)
		require.NotEmpty(t, reply.Seed)
	})
	t.Run("service not available", func(t *testing.T) {
		availableFlag = false

		err := s.GetSeed(&http.Request{}, &SeedArgs{}, &body, &requester.SeedReply{})
		require.Error(t, err)
		require.Equal(t, ServiceUnavailableErrorMessage, err.Error())
	})
	t.Run("pulse not found", func(t *testing.T) {
		availableFlag = true
		pulseError = 1

		err := s.GetSeed(&http.Request{}, &SeedArgs{}, &body, &requester.SeedReply{})
		require.Error(t, err)
		require.Equal(t, ServiceUnavailableErrorMessage, err.Error())
	})
	t.Run("pulse internal error", func(t *testing.T) {
		availableFlag = true
		pulseError = 2

		err := s.GetSeed(&http.Request{}, &SeedArgs{}, &body, &requester.SeedReply{})
		require.Error(t, err)
		require.Equal(t, InternalErrorMessage, err.Error())

		res, ok := err.(*json2.Error)
		require.True(t, ok)

		data, ok := res.Data.(requester.Data)
		require.True(t, ok)

		require.Equal(t, []string{"couldn't receive pulse", "fake error"}, data.Trace)
	})
}
