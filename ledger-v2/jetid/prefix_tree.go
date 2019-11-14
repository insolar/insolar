//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package jetid

import (
	"fmt"
	"io"
	"math"
	"math/bits"
	"strings"

	"github.com/insolar/insolar/longbits"
)

type Prefix uint32

func NewPrefixTree(autoPropagate bool) PrefixTree {
	return PrefixTree{autoPropagate: autoPropagate, leafCounts: [17]uint16{0: 1}}
}

// limited to 65k Jets
type PrefixTree struct {
	lenNibles     [32768]uint8
	leafCounts    [17]uint16
	minDepth      uint8
	maxDepth      uint8
	autoPropagate bool
	mask          Prefix
}

func (p *PrefixTree) MaxDepth() uint8 {
	return p.maxDepth
}

func (p *PrefixTree) MinDepth() uint8 {
	return p.minDepth
}

func (p *PrefixTree) IsZero() bool {
	return p.minDepth == 0 && p.maxDepth == 0 && p.leafCounts[0] == 0
}

// it is not necessary to call it, but only initialized empty PrefixTree properly supports comparability
func (p *PrefixTree) Init() {
	if p.IsZero() {
		p.leafCounts[0] = 1
	}
}

func (p *PrefixTree) IsEmpty() bool {
	return p.minDepth == 0 && p.maxDepth == 0
}

func (p *PrefixTree) Count() int {
	if p.minDepth == p.maxDepth {
		return 1 << p.minDepth
	}
	total := 0
	for _, v := range p.leafCounts {
		total += int(v)
	}
	return total
}

func (p *PrefixTree) _getPrefixLength(prefix uint16) (uint8, bool) {
	depth := p.lenNibles[prefix>>1]
	if prefix&1 != 0 {
		depth >>= 4
	} else {
		depth &= 0x0F
	}
	switch {
	case depth != 0:
		return depth + 1, true
	case p.maxDepth == 0:
		return 0, prefix == 0
	default:
		return 1, prefix <= 1
	}
}

func (p *PrefixTree) getPrefixLength(prefix uint16) (uint8, bool) {
	switch depth, ok := p._getPrefixLength(prefix); {
	case !ok || depth == 0:
		return depth, ok
	case !p.autoPropagate:
		return depth, ok
	case bits.Len16(prefix) > int(depth):
		return depth, false
	default:
		return depth, ok
	}
}

func (p *PrefixTree) setPrefixLength(prefix uint16, depth uint8) {
	switch {
	case depth > 16:
		panic("illegal value")
	case depth != 0:
		depth--
	case prefix != 0:
		panic("illegal value")
	default:
		p.lenNibles[0] = 0
		return
	}

	idx := prefix >> 1
	d := p.lenNibles[idx]
	if prefix&1 != 0 {
		d = (d & 0x0F) | (depth << 4)
	} else {
		d = (d & 0xF0) | (depth & 0x0F)
	}
	p.lenNibles[idx] = d
}

func (p *PrefixTree) resetPrefixLength(prefix uint16) {
	if prefix&1 != 0 {
		p.lenNibles[prefix>>1] &= 0x0F
	} else {
		p.lenNibles[prefix>>1] &= 0xF0
	}
}

func (p *PrefixTree) FindPrefixLength(prefix Prefix) uint8 {
	_, l := p.findPrefixLength(prefix)
	return l
}

func (p *PrefixTree) findPrefixLength(prefix Prefix) (uint16, uint8) {
	maskedPrefix := uint16(prefix & p.mask)

	switch depth, ok := p._getPrefixLength(maskedPrefix); {
	case ok:
		return maskedPrefix, depth
	case p.autoPropagate:
		// with auto-propagation all entries must have a value
		panic("illegal state")
	}
	return p._findPrefixLength(maskedPrefix)
}

func (p *PrefixTree) _findPrefixLength(maskedPrefix uint16) (uint16, uint8) {
	for maskedPrefix > math.MaxUint8 {
		level := 7 + bits.Len8(uint8(maskedPrefix>>8))
		hiBit := uint16(1) << uint8(level)
		maskedPrefix ^= hiBit

		if depth, ok := p._getPrefixLength(maskedPrefix); ok {
			return maskedPrefix, depth
		}
	}

	for maskedPrefix > 0 {
		level := bits.Len8(uint8(maskedPrefix)) - 1
		hiBit := uint16(1) << uint8(level)
		maskedPrefix ^= hiBit

		if depth, ok := p._getPrefixLength(maskedPrefix); ok {
			return maskedPrefix, depth
		}
	}

	panic("illegal state")
}

