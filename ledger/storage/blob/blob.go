/*
 *    Copyright 2019 Insolar Technologies
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

package blob

import (
	"context"

	"github.com/insolar/insolar/core"
)

// Accessor provides info about Blob-values from storage
//go:generate minimock -i github.com/insolar/insolar/ledger/storage/blob.Accessor -o ./ -s _mock.go
type Accessor interface {
	// Get returns Blob for provided id
	Get(ctx context.Context, id core.RecordID) (Blob, error)
}

// Modifier provides provides methods for setting Blob-values to storage
//go:generate minimock -i github.com/insolar/insolar/ledger/storage/blob.Modifier -o ./ -s _mock.go
type Modifier interface {
	// Set saves new Blob-value in storage
	Set(ctx context.Context, id core.RecordID, blob Blob) error
}

type Blob struct {
	Value []byte
	JetID core.JetID
}
