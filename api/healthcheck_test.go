//
// Copyright 2019 Insolar Technologies GmbH
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
//

package api

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/insolar/insolar/insolar"
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
		r := testutils.RandomRef()
		dn.GetNodeRefFunc = func() *insolar.Reference {
			return &r
		}
		list[i] = dn
	}
	return list
}

func mockCertManager(t *testing.T, nodeList []insolar.DiscoveryNode) *testutils.CertificateManagerMock {
	cm := testutils.NewCertificateManagerMock(t)
	cm.GetCertificateFunc = func() insolar.Certificate {
		c := testutils.NewCertificateMock(t)
		c.GetDiscoveryNodesFunc = func() []insolar.DiscoveryNode {
			return nodeList
		}
		return c
	}
	return cm
}

func mockNodeNetwork(t *testing.T, nodeList []insolar.DiscoveryNode) *network.NodeNetworkMock {
	nn := network.NewNodeNetworkMock(t)
	nodeMap := make(map[insolar.Reference]insolar.DiscoveryNode)
	for _, node := range nodeList {
		nodeMap[*node.GetNodeRef()] = node
	}
	nn.GetWorkingNodeFunc = func(ref insolar.Reference) insolar.NetworkNode {
		if _, ok := nodeMap[ref]; ok {
			return network.NewNetworkNodeMock(t)
		}
		return nil
	}
	return nn
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
			)
			w := newMockResponseWriter()
			hc.CheckHandler(w, new(http.Request))

			assert.Equal(t, w.header["status"], []string{test.status})
			assert.Equal(t, *w.body, []byte(test.body))
		})
	}
}
