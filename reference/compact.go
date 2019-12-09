//
// Copyright 2019 Insolar Technologies GmbH
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
//

package reference

func Empty() Holder {
	return emptyHolder
}

func EmptyLocal() *Local {
	return &emptyLocal
}

func NewRecord(local Local) Holder {
	if local.IsEmpty() {
		return Empty()
	}
	return NewNoCopy(&local, &emptyLocal)
}

func NewSelf(local Local) Holder {
	if local.IsEmpty() {
		return Empty()
	}
	return compact{&local, &local}
}

func New(local, base Local) Holder {
	return NewNoCopy(&local, &base)
}

func NewNoCopy(local, base *Local) Holder {
	switch {
	case local.IsEmpty():
		if base.IsEmpty() {
			return Empty()
		}
		local = &emptyLocal
	case base.IsEmpty():
		base = &emptyLocal
	}
	return compact{local, base}
}

var emptyLocal Local
var emptyHolder = compact{&emptyLocal, &emptyLocal}

type compact struct {
	addressLocal *Local
	addressBase  *Local
}

func (v compact) IsZero() bool {
	return v.addressLocal == nil
}

func (v compact) IsEmpty() bool {
	return v.addressLocal.IsEmpty() && v.addressBase.IsEmpty()
}

func (v compact) GetScope() Scope {
	return Scope(v.addressBase.getScope()<<2 | v.addressLocal.getScope())
}

func (v compact) GetBase() *Local {
	return v.addressBase
}

func (v compact) GetLocal() *Local {
	return v.addressLocal
}
