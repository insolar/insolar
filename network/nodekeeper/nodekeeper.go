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

package nodekeeper

import (
	"github.com/insolar/insolar/core"
)

type NodeKeeper interface {
	// GetActiveNodes get active nodes.
	GetActiveNodes() []*core.ActiveNode
	// GetUnsyncHash get hash computed based on the list of unsync nodes, and the size of this list.
	GetUnsyncHash() (hash []byte, unsyncCount int)
	// GetUnsync gets the local unsync list (excluding other nodes unsync lists)
	GetUnsync() []*core.ActiveNode
	// Sync initiate transferring unsync -> sync, sync -> active. If approved is false, unsync is not transferred to sync
	Sync(approved bool)
	// AddUnsync add unsync node to the local unsync list
	AddUnsync(*core.ActiveNode)
	// AddUnsyncGossip merge unsync list from another node to the local unsync list
	AddUnsyncGossip([]*core.ActiveNode)
}
