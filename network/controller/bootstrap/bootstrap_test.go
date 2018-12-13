/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package bootstrap

import (
	"testing"
	"time"

	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/transport/host"
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

func mockBoostrap(string) (*host.Host, error) {
	return nil, BootstrapError
}

func mockInfinityBootstrap(string) (*host.Host, error) {
	bootstrapRetries++
	if bootstrapRetries >= 5 {
		return nil, nil
	}
	return nil, InfinityBootstrapError
}

func TestBootstrap(t *testing.T) {
	_, err := bootstrap("192.180.0.1:1234", getOptions(false), mockBoostrap)
	assert.Error(t, err, BootstrapError)

	startTime := time.Now()
	expectedTime := startTime.Add(time.Millisecond * 700) // 100ms, 200ms, 200ms, 200ms, return nil error
	_, err = bootstrap("192.180.0.1:1234", getOptions(true), mockInfinityBootstrap)
	endTime := time.Now()
	assert.NoError(t, err)
	assert.WithinDuration(t, expectedTime.Round(time.Millisecond), endTime.Round(time.Millisecond), time.Millisecond*100)
}
