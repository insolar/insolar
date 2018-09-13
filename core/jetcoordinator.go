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

package core

type JetRole int

const (
	NoRole = JetRole(iota)
	VirtualExecutor
	VirtualValidator
	HeavyExecutor
	LightExecutor
	LightValidator
)

type NetworkAddress string

type JetID RecordRef

type JetCoordinator interface {
	Component
	// AmI Checks Me for role on concrete pulse for this address
	AmI(role JetRole, ref RecordRef, number PulseNumber) bool
	IsIt(role JetRole, ref RecordRef, number PulseNumber) bool

	GetVirtualExecutor(pulse PulseNumber, ref RecordRef) NetworkAddress
	GetVirtualValidators(pulse PulseNumber, ref RecordRef) []NetworkAddress

	// TODO: depends on JetTree
	//GetJetID(ref RecordRef) JetID

	// TODO: calc JetID from RecordRef inside
	GetLightExecutor(pulse PulseNumber, ref RecordRef) NetworkAddress
	GetLightValidators(pulse PulseNumber, ref RecordRef) []NetworkAddress

	// TODO: calc JetID from RecordRef inside
	GetHeavyExecutor(pulse PulseNumber, ref RecordRef) NetworkAddress
}
