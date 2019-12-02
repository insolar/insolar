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

package dropbag

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/ledger-v2/dropbag/dbcommon"
)

func TestWriteRead(t *testing.T) {
	out := BufferWriter{}
	pw, err := OpenWriteStorage(&out, 1, 0)
	require.NoError(t, err)

	_, err = pw.WritePrelude([]byte{}, 0)
	require.NoError(t, err)

	_, err = pw.WriteChapter([]byte{}, dbcommon.ChapterDetails{1, 11111, 1})
	require.NoError(t, err)

	_, err = pw.WriteChapter([]byte{}, dbcommon.ChapterDetails{2, 222222, 2})
	require.NoError(t, err)

	_, err = pw.WriteChapter([]byte{}, dbcommon.ChapterDetails{3, 999999999, 11})
	require.NoError(t, err)

	_, err = pw.WriteConclude([]byte{})
	require.NoError(t, err)
	require.NoError(t, pw.Close())

	fmt.Println()
	fmt.Println(hex.Dump(out.Bytes()))

	in := NewBufferReader(out.Bytes())
	_, err = OpenReadStorage(&in, dbcommon.ReadConfig{ReadAllEntries: true}, payloadFactory{})
	require.NoError(t, err, "offset=%d", in.Offset())

	in = NewBufferReader(out.Bytes())
	_, err = OpenReadStorage(&in, dbcommon.ReadConfig{ReadAllEntries: false}, payloadFactory{})
	require.NoError(t, err, "offset=%d", in.Offset())
}

type payloadFactory struct {
}

func (payloadFactory) CreatePayloadBuilder(dbcommon.FileFormat, dbcommon.StorageSeqReader) dbcommon.PayloadBuilder {
	return payloadBuilder{}
}

type payloadBuilder struct {
}

func (p payloadBuilder) AddPrelude(bytes []byte, pos dbcommon.StorageEntryPosition) error {
	fmt.Printf("Prelude: len=%d pos=%v\n", len(bytes), pos)
	return nil
}

func (p payloadBuilder) AddConclude(bytes []byte, pos dbcommon.StorageEntryPosition, totalCount uint32) error {
	fmt.Printf("Conclude: len=%d pos=%v total=%d\n", len(bytes), pos, totalCount)
	return nil
}

func (p payloadBuilder) AddChapter(bytes []byte, pos dbcommon.StorageEntryPosition, details dbcommon.ChapterDetails) error {
	fmt.Printf("Chapter: len=%d pos=%v details=%+v\n", len(bytes), pos, details)
	return nil
}

func (p payloadBuilder) NeedsNextChapter() bool {
	return false
}

func (p payloadBuilder) Finished() error {
	return nil
}

func (p payloadBuilder) Failed(err error) error {
	return err
}

type BufferWriter struct {
	bytes.Buffer
}

func (b *BufferWriter) Offset() int64 {
	return int64(b.Len())
}

func NewBufferReader(buf []byte) BufferReader {
	return BufferReader{buf, bytes.NewReader(buf)}
}

type BufferReader struct {
	buf []byte
	*bytes.Reader
}

func (b *BufferReader) CanSeek() bool {
	return true
}

func (b *BufferReader) CanReadMapped() bool {
	return true
}

func (b *BufferReader) ReadMapped(n int64) ([]byte, error) {
	switch {
	case n > 0:
	case n == 0:
		return nil, nil
	default:
		panic("illegal value")
	}

	p := b.Offset()
	switch end, err := b.Seek(n, io.SeekCurrent); {
	case err != nil:
		return nil, err
	case end != p+n:
		return nil, io.ErrShortBuffer
	default:
		return b.buf[p:end], nil
	}

}

func (b *BufferReader) Offset() int64 {
	if n, err := b.Seek(0, io.SeekCurrent); err != nil {
		panic(err)
	} else {
		return n
	}
}
