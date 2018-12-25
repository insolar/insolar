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

package pulsar

import (
	"context"
	"sort"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/utils/entropy"
	"github.com/pkg/errors"
)

// FetchNeighbour searches neighbour of the pulsar by pubKey of a neighbout
func (currentPulsar *Pulsar) FetchNeighbour(pubKey string) (*Neighbour, error) {
	neighbour, ok := currentPulsar.Neighbours[pubKey]
	if !ok {
		return nil, errors.New("forbidden connection")
	}
	return neighbour, nil
}

// IsStateFailed checks if state of the pulsar is failed or not
func (currentPulsar *Pulsar) IsStateFailed() bool {
	return currentPulsar.StateSwitcher.GetState() == Failed
}

func (currentPulsar *Pulsar) isStandalone() bool {
	return len(currentPulsar.Neighbours) == 0
}

func (currentPulsar *Pulsar) getMaxTraitorsCount() int {
	nodes := len(currentPulsar.Neighbours) + 1
	return (nodes - 1) / 3
}

func (currentPulsar *Pulsar) getMinimumNonTraitorsCount() int {
	nodes := len(currentPulsar.Neighbours) + 1
	return nodes - currentPulsar.getMaxTraitorsCount()
}

func (currentPulsar *Pulsar) handleErrorState(ctx context.Context, err error) {
	inslogger.FromContext(ctx).Error(err)

	currentPulsar.clearState()
}

func (currentPulsar *Pulsar) clearState() {
	currentPulsar.SetGeneratedEntropy(nil)
	currentPulsar.GeneratedEntropySign = []byte{}

	currentPulsar.SetCurrentSlotEntropy(nil)
	currentPulsar.CurrentSlotPulseSender = ""

	currentPulsar.currentSlotSenderConfirmationsLock.Lock()
	currentPulsar.CurrentSlotSenderConfirmations = map[string]core.PulseSenderConfirmation{}
	currentPulsar.currentSlotSenderConfirmationsLock.Unlock()

	currentPulsar.ClearVector()
	currentPulsar.BftGridLock.Lock()
	currentPulsar.bftGrid = map[string]map[string]*BftCell{}
	currentPulsar.BftGridLock.Unlock()
}

func (currentPulsar *Pulsar) generateNewEntropyAndSign() error {
	e := currentPulsar.EntropyGenerator.GenerateEntropy()
	currentPulsar.SetGeneratedEntropy(&e)

	sign, err := currentPulsar.CryptographyService.Sign(currentPulsar.GetGeneratedEntropy()[:])
	if err != nil {
		return err
	}
	currentPulsar.GeneratedEntropySign = sign.Bytes()

	return nil
}

func (currentPulsar *Pulsar) preparePayload(body PayloadData) (*Payload, error) {
	hashProvider := currentPulsar.PlatformCryptographyScheme.IntegrityHasher()
	hash, err := body.Hash(hashProvider)
	if err != nil {
		return nil, err
	}
	sign, err := currentPulsar.CryptographyService.Sign(hash)
	if err != nil {
		return nil, err
	}

	return &Payload{Body: body, PublicKey: currentPulsar.PublicKeyRaw, Signature: sign.Bytes()}, nil
}

func (currentPulsar *Pulsar) checkPayloadSignature(request *Payload) (bool, error) {
	publicKey, err := currentPulsar.KeyProcessor.ImportPublicKeyPEM([]byte(request.PublicKey))
	if err != nil {
		return false, err
	}

	hash, err := request.Body.Hash(currentPulsar.PlatformCryptographyScheme.IntegrityHasher())
	if err != nil {
		return false, err
	}

	return currentPulsar.CryptographyService.Verify(publicKey, core.SignatureFromBytes(request.Signature), hash), nil
}

