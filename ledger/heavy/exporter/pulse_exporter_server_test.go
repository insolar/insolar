/*
 *    Copyright 2019 Insolar Technologies
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

package exporter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

type pulseStreamMock struct {
	checker func(*Pulse) error
}

func (p *pulseStreamMock) Send(pulse *Pulse) error {
	return p.checker(pulse)
}

func (p *pulseStreamMock) SetHeader(metadata.MD) error {
	panic("implement me")
}

func (p *pulseStreamMock) SendHeader(metadata.MD) error {
	panic("implement me")
}

func (p *pulseStreamMock) SetTrailer(metadata.MD) {
	panic("implement me")
}

func (p *pulseStreamMock) Context() context.Context {
	panic("implement me")
}

func (p *pulseStreamMock) SendMsg(m interface{}) error {
	panic("implement me")
}

func (p *pulseStreamMock) RecvMsg(m interface{}) error {
	panic("implement me")
}

func TestPulseServer_Export(t *testing.T) {
	t.Run("fails if count is 0", func(t *testing.T) {
		server := NewPulseServer(nil, nil)

		err := server.Export(&GetPulses{Count: 0}, nil)

		require.NoError(t, err)
	})
}
