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

package packet

import (
	"github.com/insolar/insolar/core"
)

// RequestPulse is data received from a pulsar.
type RequestPulse struct {
	Pulse core.Pulse
}

// RequestGetRandomHosts is data for the call that returns random hosts of the Insolar network.
type RequestGetRandomHosts struct {
	HostsNumber int
}
