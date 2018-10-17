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

package certificate

import (
	"crypto/ecdsa"

	"github.com/insolar/insolar/core"
	ecdsahelper "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/pkg/errors"
)

// NewCertificate constructor creates new Certificate component
func NewCertificate(keysPath string) (*Certificate, error) {
	private, err := ecdsahelper.ImportPrivateKey("")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to import private key.")
	}

	public, err := ecdsahelper.ImportPublicKey("")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to import public key.")
	}

	return &Certificate{privateKey: private, publicKey: public}, nil
}

// Certificate component
type Certificate struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

// Start is method from Component interface and it do nothing
func (c *Certificate) Start(components core.Components) error {
	return nil
}

// Stop is method from Component interface and it do nothing
func (c *Certificate) Stop() error {
	return nil
}

// GetPublicKey returns public key as string
func (c *Certificate) GetPublicKey() (string, error) {
	return ecdsahelper.ExportPublicKey(c.publicKey)
}

// GetPublicKey returns private key as string
func (c *Certificate) GetPrivateKey() (string, error) {
	return ecdsahelper.ExportPrivateKey(c.privateKey)
}

// GetEcdsaPrivateKey returns private key in ecdsa format
func (c *Certificate) GetEcdsaPrivateKey() *ecdsa.PrivateKey {
	return c.privateKey
}
