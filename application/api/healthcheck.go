// Copyright 2020 Insolar Network Ltd.
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

package api

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
)

// HealthChecker allows to check network status of a node.
type HealthChecker struct {
	CertificateManager insolar.CertificateManager
	NodeNetwork        network.NodeNetwork // nolint: staticcheck
	PulseAccessor      pulse.Accessor
}

// NewHealthChecker creates new HealthChecker.
func NewHealthChecker(cm insolar.CertificateManager, nn network.NodeNetwork, pa pulse.Accessor) *HealthChecker { // nolint: staticcheck
	return &HealthChecker{CertificateManager: cm, NodeNetwork: nn, PulseAccessor: pa}
}

// CheckHandler is a HTTP handler for health check.
func (hc *HealthChecker) CheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	ctx := r.Context()
	p, err := hc.PulseAccessor.Latest(ctx)
	if err != nil {
		err := errors.Wrap(err, "failed to get latest pulse")
		inslogger.FromContext(ctx).Errorf("[ NodeService.GetStatus ] %s", err.Error())
		_, _ = w.Write([]byte("FAIL"))
		return
	}
	for _, node := range hc.CertificateManager.GetCertificate().GetDiscoveryNodes() {
		if hc.NodeNetwork.GetAccessor(p.PulseNumber).GetWorkingNode(*node.GetNodeRef()) == nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("FAIL"))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
