package genesis

import (
	"github.com/insolar/insolar/applicationbase/builtin/contract/nodedomain"
	"github.com/insolar/insolar/applicationbase/genesisrefs"
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
