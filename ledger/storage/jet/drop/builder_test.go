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

package drop

import (
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/ledger/storage/record"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/require"
)

func TestNewBuilder(t *testing.T) {
	t.Parallel()

	nb := NewBuilder(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher())
	cb, ok := nb.(*builder)

	require.Equal(t, true, ok)
	require.NotNil(t, cb.Hasher)
}

func TestBuilder_Append(t *testing.T) {
	t.Parallel()

	p := platformpolicy.NewPlatformCryptographyScheme()

	firstRecord := record.CodeRecord{
		Code: core.NewRecordID(123, []byte{76}),
	}
	secondRecord := record.ChildRecord{
		Ref: *core.NewRecordRef(core.RecordID{}, *core.NewRecordID(321, []byte{12})),
	}
	eh := p.ReferenceHasher()
	_, err := firstRecord.WriteHashData(eh)
	require.NoError(t, err)
	_, err = secondRecord.WriteHashData(eh)
	require.NoError(t, err)

	nb := NewBuilder(p.ReferenceHasher())

	err = nb.Append(&firstRecord)
	require.NoError(t, err)
	err = nb.Append(&secondRecord)
	require.NoError(t, err)

	cb := nb.(*builder)

	require.Equal(t, eh.Sum(nil), cb.Sum(nil))
}

func TestBuilder_Append_ReturnsErrIfAppendNil(t *testing.T) {
	t.Parallel()

	nb := NewBuilder(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher())

	err := nb.Append(nil)

	require.Error(t, err)
}

func TestBuilder_Size(t *testing.T) {
	t.Parallel()

	nb := NewBuilder(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher())
	cb := nb.(*builder)

	nb.Size(555)

	require.Equal(t, uint64(555), *cb.dropSize)
}

func TestBuilder_PrevHash(t *testing.T) {
	t.Parallel()

	nb := NewBuilder(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher())
	cb := nb.(*builder)

	nb.PrevHash([]byte{8, 7, 1, 2})

	require.Equal(t, []byte{8, 7, 1, 2}, cb.prevHash)
}

func TestBuilder_Pulse(t *testing.T) {
	t.Parallel()

	nb := NewBuilder(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher())
	cb := nb.(*builder)

	nb.Pulse(123)

	require.Equal(t, core.PulseNumber(123), *cb.pn)
}

func TestBuilder_Build_NilPn(t *testing.T) {
	t.Parallel()

	nb := NewBuilder(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher())
	nb.Size(555)
	nb.PrevHash([]byte{8, 7, 1, 2})

	_, err := nb.Build()

	require.Error(t, err)
}

func TestBuilder_Build_NilPrevHash(t *testing.T) {
	t.Parallel()

	nb := NewBuilder(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher())
	nb.PrevHash([]byte{8, 7, 1, 2})
	nb.Pulse(123)

	_, err := nb.Build()

	require.Error(t, err)
}

func TestBuilder_Build_NilSize(t *testing.T) {
	t.Parallel()

	nb := NewBuilder(platformpolicy.NewPlatformCryptographyScheme().ReferenceHasher())
	nb.Size(555)
	nb.Pulse(123)

	_, err := nb.Build()

	require.Error(t, err)
}

func TestBuilder_Build(t *testing.T) {
	p := platformpolicy.NewPlatformCryptographyScheme()

	firstRecord := record.CodeRecord{
		Code: core.NewRecordID(123, []byte{76}),
	}
	secondRecord := record.ChildRecord{
		Ref: *core.NewRecordRef(core.RecordID{}, *core.NewRecordID(321, []byte{12})),
	}
	eh := p.ReferenceHasher()
	_, err := firstRecord.WriteHashData(eh)
	require.NoError(t, err)
	_, err = secondRecord.WriteHashData(eh)
	require.NoError(t, err)

	nb := NewBuilder(p.ReferenceHasher())

	err = nb.Append(&firstRecord)
	require.NoError(t, err)
	err = nb.Append(&secondRecord)
	require.NoError(t, err)
	nb.Size(555)
	nb.Pulse(123)
	nb.PrevHash([]byte{8, 7, 1, 2})

	res, err := nb.Build()

	require.NoError(t, err)
	require.Equal(
		t,
		jet.Drop{
			Pulse:    123,
			PrevHash: []byte{8, 7, 1, 2},
			Hash:     eh.Sum(nil),
			DropSize: 555,
		},
		res,
	)
}
