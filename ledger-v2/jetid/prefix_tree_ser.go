package jetid

import (
	"bytes"
	"compress/lzw"
	"fmt"
	"github.com/insolar/insolar/longbits"
	"io"
	"math"
	"math/bits"
)

const (
	RawSerializeV1 = 1
	LZWSerializeV1 = 2
)

const compressionTolerance = math.MaxUint8 * 75 / 100 // 75% of uncompressed size
// const minSizeEfficientLZWLength = 160
const minEffortEfficientLZWLength = 256
const maxUncompressedSize = 8192

type PrefixTreeSerializer struct {
	UseLZW       bool
	LzwThreshold byte // =0 => minEffortEfficientLZWLength
	LzwTolerance byte // =0 => compressionTolerance
}

//
// General idea of this serialization is based on the "mountain range" approach to visualize Catalan numbers,
// yet it is different as we have 2 top and bottom limits and the left and right bounds can be above the bottom limit.
// https://en.wikipedia.org/wiki/Catalan_number
// https://brilliant.org/wiki/catalan-numbers/
//
// This implementation is suboptimal and consumes extra >40% of theoretical minimum,
// but it takes less for balanced trees (down to 2 bytes for a perfect tree).
//
// Approximate size of uncompressed serialized binary is 2 + 0.85*Count(), approx max is 6700 byte.
// When LZW compression is enabled, it can reduce by 2-3 times, approx max is down to 1400 bytes.
//
// First byte is either =RawSerializeV1 or =LZWSerializeV1
//
// O(n log n)
//
func (v PrefixTreeSerializer) Serialize(p *PrefixTree, w io.Writer) error {
	b := v.SerializeToRawBytes(p)
	return v.postSerialize(b, w)
}

func (v PrefixTreeSerializer) postSerialize(b []byte, w io.Writer) error {
	switch {
	case !v.UseLZW || len(b) <= 2:
		// don't try LZW
	case v.LzwThreshold == 0:
		if len(b) < minEffortEfficientLZWLength {
			break
		}
		fallthrough
	case len(b) >= int(v.LzwThreshold):
		b = v.tryLZW(b)
	}

	switch n, e := w.Write(b); {
	case e != nil:
		return e
	case n != len(b):
		return fmt.Errorf("incomplete write")
	}
	return nil
}

func (v PrefixTreeSerializer) tryLZW(uncompressed []byte) []byte {
	switch {
	case uncompressed[0] != RawSerializeV1:
		panic("illegal value")
	case len(uncompressed) <= 2:
		return uncompressed
	case len(uncompressed) > maxUncompressedSize:
		panic("illegal value")
	}

	compressed := bytes.NewBuffer(make([]byte, 0, len(uncompressed)+2))
	compressed.WriteByte(LZWSerializeV1)
	compressed.WriteByte(uncompressed[1])

	raw := uncompressed[2:] // ignore first 2 bytes

	// LittleEndian uint16
	compressed.WriteByte(byte(len(raw)))
	compressed.WriteByte(byte(len(raw) >> 8))

	compressor := lzw.NewWriter(compressed, lzw.LSB, 8)
	if _, err := compressor.Write(raw); err != nil {
		panic(err)
	}
	if err := compressor.Close(); err != nil {
		panic(err)
	}

	limit := compressionTolerance
	switch v.LzwTolerance {
	case 0:
	case math.MaxUint8:
		return compressed.Bytes()
	default:
		limit = int(v.LzwTolerance)
	}

	if compressed.Len() <= (limit*len(raw))/255 {
		return compressed.Bytes()
	}
	return uncompressed
}

const encodedDepthZeroTree = 0xFF

