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

package insolar

import (
	"fmt"
	"testing"

	"github.com/insolar/insolar/reference"
)

func TestStub(t *testing.T) {
	str := "11tJDPHz1yWzKi4PoKybBDjLJmFeqH67qyKmwGECeMy.11111111111111111111111111111111"

	dec := reference.NewDefaultDecoder(reference.AllowLegacy | reference.AllowRecords)
	ref, err := dec.Decode(str)
	fmt.Println(ref)
	fmt.Println(err)
}
