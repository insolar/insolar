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

package pulsenetwork

import (
	"context"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/transport/host"
)
type distributor struct {
	transport transport.Transport

	pulsarHost     *host.Host
	bootstrapHosts []*host.Host
}

func NewDistributor(conf configuration.PulseDistributor, transport transport.Transport) (core.PulseDistributor, error) {
	return &distributor{
	}, nil
}
func (d *distributor) Distribute(ctx context.Context, pulse *core.Pulse) {
	panic("not implemented")
}
