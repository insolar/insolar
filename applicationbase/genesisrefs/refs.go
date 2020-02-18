// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package genesisrefs

import (
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulse"
)

const (
	PrototypeType   = "prototype"
	PrototypeSuffix = "_proto"
)

// This constants are here temporary, after PENV-36 they will be moved
const (
	// GenesisNameNodeDomain is the name of node domain contract for genesis record.
	GenesisNameNodeDomain = "nodedomain"
	// GenesisNameNodeRecord is the name of node contract for genesis record.
	GenesisNameNodeRecord = "noderecord"
)

var PredefinedPrototypes = map[string]insolar.Reference{
	GenesisNameNodeDomain + PrototypeSuffix: *GenerateProtoReferenceFromContractID(PrototypeType, GenesisNameNodeDomain, 0),
	GenesisNameNodeRecord + PrototypeSuffix: *GenerateProtoReferenceFromContractID(PrototypeType, GenesisNameNodeRecord, 0),
}

var (
	// ContractNodeDomain is the node domain contract reference.
	ContractNodeDomain = GenesisRef(GenesisNameNodeDomain)
	// ContractNodeRecord is the node contract reference.
	ContractNodeRecord = GenesisRef(GenesisNameNodeRecord)
)

// Generate reference from hash code.
func GenerateProtoReferenceFromCode(pulse insolar.PulseNumber, code []byte) *insolar.Reference {
	hasher := platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher()
	codeHash := hasher.Hash(code)
	id := insolar.NewID(pulse, codeHash)
	return insolar.NewReference(*id)
}

// Generate prototype reference from contract id.
func GenerateProtoReferenceFromContractID(typeContractID string, name string, version int) *insolar.Reference {
	contractID := fmt.Sprintf("%s::%s::v%02d", typeContractID, name, version)
	return GenerateProtoReferenceFromCode(pulse.BuiltinContract, []byte(contractID))
}

// Generate contract reference from contract id.
func GenerateCodeReferenceFromContractID(typeContractID string, name string, version int) *insolar.Reference {
	contractID := fmt.Sprintf("%s::%s::v%02d", typeContractID, name, version)
	hasher := platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher()
	codeHash := hasher.Hash([]byte(contractID))
	id := insolar.NewID(pulse.BuiltinContract, codeHash)
	return insolar.NewRecordReference(*id)
}

// GenesisRef returns reference to any genesis records.
func GenesisRef(name string) insolar.Reference {
	if ref, ok := PredefinedPrototypes[name]; ok {
		return ref
	}
	pcs := platformpolicy.NewPlatformCryptographyScheme()
	req := record.IncomingRequest{
		CallType: record.CTGenesis,
		Method:   name,
	}
	virtRec := record.Wrap(&req)
	hash := record.HashVirtual(pcs.ReferenceHasher(), virtRec)
	id := insolar.NewID(pulse.MinTimePulse, hash)
	return *insolar.NewReference(*id)
}