// Returns uncompressed form only
// First byte is always =RawSerializeV1
func (v PrefixTreeSerializer) SerializeToRawBytes(p *PrefixTree) []byte {
	encodedDepth := uint8(encodedDepthZeroTree)
	switch {
	case p.maxDepth < p.minDepth:
		panic("illegal state")
	case p.minDepth > 0:
		encodedDepth = p.minDepth - 1 | (p.maxDepth-p.minDepth)<<4
	}

	bb := longbits.NewBitBuilder(longbits.LSB, len(p.lenNibles))
	bb.AppendByte(RawSerializeV1)
	bb.AppendByte(encodedDepth)

	if p.maxDepth != p.minDepth {
		maxPrefix := 1 << p.minDepth
		for prefix := 0; prefix < maxPrefix; prefix++ {
			v.serializeBranch(p, &bb, uint16(prefix), p.minDepth)
		}
	}

	return bb.DoneToBytes()
}

const shallowBitCount = 3 // Meaningful values are 2 or 3. Factually disables use of shallow bit when =4

func (v PrefixTreeSerializer) serializeBranch(p *PrefixTree, bb *longbits.BitBuilder, prefix uint16, minDepth uint8) {
	depth, ok := p.getPrefixLength(prefix)
	maxDelta := p.maxDepth - minDepth
	//fmt.Printf("S: %04x %2d %2d %v\n", prefix, minDepth, depth, isShallow)
	switch {
	case !ok:
		panic("illegal state")
	case p.maxDepth < depth:
		panic("illegal state")
	case depth < minDepth:
		panic("illegal state")
	case maxDelta < 1<<shallowBitCount:
	case depth == minDepth:
		bb.AppendBit(0)
		return
	default:
		bb.AppendBit(1)
	}

	//fmt.Println(bb.AlignOffset(), bb.CompletedByteCount())
	bb.AppendSubByte(depth-minDepth, uint8(bits.Len8(maxDelta)))

	// zero-branch is accompanied by one-branches, one at each depth level
	for branchDepth := depth; branchDepth > minDepth; branchDepth-- {
		subBranchBit := uint16(1) << (branchDepth - 1)
		if prefix&subBranchBit != 0 { // TODO can this ever be true?
			continue
		}
		if branchDepth == p.maxDepth {
			continue
		}

		branchPrefix := prefix | subBranchBit
		v.serializeBranch(p, bb, branchPrefix, branchDepth)
	}
}

type PrefixTreeDeserializer struct {
}

// Reads the serialized content. Doesn't change propagation mode.
//
// O(n log n)
//

func (v PrefixTreeDeserializer) Deserialize(r io.ByteReader) (*PrefixTree, error) {
	tree := PrefixTree{}
	return &tree, v.deserializeTo(&tree, r)
}

// Can only be called on an empty tree otherwise panics.
func (v PrefixTreeDeserializer) DeserializeTo(p *PrefixTree, r io.ByteReader) error {
	if p.maxDepth != 0 || p.minDepth != 0 {
		panic("illegal state")
	}
	return v.deserializeTo(p, r)
}

func (v PrefixTreeDeserializer) unpackLZW(r io.ByteReader) (io.ByteReader, error) {
	if b0, err := r.ReadByte(); err != nil {
		return nil, err
	} else if b1, err := r.ReadByte(); err != nil {
		return nil, err
	} else {
		uncompressedSize := int(b0) | int(b1)<<8
		if uncompressedSize == 0 || uncompressedSize > maxUncompressedSize {
			return nil, fmt.Errorf("invalid content: format=lzw uncompressedSize=%d", uncompressedSize)
		}
		uncompressed := make([]byte, uncompressedSize)

		var wrapped io.Reader
		if ir, ok := r.(io.Reader); ok {
			wrapped = ir
		} else {
			// Only ByteReader is used inside lzw decoder, but declaration of NewReader doesn't allow to pass it
			wrapped = stubReader{r}
		}

		unpacker := lzw.NewReader(wrapped, lzw.LSB, 8)
		if _, err := io.ReadFull(unpacker, uncompressed); err != nil {
			return nil, err
		}
		if err := unpacker.Close(); err != nil {
			return nil, err
		}

		return bytes.NewReader(uncompressed), nil
	}
}

