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

package state

import (
	"github.com/insolar/insolar/network"
)

const (
	// NoNetwork state means that nodes doesn`t match majority_rule
	NoNetwork = iota
	// VoidNetwork state means that nodes have not complete min_role_count rule for proper work
	VoidNetwork
	// JetlessNetwork state means that every Jet need proof completeness of stored data
	JetlessNetwork
	// AuthorizationNetwork state means that every node need to validate ActiveNodeList using NodeDomain
	AuthorizationNetwork
	// CompleteNetwork state means network is ok and ready for proper work
	CompleteNetwork
)

type NetworkSwitcher struct {
}

func NewNetworkSwitcher() (*NetworkSwitcher, error) {
	return &NetworkSwitcher{}, nil
}

func (ns *NetworkSwitcher) GetState() network.State {
	return CompleteNetwork
}
