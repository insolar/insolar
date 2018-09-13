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
	RoleVirtualExecutor  = JetRole(iota + 1) // Role responsible for current CPU operations
	RoleVirtualValidator                     // Role responsible for past CPU operations
	RoleLightExecutor                        // TODO: add docs
	RoleLightValidator                       // TODO: add docs
	RoleHeavyExecutor                        // TODO: add docs
)

type JetCoordinator interface {
	// IsAuthorized checks for role on concrete pulse for the address
	IsAuthorized(role JetRole, obj RecordRef, pulse PulseNumber, node RecordRef) bool

	// TODO: add docs
	QueryRole(role JetRole, obj RecordRef, pulse PulseNumber) []RecordRef
}
