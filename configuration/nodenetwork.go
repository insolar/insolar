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

package configuration

// KeyPair holds public and private key for Node
type KeyPair struct {
	PublicKey  string
	PrivateKey string
}

// Node holds configuration for one Node
type Node struct {
	keys KeyPair
	role string
}

// NodeNetwork holds configuration for NodeNetwork
type NodeNetwork struct {
	Nodes []Node
}

// NewNodeNetwork creates new default NodeNetwork configuration
func NewNodeNetwork() NodeNetwork {
	return NodeNetwork{}
}
