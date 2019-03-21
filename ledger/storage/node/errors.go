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

package node

import (
	"github.com/pkg/errors"
)

var (
	// ErrOverride is returned when trying to set nodes for non-empty pulse.
	ErrOverride = errors.New("node override is forbidden")
	// ErrNoNodes is returned when nodes for specified criteria could not be found.
	ErrNoNodes = errors.New("matching nodes not found")
)
