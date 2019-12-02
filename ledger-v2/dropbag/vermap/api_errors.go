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

package vermap

import "errors"

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrKeyNotFound    = errors.New("key not found")
	ErrEmptyKey       = errors.New("key cannot be empty")
	ErrExistingKey    = errors.New("existing key")
	ErrInvalidKey     = errors.New("key is restricted")
	ErrConflict       = errors.New("transaction conflict")
	ErrReadOnlyTxn    = errors.New("read-only transaction")
	ErrNoDelete       = errors.New("delete is not allowed")
	ErrDiscardedTxn   = errors.New("transaction has been discarded")
	ErrTxnTooBig      = errors.New("tx is too big")
)