// copied from jetcoordinator
// (the only difference is type of input/output arrays)
func selectByEntropy(
	scheme core.PlatformCryptographyScheme,
	e core.Entropy,
	values []string,
	count int,
) ([]string, error) { // nolint: megacheck
	// TODO: remove sort when network provides sorted result from GetActiveNodesByRole (INS-890) - @nordicdyno 5.Dec.2018
	sort.Strings(values)
	in := make([]interface{}, 0, len(values))
	for _, value := range values {
		in = append(in, interface{}(value))
	}

	res, err := entropy.SelectByEntropy(scheme, e[:], in, count)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(res))
	for _, value := range res {
		out = append(out, value.(string))
	}
	return out, nil
}

// GetLastPulse returns last pulse in the thread-safe mode
func (currentPulsar *Pulsar) GetLastPulse() *core.Pulse {
	currentPulsar.lastPulseLock.RLock()
	defer currentPulsar.lastPulseLock.RUnlock()
	return currentPulsar.lastPulse
}

// SetLastPulse sets last pulse in the thread-safe mode
func (currentPulsar *Pulsar) SetLastPulse(newPulse *core.Pulse) {
	currentPulsar.lastPulseLock.Lock()
	defer currentPulsar.lastPulseLock.Unlock()
	currentPulsar.lastPulse = newPulse
}

// GetCurrentSlotEntropy returns currentSlotEntropy in the thread-safe mode
func (currentPulsar *Pulsar) GetCurrentSlotEntropy() *core.Entropy {
	currentPulsar.currentSlotEntropyLock.RLock()
	defer currentPulsar.currentSlotEntropyLock.RUnlock()
	return currentPulsar.currentSlotEntropy
}

// SetCurrentSlotEntropy sets currentSlotEntropy in the thread-safe mode
func (currentPulsar *Pulsar) SetCurrentSlotEntropy(currentSlotEntropy *core.Entropy) {
	currentPulsar.currentSlotEntropyLock.Lock()
	defer currentPulsar.currentSlotEntropyLock.Unlock()
	currentPulsar.currentSlotEntropy = currentSlotEntropy
}

// GetGeneratedEntropy returns generatedEntropy in the thread-safe mode
func (currentPulsar *Pulsar) GetGeneratedEntropy() *core.Entropy {
	currentPulsar.generatedEntropyLock.RLock()
	defer currentPulsar.generatedEntropyLock.RUnlock()
	return currentPulsar.generatedEntropy
}

// SetGeneratedEntropy sets generatedEntropy in the thread-safe mode
func (currentPulsar *Pulsar) SetGeneratedEntropy(currentSlotEntropy *core.Entropy) {
	currentPulsar.generatedEntropyLock.Lock()
	defer currentPulsar.generatedEntropyLock.Unlock()
	currentPulsar.generatedEntropy = currentSlotEntropy
}


func (currentPulsar *Pulsar) CreateVectorCopy() map[string]*BftCell {
	currentPulsar.ownedBtfRowLock.Lock()
	defer 	currentPulsar.ownedBtfRowLock.Unlock()

	newMap := map[string]*BftCell{}

	for key, value := range currentPulsar.ownedBftRow{
		newMap[key]= &BftCell{
			Entropy:value.GetEntropy(),
			IsEntropyReceived:value.GetIsEntropyReceived(),
			Sign:value.GetSign(),
		}
	}

	return newMap
}

func (currentPulsar *Pulsar) AddItemToVector(pubKey string, cell *BftCell){
	currentPulsar.ownedBtfRowLock.Lock()
	defer 	currentPulsar.ownedBtfRowLock.Unlock()
	currentPulsar.ownedBftRow[pubKey] = cell
}

func (currentPulsar *Pulsar) GetItemFromVector(pubKey string) (*BftCell, bool){
	currentPulsar.ownedBtfRowLock.RLock()
	defer 	currentPulsar.ownedBtfRowLock.RUnlock()
	btfCell, ok := currentPulsar.ownedBftRow[pubKey]
	return btfCell, ok
}

func (currentPulsar *Pulsar) ClearVector(){
	currentPulsar.ownedBtfRowLock.Lock()
	defer 	currentPulsar.ownedBtfRowLock.Unlock()
	currentPulsar.ownedBftRow = map[string]*BftCell{}
}