func (p *PrefixTree) Split(prefix Prefix, prefixLimit uint8) {
	switch maskedPrefix, prefixLen := p.findPrefixLength(prefix); {
	case prefixLimit != prefixLen:
		panic("illegal value")
	case int(prefixLen) >= len(p.leafCounts):
		panic("illegal value") // TODO return as error?
	default:
		p._split(maskedPrefix, prefixLen, p.autoPropagate)
	}
}

func (p *PrefixTree) splitForDeserialize(maskedPrefix uint16, prefixLimit uint8) {
	switch prefixLen, ok := p.getPrefixLength(maskedPrefix); {
	case !ok:
		panic("illegal value")
	case prefixLen != prefixLimit:
		panic("illegal value")
	default:
		p._split(maskedPrefix, prefixLen, false)
	}
}

func (p *PrefixTree) _split(maskedPrefix uint16, prefixLen uint8, doPropagate bool) {
	switch n := p.leafCounts[prefixLen]; {
	case n > 1:
		p.leafCounts[prefixLen] = n - 1
	case n == 1:
		p.leafCounts[prefixLen] = 0
		if prefixLen == p.minDepth {
			p.minDepth++
		}
	case prefixLen == 0 && p.minDepth == 0 && p.maxDepth == 0:
		// zero state
		p.minDepth++
	default:
		panic("illegal state")
	}

	if prefixLen == p.maxDepth {
		p.maxDepth++
		p.mask = (p.mask << 1) | 1
		if p.autoPropagate {
			p.propagateNewDepth(p.maxDepth - 1)
		}
	}

	pairedPrefix := maskedPrefix | (1 << prefixLen)
	prefixLen++
	p.setPrefixLength(maskedPrefix, prefixLen)
	p.setPrefixLength(pairedPrefix, prefixLen)
	p.leafCounts[prefixLen] += 2

	if doPropagate {
		p.propagate(maskedPrefix, prefixLen)
		p.propagate(pairedPrefix, prefixLen)
	}
}

func (p *PrefixTree) Merge(prefix Prefix, prefixLimit uint8) {
	switch maskedPrefix, prefixLen := p.findPrefixLength(prefix); {
	case prefixLimit != prefixLen:
		panic("illegal value")
	case prefixLen == 0:
		panic("illegal value")
	default:
		p._merge(maskedPrefix, prefixLen)
	}
}

func (p *PrefixTree) merge(maskedPrefix uint16, prefixLimit uint8) {
	switch prefixLen, ok := p.getPrefixLength(maskedPrefix); {
	case !ok:
		panic("illegal value")
	case prefixLen != prefixLimit:
		panic("illegal value")
	case prefixLen == 0:
		panic("illegal value")
	default:
		p._merge(maskedPrefix, prefixLen)
	}
}

func (p *PrefixTree) _merge(maskedPrefix uint16, prefixLen uint8) {
	pairedPrefix := maskedPrefix | (1 << (prefixLen - 1))

	switch pairedPrefixLen, ok := p.getPrefixLength(pairedPrefix); {
	case maskedPrefix == pairedPrefix:
		panic("illegal value - only the zero-side is allowed to merge") // TODO return as error?
	case !ok:
		panic("illegal value - missing prefix") // TODO return as error?
	case pairedPrefixLen != prefixLen:
		panic("illegal value - unbalanced merge pair") // TODO return as error?
	}

	switch n := p.leafCounts[prefixLen]; {
	case n > 2:
		p.leafCounts[prefixLen] = n - 2
	case n == 2:
		switch p.maxDepth {
		case 0:
			panic("illegal state")
		case prefixLen:
			p.maxDepth--
			p.mask >>= 1
		}
		p.leafCounts[prefixLen] = 0
	default:
		panic("illegal state")
	}

	if prefixLen == p.minDepth {
		p.minDepth--
	}
	prefixLen--

	p.resetPrefixLength(pairedPrefix)
	p.resetPrefixLength(maskedPrefix)
	if prefixLen > 0 {
		p.setPrefixLength(maskedPrefix, prefixLen)
		if p.autoPropagate {
			p.propagate(maskedPrefix, prefixLen)
		}
	}
	p.leafCounts[prefixLen]++

}

