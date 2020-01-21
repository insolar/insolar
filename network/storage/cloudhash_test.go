// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

package storage

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryCloudHashStorage(t *testing.T) {
	cs := NewMemoryCloudHashStorage()

	pulse := insolar.Pulse{PulseNumber: 15}
	cloudHash := []byte{1, 2, 3, 4, 5}

	err := cs.Append(pulse.PulseNumber, cloudHash)
	assert.NoError(t, err)

	cloudHash2, err := cs.ForPulseNumber(pulse.PulseNumber)
	assert.NoError(t, err)

	assert.Equal(t, cloudHash, cloudHash2)
}
