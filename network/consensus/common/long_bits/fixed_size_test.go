package long_bits

import (
	"errors"
	io "io"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFoldUint64(t *testing.T) {
	require.Equal(t, uint32(0), FoldUint64(0))

	require.Equal(t, uint32(2), FoldUint64(2))

	require.Equal(t, uint32(math.MaxUint32), FoldUint64(math.MaxUint32))

	require.Equal(t, uint32(1), FoldUint64(math.MaxUint32+1))

	require.Equal(t, uint32(0), FoldUint64(math.MaxUint64))
}

func TestEqualFixedLenWriterTo(t *testing.T) {
	require.False(t, EqualFixedLenWriterTo(nil, nil))

	bits1 := NewBits64(0)
	require.False(t, EqualFixedLenWriterTo(&bits1, nil))

	require.False(t, EqualFixedLenWriterTo(nil, &bits1))

	bits2 := NewBits64(0)
	require.True(t, EqualFixedLenWriterTo(&bits1, &bits2))

	bits2 = NewBits64(1)
	require.False(t, EqualFixedLenWriterTo(&bits1, &bits2))
}

func TestCompare(t *testing.T) {
	require.False(t, (&writerToComparer{}).compare(nil, nil))

	bits1 := NewBits64(0)
	require.False(t, (&writerToComparer{}).compare(&bits1, nil))

	require.False(t, (&writerToComparer{}).compare(nil, &bits1))

	bits2 := NewBits64(0)
	require.True(t, (&writerToComparer{}).compare(&bits1, &bits2))

	bits3 := NewBits128(0, 0)
	require.False(t, (&writerToComparer{}).compare(&bits1, &bits3))

	bits1 = NewBits64(1)
	require.False(t, (&writerToComparer{}).compare(&bits1, &bits2))
}

func TestWrite(t *testing.T) {
	require.Panics(t, func() { (&writerToComparer{}).Write(nil) })

	bits := NewBits64(0)
	fr := NewFixedReaderMock(t)
	fr.WriteToMock.Set(func(io.Writer) (int64, error) { return 0, errors.New("test") })
	n, err := (&writerToComparer{other: fr}).Write(bits.AsBytes())
	require.NotEqual(t, nil, err)

	require.Equal(t, 0, n)

	n, err = (&writerToComparer{other: &bits}).Write(bits.AsBytes())
	require.Equal(t, nil, err)

	require.Equal(t, 8, n)
}

func TestAsByteString(t *testing.T) {
	fs := &fixedSize{}
	require.Equal(t, "", fs.AsByteString())

	fs = &fixedSize{data: []byte{'a', 'b', 'c'}}
	require.Equal(t, "abc", fs.AsByteString())
}

/*func TestWriteTo(t *testing.T) {
	fs := &fixedSize{data: []byte{}}

	bits := NewBits64(0)

	n, err := fs.WriteTo(&writerToComparer{other: &bits})
	require.True(t, err == nil)

	require.Equal(t, 8, n)
}*/

func TestRead(t *testing.T) {
	item := byte(3)
	fs := &fixedSize{data: []byte{item}}
	buf := make([]byte, 2)
	n, err := fs.Read(buf)
	require.Equal(t, 1, n)

	require.Equal(t, nil, err)

	require.Equal(t, item, buf[0])
}

func TestFoldToUint64(t *testing.T) {
	fs := &fixedSize{data: []byte{1}}
	require.Panics(t, func() { fs.FoldToUint64() })

	fs.data = append(fs.data, 2, 3, 4, 5, 6, 7, 8)
	require.Equal(t, uint64(0x807060504030201), fs.FoldToUint64())
}

func TestFixedByteSize(t *testing.T) {
	fs := &fixedSize{data: []byte{1, 2}}
	require.Equal(t, len(fs.data), fs.FixedByteSize())
}

func TestAsBytes(t *testing.T) {
	fs := &fixedSize{data: []byte{1, 2}}
	require.Len(t, fs.AsBytes(), len(fs.data))

	require.Equal(t, fs.data, fs.AsBytes())
}

func TestNewFixedReader(t *testing.T) {
	data := []byte{1, 2, 3}
	fr := NewFixedReader(data)
	require.Len(t, fr.AsBytes(), len(data))

	require.Equal(t, data[1], fr.AsBytes()[1])
}

func TestCopyFixedSize(t *testing.T) {
	item := 0x7777
	bits := NewBits64(uint64(item))
	fr := CopyFixedSize(&bits)

	require.Len(t, fr.AsBytes(), len(bits))

	require.Equal(t, uint8(item), fr.AsBytes()[0])

	require.Equal(t, bits[0], fr.AsBytes()[0])
}
