// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package api

import (
	"context"
	"errors"
	"github.com/gojuno/minimock/v3"
	"github.com/insolar/insolar/api/seedmanager"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/rpc/v2"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestSchemaService_Get(t *testing.T) {
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
	cfg := configuration.NewAPIRunner(false)
	cfg.Address = "localhost:19192"
	cfg.SwaggerPath = "testdata/api-exported.yaml"

	runner := Runner{
		cfg:                 &cfg,
		AvailabilityChecker: checker,
		PulseAccessor:       accessor,
		SeedManager:         seedmanager.New(),
		SeedGenerator:       seedmanager.SeedGenerator{},
	}
	s := NewSchemaService(&runner)
	defer runner.SeedManager.Stop()

	body := rpc.RequestBody{Raw: []byte(`{}`)}

	t.Run("success", func(t *testing.T) {
		availableFlag = true
		reply := new(interface{})

		err := s.Get(&http.Request{}, &SeedArgs{}, &body, reply)
		require.Nil(t, err)
		r := (*reply).(map[string]interface{})
		require.IsType(t, "string", r["openapi"], "openapi presented and it's string")
	})
}

func TestNewSchemaService(t *testing.T) {
	cfg := configuration.NewAPIRunner(false)
	cfg.SwaggerPath = "mamamylaramu"
	runner := Runner{
		cfg: &cfg,
	}

	defer func() {
		pan, ok := recover().(string)
		require.True(t, ok, "expected panic on unexistent file")
		require.Equal(t, "Can't read schema from 'mamamylaramu'", pan, "expected panic from file reading routine")
	}()
	NewSchemaService(&runner)
}
