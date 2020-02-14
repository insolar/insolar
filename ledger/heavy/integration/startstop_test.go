// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build slowtest

package integration_test

import (
	"context"
	"os"
	"testing"

	"github.com/insolar/insolar/applicationbase/genesis"
	"github.com/stretchr/testify/assert"
)

func TestStartStop(t *testing.T) {
	cfg := DefaultHeavyConfig()
	defer os.RemoveAll(cfg.Ledger.Storage.DataDirectory)
	heavyConfig := genesis.HeavyConfig{}
	s, err := NewBadgerServer(context.Background(), cfg, heavyConfig, nil)
	assert.NoError(t, err)
	s.Stop()
}
