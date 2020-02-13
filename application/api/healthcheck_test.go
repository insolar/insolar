// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package api

import (
	"context"
	"net/http"
	"strconv"
	"testing"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/pulse"

	"github.com/insolar/insolar/insolar"
	network2 "github.com/insolar/insolar/network"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/network"

	"github.com/stretchr/testify/assert"
)

type mockResponseWriter struct {
	header http.Header
	body   *[]byte
}

func newMockResponseWriter() mockResponseWriter {
	return mockResponseWriter{
		header: make(http.Header, 0),
		body:   new([]byte),
	}
}

func (w mockResponseWriter) Header() http.Header {
	return w.header
}

func (w mockResponseWriter) Write(data []byte) (int, error) {
	*w.body = append(*w.body, data...)
	return len(data), nil
}

func (w mockResponseWriter) WriteHeader(statusCode int) {
	w.header["status"] = []string{strconv.Itoa(statusCode)}
}

func randomNodeList(t *testing.T, size int) []insolar.DiscoveryNode {
	list := make([]insolar.DiscoveryNode, size)
	for i := 0; i < size; i++ {
		dn := testutils.NewDiscoveryNodeMock(t)
		r := gen.Reference()
		dn.GetNodeRefMock.Set(func() *insolar.Reference {
			return &r
		})
		list[i] = dn
	}
	return list
}

func mockCertManager(t *testing.T, nodeList []insolar.DiscoveryNode) *testutils.CertificateManagerMock {
	cm := testutils.NewCertificateManagerMock(t)
	cm.GetCertificateMock.Set(func() insolar.Certificate {
		c := testutils.NewCertificateMock(t)
		c.GetDiscoveryNodesMock.Set(func() []insolar.DiscoveryNode {
			return nodeList
		})
		return c
	})
	return cm
}

func mockNodeNetwork(t *testing.T, nodeList []insolar.DiscoveryNode) *network.NodeNetworkMock {
	nn := network.NewNodeNetworkMock(t)
	nodeMap := make(map[insolar.Reference]insolar.DiscoveryNode)
	for _, node := range nodeList {
		nodeMap[*node.GetNodeRef()] = node
	}

	accessorMock := network.NewAccessorMock(t)
	accessorMock.GetWorkingNodeMock.Set(func(ref insolar.Reference) insolar.NetworkNode {
		if _, ok := nodeMap[ref]; ok {
			return network.NewNetworkNodeMock(t)
		}
		return nil
	})

	nn.GetAccessorMock.Set(func(p1 insolar.PulseNumber) network2.Accessor {
		return accessorMock
	})

	return nn
}

func mockPulseAccessor(t *testing.T) *pulse.AccessorMock {
	pa := pulse.NewAccessorMock(t)
	pa.LatestMock.Set(func(context.Context) (insolar.Pulse, error) {
		return *insolar.GenesisPulse, nil
	})
	return pa
}

func TestHealthChecker_CheckHandler(t *testing.T) {
	tests := []struct {
		name         string
		from, to     int
		status, body string
	}{
		{"all discovery", 0, 20, "200", "OK"},
		{"not enough discovery", 0, 11, "500", "FAIL"},
		{"extra nodes", 0, 40, "200", "OK"},
		{"not enough discovery and extra nodes", 5, 40, "500", "FAIL"},
	}
	nodes := randomNodeList(t, 40)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hc := NewHealthChecker(
				mockCertManager(t, nodes[:20]),
				mockNodeNetwork(t, nodes[test.from:test.to]),
				mockPulseAccessor(t),
			)
			w := newMockResponseWriter()
			hc.CheckHandler(w, new(http.Request))

			assert.Equal(t, w.header["status"], []string{test.status})
			assert.Equal(t, *w.body, []byte(test.body))
		})
	}
}
