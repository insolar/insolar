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

package configuration

// Metrics holds configuration for metrics publishing.
type Metrics struct {
	ListenAddress string
	Namespace     string
	ZpagesEnabled bool
}

// NewMetrics creates new default configuration for metrics publishing.
func NewMetrics() Metrics {
	return Metrics{
		ListenAddress: "0.0.0.0:9090",
		Namespace:     "insolar",
		ZpagesEnabled: true,
	}
}
