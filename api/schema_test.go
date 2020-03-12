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
		reply := MAP{}

		err := s.Get(&http.Request{}, &SeedArgs{}, &body, &reply)
		require.Nil(t, err)
		require.IsType(t, "string", reply["openapi"], "right openapi")
	})
}

func TestNewSchemaService(t *testing.T) {
	cfg := configuration.NewAPIRunner(false)
	cfg.SwaggerPath = "mamamylaramu"
	runner := Runner{
		cfg: &cfg,
	}

	defer func() {
		require.NotNil(t, recover(), "panic on unexistent file")
	}()
	NewSchemaService(&runner)
}
