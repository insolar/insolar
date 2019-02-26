/*
 *    Copyright 2019 Insolar
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

package drop

import (
	"io"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/pkg/errors"
)

type Builder interface {
	Append(item Hashable) error
	Size(size uint64)
	PrevHash(prevHash []byte)
	Pulse(pn core.PulseNumber)

	Build() (jet.JetDrop, error)
}

type Hashable interface {
	WriteHashData(w io.Writer) (int, error)
}

type builder struct {
	core.Hasher
	dropSize *uint64
	prevHash []byte
	pn       *core.PulseNumber
}

func NewBuilder(hasher core.Hasher) Builder {
	return &builder{
		Hasher: hasher,
	}
}

func (b *builder) Append(item Hashable) (err error) {
	_, err = item.WriteHashData(b.Hasher)
	return
}

func (b *builder) Size(size uint64) {
	b.dropSize = &size
}

func (b *builder) PrevHash(prevHash []byte) {
	b.prevHash = prevHash
}

func (b *builder) Pulse(pn core.PulseNumber) {
	b.pn = &pn
}

func (b *builder) Build() (jet.JetDrop, error) {
	if b.prevHash == nil {
		return jet.JetDrop{}, errors.New("prevHash is required")
	}
	if b.dropSize == nil {
		return jet.JetDrop{}, errors.New("dropSize is required")
	}
	if b.pn == nil {
		return jet.JetDrop{}, errors.New("pulseNumber is required")
	}

	return jet.JetDrop{
		Pulse:    *b.pn,
		PrevHash: b.prevHash,
		Hash:     b.Hasher.Sum(nil),
		DropSize: *b.dropSize,
	}, nil
}
