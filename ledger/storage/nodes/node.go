/*
 *    Copyright 2019 Insolar
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

package nodes

import (
	"crypto"

	"github.com/insolar/insolar/core"
)

type Node struct {
	FID   core.RecordRef
	FRole core.StaticRole
}

func (n Node) Address() string {
	panic("implement me")
}

func (n Node) ConsensusAddress() string {
	panic("implement me")
}

func (Node) GetGlobuleID() core.GlobuleID {
	panic("implement me")
}

func (n Node) ID() core.RecordRef {
	return n.FID
}

func (Node) PublicKey() crypto.PublicKey {
	panic("implement me")
}

func (n Node) Role() core.StaticRole {
	return n.FRole
}

func (Node) ShortID() core.ShortNodeID {
	panic("implement me")
}

func (Node) Version() string {
	panic("implement me")
}

func (Node) Leaving() bool {
	panic("implement me")
}

func (Node) LeavingETA() core.PulseNumber {
	panic("implement me")
}

func (Node) IsWorking() bool {
	panic("implement me")
}

func (Node) GetState() core.NodeState {
	panic("implement me")
}
