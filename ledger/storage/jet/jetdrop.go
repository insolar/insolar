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

package jet

import (
	"github.com/insolar/insolar/core"
)

// JetDrop is a blockchain block.
// It contains hashes of the current block and the previous one.
type JetDrop struct {
	// Pulse number (probably we should save it too).
	Pulse core.PulseNumber

	// PrevHash is a hash of all record hashes belongs to previous pulse.
	PrevHash []byte

	// Hash is a hash of all record hashes belongs to one pulse and previous drop hash.
	Hash []byte
}
