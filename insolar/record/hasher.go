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

package record

import (
	"hash"

	"github.com/pkg/errors"
)

// HashVirtual returns hash for virtual record.
func HashVirtual(h hash.Hash, rec Virtual) []byte {
	// Signature must not affects material record hash calculating.
	rec.Signature = nil
	buf, err := rec.Marshal()
	if err != nil {
		panic(err)
	}
	_, err = h.Write(buf)
	if err != nil {
		panic(err)
	}
	return h.Sum(nil)
}

// HashMaterial returns hash for material record.
func HashMaterial(h hash.Hash, rec Material) ([]byte, error) {
	if rec.Virtual == nil {
		return nil, errors.New("virtual record is nil")
	}
	// Calculate virtual hash separately from the material
	// because changing material record fields must not affects
	// hash from virtual field.
	virtHash := HashVirtual(h, *rec.Virtual)
	rec.Virtual = nil

	// Signature must not affects material record hash calculating.
	rec.Signature = nil

	buf, err := rec.Marshal()
	if err != nil {
		panic(err)
	}

	// Appends virtual hash with other material fields.
	_, err = h.Write(virtHash)
	if err != nil {
		return nil, errors.Wrap(err, "can't write virtual-part record hash")
	}
	_, err = h.Write(buf)
	if err != nil {
		return nil, errors.Wrap(err, "can't write material-part record hash")
	}

	return h.Sum(nil), nil
}