type stubReader struct {
	io.ByteReader
}

func (stubReader) Read(p []byte) (n int, err error) {
	panic("unexpected")
}

func (v PrefixTreeDeserializer) deserializeTo(p *PrefixTree, r io.ByteReader) error {

	var treeFn func() error
	switch b, err := r.ReadByte(); {
	case err != nil:
		return err
	case b == RawSerializeV1:
		treeFn = func() error {
			return v.deserializeTree(p, r)
		}
	case b == LZWSerializeV1:
		treeFn = func() error {
			if unpacked, err := v.unpackLZW(r); err != nil {
				return err
			} else {
				return v.deserializeTree(p, unpacked)
			}
		}
	default:
		return fmt.Errorf("unsupported type: %d", b)
	}

	switch encodedDepth, err := r.ReadByte(); {
	case err != nil:
		return err
	case encodedDepth == encodedDepthZeroTree:
		// empty tree
		p.leafCounts[0] = 1
		return nil
	default:
		p.minDepth = encodedDepth&0x0F + 1
		p.maxDepth = encodedDepth>>4 + p.minDepth
		if p.minDepth > p.maxDepth {
			return fmt.Errorf("invalid content: encodedDepth=%x", encodedDepth)
		}
		p.mask = (Prefix(1) << p.maxDepth) - 1

		p.generatePrefectTree()
		if p.minDepth == p.maxDepth {
			return nil
		}
	}

	return treeFn()
}

func (v PrefixTreeDeserializer) deserializeTree(p *PrefixTree, r io.ByteReader) error {
	br := longbits.NewBitIoReader(longbits.LSB, r)

	maxPrefix := 1 << p.minDepth
	for prefix := 0; prefix < maxPrefix; prefix++ {
		if err := v.deserializeBranch(p, br, uint16(prefix), p.minDepth); err != nil {
			return err
		}
	}

	if p.autoPropagate {
		p.propagateAll()
	}

	return nil
}

func (v PrefixTreeDeserializer) deserializeBranch(p *PrefixTree, br *longbits.BitIoReader, prefix uint16, minDepth uint8) error {
	maxDelta := p.maxDepth - minDepth
	switch {
	case p.maxDepth < minDepth:
		return fmt.Errorf("maxDepth < minDepth")
	case maxDelta < 1<<shallowBitCount:
	default:
		switch b, e := br.ReadBool(); {
		case e != nil:
			return e
		case !b:
			//fmt.Printf("D: %04x %2d -- %v\n", prefix, minDepth, isShallow)
			return nil
		}
	}

	depth := minDepth
	//fmt.Println(br.AlignOffset())
	if delta, e := br.ReadSubByte(uint8(bits.Len8(maxDelta))); e != nil {
		return e
	} else {
		depth += delta
	}
	switch {
	case depth > p.maxDepth:
		return fmt.Errorf("depth > p.maxDepth")
	case depth < minDepth:
		return fmt.Errorf("depth < minDepth")
	}
	//fmt.Printf("D: %04x %2d %2d %v\n", prefix, minDepth, depth, isShallow)

	// add a zero-branch and all relevant one-branches
	for i := minDepth; i < depth; i++ {
		p.splitForDeserialize(prefix, i)
	}

	// zero-branch is accompanied by one-branches, one at each depth level
	for branchDepth := depth; branchDepth > minDepth; branchDepth-- {
		subBranchBit := uint16(1) << (branchDepth - 1)
		if prefix&subBranchBit != 0 { // TODO can this ever be true?
			continue
		}
		if branchDepth == p.maxDepth {
			continue
		}

		branchPrefix := prefix | subBranchBit
		if e := v.deserializeBranch(p, br, branchPrefix, branchDepth); e != nil {
			return e
		}
	}

	return nil
}
