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

package genesis

import (
	"github.com/insolar/insolar/application/genesisrefs"
	"github.com/insolar/insolar/applicationbase/builtin/contract/nodedomain"
	"github.com/insolar/insolar/insolar"
)

func NodeDomain(parentName string) ContractState {
	nd, _ := nodedomain.NewNodeDomain()
	return ContractState{
		Name:       genesisrefs.GenesisNameNodeDomain,
		Prototype:  genesisrefs.GenesisNameNodeDomain,
		ParentName: parentName,
		Memory:     MustGenMemory(nd),
	}
}

func MustGenMemory(data interface{}) []byte {
	b, err := insolar.Serialize(data)
	if err != nil {
		panic("failed to serialize contract data")
	}
	return b
}
