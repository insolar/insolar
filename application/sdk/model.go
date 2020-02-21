// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package sdk

import (
	"fmt"
	"math/big"

	"github.com/insolar/insolar/api/requester"
)

type Member interface {
	GetReference() string
	GetPrivateKey() string
	GetPublicKey() string
	GetBalance() *big.Int
	SetBalance(*big.Int)
}

// Member model object
type CommonMember struct {
	Reference  string
	PrivateKey string
	PublicKey  string
	Balance    *big.Int
}

// MigrationMember model object
type MigrationMember struct {
	CommonMember
	MigrationAddress string
}

// NewMember creates new Member
func NewMember(ref string, privateKey string, publicKey string) *CommonMember {
	return &CommonMember{
		Reference:  ref,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}

func (m *CommonMember) GetReference() string {
	return m.Reference
}

func (m *CommonMember) GetPrivateKey() string {
	return m.PrivateKey
}

func (m *CommonMember) GetPublicKey() string {
	return m.PublicKey
}

func (m *CommonMember) GetBalance() *big.Int {
	return m.Balance
}

func (m *CommonMember) SetBalance(b *big.Int) {
	m.Balance = b
}

func (m *CommonMember) String() string {
	return fmt.Sprintf("Reference: %s; Private key: %s, Public key: %s. \n", m.Reference, m.PrivateKey, m.PublicKey)
}

// NewMigrationMember creates new MigrationMember
func NewMigrationMember(ref string, migrationAddress string, privateKey string, publicKey string) *MigrationMember {
	return &MigrationMember{
		CommonMember: CommonMember{
			Reference:  ref,
			PrivateKey: privateKey,
			PublicKey:  publicKey,
		},
		MigrationAddress: migrationAddress,
	}
}

func (m *MigrationMember) String() string {
	return fmt.Sprintf("Reference: %s; Private key: %s, Public key: %s, Migration address: %s. \n", m.Reference, m.PrivateKey, m.PublicKey, m.MigrationAddress)
}

// InfoResponse represents response from rpc on network.getInfo method
type InfoResponse struct {
	RootDomain             string   `json:"rootDomain"`
	RootMember             string   `json:"rootMember"`
	MigrationAdminMember   string   `json:"migrationAdminMember"`
	FeeMember              string   `json:"feeMember"`
	MigrationDaemonMembers []string `json:"migrationDaemonMembers"`
	NodeDomain             string   `json:"nodeDomain"`
	TraceID                string   `json:"traceID"`
}

type rpcInfoResponse struct {
	requester.Response
	Result InfoResponse `json:"result"`
}
