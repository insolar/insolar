package member

import (
	mfm "ilya/v2/mockMagic"
)

type Member struct {
	mfm.MockMagic
	name string
	publicKey []byte
}

func (m *Member) GetName() string {
	return m.name
}

func (m *Member) GetPublicKey() []byte {
	return m.publicKey
}
