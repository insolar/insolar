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
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/consensus"
	"github.com/pkg/errors"
)

type unsyncList struct {
	unsync []*core.ActiveNode
	pulse  core.PulseNumber
}

// NewUnsyncHolder create new object to hold data for consensus
func NewUnsyncHolder(pulse core.PulseNumber, unsync []*core.ActiveNode) consensus.UnsyncHolder {
	return &unsyncList{pulse: pulse, unsync: unsync}
}

// GetUnsync returns list of local unsync nodes. This list is created
func (u *unsyncList) GetUnsync() []*core.ActiveNode {
	return u.unsync
}

// GetPulse returns actual pulse for current consensus process.
func (u *unsyncList) GetPulse() core.PulseNumber {
	return u.pulse
}

// SetHash sets hash of unsync lists for each node of consensus.
func (u *unsyncList) SetHash([]*consensus.NodeUnsyncHash) {

}

// GetHash get hash of unsync lists for each node of consensus. If hash is not calculated yet, then this call blocks
// until the hash is calculated with SetHash() call
func (u *unsyncList) GetHash(blockTimeout time.Duration) ([]*consensus.NodeUnsyncHash, error) {
	return nil, errors.New("not implemented")
}
