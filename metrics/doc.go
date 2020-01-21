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

/*
Package metrics is a gateway for Prometheus monitoring system, it based on Prometheus golang client.
Package contains metrics collectors descriptions of entire project.
Component starts http server on http://0.0.0.0:8080/metrics by default(can be changed in configuration)

Example:

	// starts metrics server
	cfg := configuration.NewMetrics()
	m, _ := NewMetrics(cfg)
    m.Start(nil)

    // manipulate with network metrics
	NetworkMessageSentTotal.Inc()
	NetworkPacketSentTotal.WithLabelValues("ping").Add(55)

*/
package metrics
