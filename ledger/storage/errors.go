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

package storage

import (
	"errors"

	"github.com/dgraph-io/badger"
)

var (
	// ErrNotFound returns if record/index not found in storage.
	ErrNotFound = errors.New("storage object not found")

	// ErrConflictRetriesOver is returned if Update transaction fails on all retry attempts.
	ErrConflictRetriesOver = errors.New("transaction conflict retries limit exceeded")

	// ErrConflict is the alias for badger.ErrConflict.
	ErrConflict = badger.ErrConflict

	// ErrOverride is returned if SetRecord tries to update existing record.
	ErrOverride = errors.New("records override is forbidden")

	// ErrClosed is returned when attempt to read or write to closed db.
	ErrClosed = errors.New("db is closed")
)
