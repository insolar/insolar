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

package object

import (
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is returned when value not found.
	ErrNotFound = errors.New("object not found")
	// ErrOverride is returned when trying to update existing record with the same id.
	ErrOverride = errors.New("record override is forbidden")

	// ErrIndexNotFound is returned when an index not found.
	ErrIndexNotFound = errors.New("index not found")
)