func (p *PrefixTree) propagate(prefix uint16, baseDepth uint8) {
	switch {
	case baseDepth == 0 || baseDepth > p.maxDepth:
		panic("illegal state")
	case baseDepth == p.maxDepth:
		return
	}
	incStep := 1 << baseDepth
	if int(prefix) >= incStep {
		panic("illegal state")
	}
	setDepth := (baseDepth - 1) & 0xF
	maxStep := 1 << p.maxDepth
	if prefix&1 == 0 {
		for i := incStep; i < maxStep; i += incStep {
			idx := (i + int(prefix)) >> 1 // as we count nibles, not bytes
			p.lenNibles[idx] = p.lenNibles[idx]&0xF0 | setDepth
		}
	} else {
		setDepth <<= 4
		for i := incStep; i < maxStep; i += incStep {
			idx := (i + int(prefix)) >> 1 // as we count nibles, not bytes
			p.lenNibles[idx] = p.lenNibles[idx]&0x0F | setDepth
		}
	}
}

func (p *PrefixTree) propagateNewDepth(prevMaxDepth uint8) {
	switch {
	case p.maxDepth < prevMaxDepth:
		panic("illegal state")
	case p.maxDepth == prevMaxDepth:
		return
	case prevMaxDepth == 0:
		p.lenNibles[0] = p.lenNibles[0]&0x0F | p.lenNibles[0]<<4
		if p.maxDepth == 1 {
			return
		}
		fallthrough
	case prevMaxDepth == 1:
		for i := 1<<(p.maxDepth-1) - 1; i > 0; i-- {
			p.lenNibles[i] = p.lenNibles[0]
		}
		return
	}

	source := p.lenNibles[:1<<(prevMaxDepth-1)]
	max := 1 << (p.maxDepth - 1)
	for i := len(source); i < max; i += len(source) {
		copy(p.lenNibles[i:], source)
	}
}

func (p *PrefixTree) propagateAll() {
	if p.maxDepth <= 1 {
		return
	}

	max := uint16(1) << (p.maxDepth - 1)
	for i := uint16(1); i < max; i++ {
		v := p.lenNibles[i]
		switch {
		case v == 0:
			p.lenNibles[i] = p._getParentNible(i)
		case v <= 0x0F:
			p.lenNibles[i] = v | p._getParentNible(i)&0xF0
		case v&0x0F == 0:
			p.lenNibles[i] = v&0xF0 | p._getParentNible(i)&0x0F
		}
	}
}

func (p *PrefixTree) _getParentNible(index uint16) uint8 {
	highBit := uint16(1) << uint8(bits.Len16(index)-1)
	index ^= highBit
	return p.lenNibles[index]
}

func (p *PrefixTree) SetPropagate() {
	if p.autoPropagate {
		return
	}
	p.autoPropagate = true
	p.propagateAll()
}

func (p *PrefixTree) Cleanup() {
	switch {
	case p.maxDepth == 16:
		return
	case p.maxDepth == 0:
		if p.autoPropagate {
			for i := len(p.lenNibles) - 1; i >= 0; i-- {
				p.lenNibles[i] = 0
			}
		}
		switch p.leafCounts[0] {
		case 0:
			p.leafCounts[0] = 1
		case 1:
		default:
			panic("illegal state")
		}
	case p.maxDepth > 16:
		panic("illegal state")
	default:
		if p.autoPropagate {
			for i := 1 << (p.maxDepth - 1); i < len(p.lenNibles); i++ {
				p.lenNibles[i] = 0
			}
		}
		if p.leafCounts[0] != 0 {
			panic("illegal state")
		}
	}

	for i := len(p.leafCounts) - 1; i > int(p.maxDepth); i-- {
		if p.leafCounts[i] != 0 {
			panic("illegal state")
		}
	}
}

func (p *PrefixTree) String() string {
	return fmt.Sprintf("min=%d max=%d cnt=%v", p.minDepth, p.maxDepth, p.leafCounts)
}

func (p *PrefixTree) PrintTable() {
	p.printTable(p.getPrefixLength)
}

func (p *PrefixTree) PrintTableAll() {
	p.printTable(p._getPrefixLength)
}

