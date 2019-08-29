//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package transport

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type testNode struct {
	udp     DatagramTransport
	address string
}

func (t *testNode) HandleDatagram(ctx context.Context, address string, buf []byte) {
	inslogger.FromContext(ctx).Info("Handle Datagram ", buf)
}

func newTestNode(port int) (*testNode, error) {
	cfg := configuration.NewHostNetwork().Transport
	cfg.Address = fmt.Sprintf("127.0.0.1:%d", port)

	node := &testNode{}
	udp, err := NewFactory(cfg).CreateDatagramTransport(node)
	if err != nil {
		return nil, err
	}
	node.udp = udp

	err = node.udp.Start(context.Background())
	if err != nil {
		return nil, err
	}

	node.address = udp.Address()
	return node, nil
}

func TestUdpTransport_SendDatagram(t *testing.T) {
	ctx := context.Background()

	node1, err := newTestNode(0)
	assert.NoError(t, err)
	node2, err := newTestNode(0)
	assert.NoError(t, err)

	err = node1.udp.SendDatagram(ctx, node2.address, []byte{1, 2, 3})
	assert.NoError(t, err)

	err = node2.udp.SendDatagram(ctx, node1.address, []byte{5, 4, 3})
	assert.NoError(t, err)

	err = node1.udp.Stop(ctx)
	assert.NoError(t, err)

	err = node1.udp.Start(ctx)
	assert.NoError(t, err)

	err = node1.udp.SendDatagram(ctx, node2.address, []byte{1, 2, 3})
	assert.NoError(t, err)

	err = node2.udp.SendDatagram(ctx, node1.address, []byte{5, 4, 3})
	assert.NoError(t, err)

	err = node1.udp.Stop(ctx)
	assert.NoError(t, err)
	err = node2.udp.Stop(ctx)
	assert.NoError(t, err)
}

func TestUdpTransport_SendDatagram_Error(t *testing.T) {
	cfg := configuration.NewHostNetwork().Transport
	cfg.Address = fmt.Sprintf("127.0.0.1:%d", 0)

	node := &testNode{}
	udp, err := NewFactory(cfg).CreateDatagramTransport(node)
	require.NoError(t, err)

	ctx := context.Background()
	err = udp.SendDatagram(ctx, udp.Address(), []byte{1, 2, 3})
	require.EqualError(t, err, "failed to send datagram: transport is not started")

	err = udp.Start(ctx)
	require.NoError(t, err)

	err = udp.SendDatagram(ctx, udp.Address(), []byte{1, 2, 3})
	require.NoError(t, err)

	err = udp.Stop(ctx)
	require.NoError(t, err)
}
