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

package object

import (
	"github.com/insolar/insolar/genesis/model/class"
)

// ReferenceContainer is a implementation of Proxy for containerization purpose.
type ReferenceContainer struct {
	reference Reference
	BaseProxy
}

// NewReferenceContainer creates new container for reference.
func NewReferenceContainer(ref Reference) *ReferenceContainer {
	return &ReferenceContainer{
		reference: ref,
	}
}

// GetClassID return string representation of object's class.
func (rc *ReferenceContainer) GetClassID() string {
	return class.ReferenceID
}

// GetStoredReference returns stored reference.
func (rc *ReferenceContainer) GetStoredReference() Reference {
	return rc.reference
}