func (p *PrefixTree) printTable(getFn func(uint16) (uint8, bool)) {
	fmt.Printf("min=%d max=%d cnt=%v\n", p.minDepth, p.maxDepth, p.leafCounts)
	for i := range p.lenNibles {
		prefix := uint16(i << 1)
		if depth, ok := getFn(prefix); ok {
			fmt.Printf("%5d [%2d]", prefix, depth)
			p.printRow(prefix, depth)
		}

		prefix++
		if depth, ok := getFn(prefix); ok {
			fmt.Printf("%5d [%2d]", prefix, depth)
			p.printRow(prefix, depth)
		}
	}
}

func (p *PrefixTree) printRow(prefix uint16, pLen uint8) {
	b := strings.Builder{}
	b.Grow(32)
	for i := uint8(0); i < pLen; i++ {
		b.WriteByte(' ')
		b.WriteByte('0' | byte(prefix)&1)
		prefix >>= 1
	}
	fmt.Println(b.String())
}

const compactSerializeV1 = 1

//
// General idea of this serialization is based on the "mountain range" approach to visualize Catalan numbers,
// yet it is different as we have 2 top and bottom limits and the left and right bounds can be above the bottom limit.
// https://en.wikipedia.org/wiki/Catalan_number
// https://brilliant.org/wiki/catalan-numbers/
//
// This implementation is suboptimal and consumes extra >40% of theoretical minimum (more for small trees),
// but it takes less for balanced trees (down to 2 bytes for a perfect tree).
//

func (p *PrefixTree) CompactSerialize(w io.Writer) error {
	b := p.CompactSerializeToBytes()
	switch n, e := w.Write(b); {
	case e != nil:
		return e
	case n != len(b):
		return fmt.Errorf("incomplete write")
	}
	return nil
}

const encodedDepthZeroTree = 0xFF

func (p *PrefixTree) CompactSerializeToBytes() []byte {
	encodedDepth := uint8(encodedDepthZeroTree)
	switch {
	case p.maxDepth < p.minDepth:
		panic("illegal state")
	case p.minDepth > 0:
		encodedDepth = p.minDepth - 1 | (p.maxDepth-p.minDepth)<<4
	}

	bb := longbits.NewBitBuilder(longbits.FirstLow, len(p.lenNibles))
	bb.AppendByte(compactSerializeV1)
	bb.AppendByte(encodedDepth)

	if p.maxDepth != p.minDepth {
		maxPrefix := 1 << p.minDepth
		for prefix := 0; prefix < maxPrefix; prefix++ {
			p.serializeBranch(&bb, uint16(prefix), p.minDepth)
		}
	}

	return bb.DoneToBytes()
}

const shallowBitCount = 3 // Meaningful values are 2 or 3. Factually disables use of shallow bit when =4

func (p *PrefixTree) serializeBranch(bb *longbits.BitBuilder, prefix uint16, minDepth uint8) {
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
		p.serializeBranch(bb, branchPrefix, branchDepth)
	}
}

func (p *PrefixTree) CompactDeserialize(r io.ByteReader) error {
	if p.maxDepth != 0 || p.minDepth != 0 {
		panic("illegal state")
	}
	switch b, e := r.ReadByte(); {
	case e != nil:
		return e
	case b != compactSerializeV1:
		return fmt.Errorf("unsupported type: %d", b)
	}

	switch encodedDepth, e := r.ReadByte(); {
	case e != nil:
		return e
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
	br := longbits.NewBitIoReader(longbits.FirstLow, r)

	maxPrefix := 1 << p.minDepth
	for prefix := 0; prefix < maxPrefix; prefix++ {
		if e := p.deserializeBranch(br, uint16(prefix), p.minDepth); e != nil {
			return e
		}
	}

	if p.autoPropagate {
		p.propagateAll()
	}

	return nil
}

func (p *PrefixTree) generatePrefectTree() {
	p.leafCounts[0] = 0
	upperBound := 1 << p.minDepth
	p.leafCounts[p.minDepth] = uint16(upperBound)
	upperBound >>= 1
	entry := (p.minDepth - 1) | (p.minDepth-1)<<4
	for i := 0; i < upperBound; i++ {
		p.lenNibles[i] = entry
	}
}

func (p *PrefixTree) deserializeBranch(br *longbits.BitIoReader, prefix uint16, minDepth uint8) error {
	maxDelta := p.maxDepth - minDepth
	switch {
	case p.maxDepth < minDepth:
		panic("illegal state")
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
		panic("illegal state")
	case depth < minDepth:
		panic("illegal state")
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
		if e := p.deserializeBranch(br, branchPrefix, branchDepth); e != nil {
			return e
		}
	}

	return nil
}
