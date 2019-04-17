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

	"github.com/insolar/insolar/insolar"
)

// HealthChecker allows to check network status of a node.
type HealthChecker struct {
	CertificateManager insolar.CertificateManager
	NodeNetwork        insolar.NodeNetwork
}

// NewHealthChecker creates new HealthChecker.
func NewHealthChecker(cm insolar.CertificateManager, nn insolar.NodeNetwork) *HealthChecker {
	return &HealthChecker{CertificateManager: cm, NodeNetwork: nn}
}

// CheckHandler is a HTTP handler for health check.
func (hc *HealthChecker) CheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	for _, node := range hc.CertificateManager.GetCertificate().GetDiscoveryNodes() {
		if hc.NodeNetwork.GetWorkingNode(*node.GetNodeRef()) == nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("FAIL"))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
