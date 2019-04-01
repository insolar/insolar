/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package bootstrap

import (
	"context"
	"testing"
	"time"

	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func getOptions(infinity bool) *common.Options {
	return &common.Options{
		TimeoutMult:       2 * time.Millisecond,
		InfinityBootstrap: infinity,
		MinTimeout:        100 * time.Millisecond,
		MaxTimeout:        200 * time.Millisecond,
		PingTimeout:       1 * time.Second,
		PacketTimeout:     10 * time.Second,
		BootstrapTimeout:  10 * time.Second,
	}
}

var BootstrapError = errors.New("bootstrap without repeat")
var InfinityBootstrapError = errors.New("infinity bootstrap")
var bootstrapRetries = 0

func mockBootstrap(context.Context, string) (*network.BootstrapResult, error) {
	return nil, BootstrapError
}

func mockInfinityBootstrap(context.Context, string) (*network.BootstrapResult, error) {
	bootstrapRetries++
	if bootstrapRetries >= 5 {
		return nil, nil
	}
	return nil, InfinityBootstrapError
}

func TestBootstrap(t *testing.T) {
	t.Skip("flaky test")
	ctx := context.Background()
	_, err := bootstrap(ctx, "192.180.0.1:1234", getOptions(false), mockBootstrap)
	assert.Error(t, err, BootstrapError)

	startTime := time.Now()
	expectedTime := startTime.Add(time.Millisecond * 700) // 100ms, 200ms, 200ms, 200ms, return nil error
	_, err = bootstrap(ctx, "192.180.0.1:1234", getOptions(true), mockInfinityBootstrap)
	endTime := time.Now()
	assert.NoError(t, err)
	assert.WithinDuration(t, expectedTime.Round(time.Millisecond), endTime.Round(time.Millisecond), time.Millisecond*100)
}
