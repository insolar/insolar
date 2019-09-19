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

package sdk

type Member interface {
	GetReference() string
	GetPrivateKey() string
	GetPublicKey() string
}

// Member model object
type CommonMember struct {
	Reference  string
	PrivateKey string
	PublicKey  string
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
