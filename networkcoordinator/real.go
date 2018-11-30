/*
 *    Copyright 2018 INS Ecosystem
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

package networkcoordinator

import (
	"context"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

type realNetworkCoordinator struct {
}

func newRealNetworkCoordinator() *realNetworkCoordinator {
	return &realNetworkCoordinator{}
}

func (rnc *realNetworkCoordinator) GetCert(ctx context.Context, nodeRef core.RecordRef) (core.Certificate, error) {
	return nil, errors.New("not implemented")
}

func (rnc *realNetworkCoordinator) ValidateCert(ctx context.Context, certificate core.AuthorizationCertificate) (bool, error) {
	return false, errors.New("not implemented")
}

func (rnc *realNetworkCoordinator) WriteActiveNodes(ctx context.Context, number core.PulseNumber, activeNodes []core.Node) error {
	return errors.New("not implemented")
}

func (rnc *realNetworkCoordinator) SetPulse(ctx context.Context, pulse core.Pulse) error {
	return errors.New("not implemented")
}
