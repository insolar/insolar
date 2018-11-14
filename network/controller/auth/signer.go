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

package auth

import (
	"fmt"
	"sync"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

type Nonce []byte

// Signer is a component which signs and checks authorization nonce.
type Signer struct {
	lock                *sync.Mutex
	nonces              map[core.RecordRef]Nonce
	cryptographyService core.CryptographyService
	coordinator         core.NetworkCoordinator
}

// NewSignHandler creates a new sign handler.
func NewSigner(cryptographyService core.CryptographyService, coordinator core.NetworkCoordinator) *Signer {
	return &Signer{
		lock:                &sync.Mutex{},
		cryptographyService: cryptographyService,
		coordinator:         coordinator,
		nonces:              make(map[core.RecordRef]Nonce),
	}
}

// AddPendingNode adds a pending node and returns a corresponding nonce.
func (signer *Signer) AddPendingNode(ref core.RecordRef) (Nonce, error) {
	nonce, err := signer.generateNonce()
	if err != nil {
		return nil, errors.Wrapf(err, "Error generating nonce for node %s", ref.String())
	}
	signer.addNonce(ref, nonce)
	return nonce, nil
}

func (signer *Signer) generateNonce() (Nonce, error) {
	// TODO: add entropy.
	return time.Now().MarshalBinary()
}

func (signer *Signer) addNonce(ref core.RecordRef, nonce Nonce) {
	signer.lock.Lock()
	defer signer.lock.Unlock()

	// TODO: add cleaner for outdated nonces.
	signer.nonces[ref] = nonce
}

// AuthorizeNode checks a nonce sign and authorizes node.
func (signer *Signer) AuthorizeNode(ref core.RecordRef, signedNonce []byte) error {
	// ctx := context.TODO()

	_, ok := signer.getNonce(ref)
	if !ok {
		errMsg := fmt.Sprintf("Failed to authorize node %s: could not find previously sent nonce", ref.String())
		return errors.New(errMsg)
	}
	// TODO: make coordinator work on zeronet discovery node.
	// _, _, err := signer.coordinator.Authorize(ctx, ref, nonce, signedNonce)
	// if err != nil {
	// 	return errors.Wrapf(err, "Failed to authorize node %s", ref.String())
	// }
	return nil
}

func (signer *Signer) getNonce(ref core.RecordRef) (Nonce, bool) {
	signer.lock.Lock()
	defer signer.lock.Unlock()

	nonce, ok := signer.nonces[ref]
	if ok {
		delete(signer.nonces, ref)
	}
	return nonce, ok
}

// SignNonce sign a nonce.
func (signer *Signer) SignNonce(nonce Nonce) (*core.Signature, error) {
	sign, err := signer.cryptographyService.Sign(nonce)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to sign nonce")
	}
	return sign, nil
}
