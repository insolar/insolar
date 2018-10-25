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

package testutils

import (
	"crypto/rand"

	"github.com/insolar/insolar/core"
	"github.com/satori/go.uuid"
)

// RandomString generates random uuid and return it as a string
func RandomString() string {
	newUUID, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return newUUID.String()
}

// RandomRef generates random object reference
func RandomRef() core.RecordRef {
	ref := [core.RecordRefSize]byte{}
	_, err := rand.Read(ref[:])
	if err != nil {
		panic(err)
	}
	return ref
}

// RandomID generates random object ID
func RandomID() core.RecordID {
	id := [core.RecordIDSize]byte{}
	_, err := rand.Read(id[:])
	if err != nil {
		panic(err)
	}
	return id
}

func TestNode(ref core.RecordRef) *core.Node {
	return &core.Node{
		NodeID:   ref,
		PulseNum: core.PulseNumber(0),
		State:    core.NodeActive,
		Roles:    []core.NodeRole{core.RoleUnknown},
		// PublicKey: &key.PublicKey,
	}
}
