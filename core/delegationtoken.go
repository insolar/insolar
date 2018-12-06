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

//go:generate minimock -i github.com/insolar/insolar/core.DelegationTokenFactory -o ../testutils -s _mock.go
type DelegationTokenFactory interface {
	IssuePendingExecution(msg Message, pulse PulseNumber) (DelegationToken, error)
	IssueGetObjectRedirect(sender *RecordRef, redirectedMessage Message) (DelegationToken, error)
	IssueGetChildrenRedirect(sender *RecordRef, redirectedMessage Message) (DelegationToken, error)
	IssueGetCodeRedirect(sender *RecordRef, redirectedMessage Message) (DelegationToken, error)
	Verify(parcel Parcel) (bool, error)
}

// DelegationToken is the base interface for delegation tokens
type DelegationToken interface {
	// Type returns token type.
	Type() DelegationTokenType

	// Verify checks against the token. See also delegationtoken.Verify(...)
	Verify(parcel Parcel) (bool, error)
}
