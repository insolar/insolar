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
	"math"
	"math/bits"
	"strings"
)

type Prefix uint32

func NewPrefixTree(autoPropagate bool) PrefixTree {
	return PrefixTree{autoPropagate: autoPropagate, leafCounts: [17]uint16{0: 1}}
}

//
// Prefix tree for jets. Limited to 65536 jets and 16 bit prefix. Root is bit[0].
// The only difference with usual binary tree is that nodes are added and removed in pairs
// by Split and Merge operations accordingly.
//
// Empty PrefixTree always contains a root (zero) jet.
//
// PrefixTree has 2 modes varies by CPU (below, n = number of jets):
// - non-propagating: Split/Merge are O(1), GetPrefix() is O(log n)
// - propagating: GetPrefix() is O(1), Split/Merge is amortized O(log n) with peaks of O(n)
//
// Memory is always O(1) for any modes (fixed size structure, zero heap activity).
//
// PrefixTree can either be created NewPrefixTree() or as a PrefixTree{}.
// Zero PrefixTree will behave same as NewPrefixTree(false) but will return IsZero()=true until first modification.
//
// PrefixTree can be copied and compared as a structure.
// PrefixTree structures are equal when they have the same mode and same set of jets, irrelevant of Split/Merge sequence.
// Empty PrefixTree is not equal to a zero PrefixTree.
//
// Serialization is supported, with O(n log n) per operation.
//

type PrefixTree struct {
	lenNibles     [32768]uint8
	leafCounts    [17]uint16
	minDepth      uint8
	maxDepth      uint8
	autoPropagate bool
	mask          Prefix
}

// Maximum prefix length of a jet in this tree.
func (p *PrefixTree) MaxDepth() uint8 {
	return p.maxDepth
}

// Minimum prefix length of a jet in this tree.
func (p *PrefixTree) MinDepth() uint8 {
	return p.minDepth
}

func (p *PrefixTree) IsZero() bool {
	return p.minDepth == 0 && p.maxDepth == 0 && p.leafCounts[0] == 0
}

// Converts zero state tree to a proper empty tree.
// Only useful for the zero state, is not necessary to call.
func (p *PrefixTree) Init() {
	if p.IsZero() {
		p.leafCounts[0] = 1
	}
}

// True when there is only a root jet.
func (p *PrefixTree) IsEmpty() bool {
	return p.minDepth == 0 && p.maxDepth == 0
}

// Returns a total number of jets in this tree. Always >= 1. O(log n)
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
	// depth == 0 requires a special handling as it is a multi-purpose value
	switch {
	case depth != 0:
		return depth + 1, true
	case p.maxDepth == 0:
		return 0, prefix == 0
	default:
		return 1, prefix <= 1 || p.autoPropagate && p.minDepth == 1
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

// Returns number of denotative bits for the given prefix and masked prefix, with denotative bits only.
// Number of denotative bits is [0..16].
// O(log n) for non-propagating and O(1) for propagating mode.
func (p *PrefixTree) GetPrefix(prefix Prefix) (Prefix, uint8) {
	pfx, l := p.findPrefixLength(prefix)
	return Prefix(pfx), l
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

// Splits the given jet into 2 sub-jets (converts a leaf node into a full node with 2 leafs).
// (prefixLimit) - number of valuable bits in the given prefix. Will panic when prefixLimit is less than actual prefix length for the given prefix.
//
// O(1) for non-propagating and amortized O(log n) for propagating mode.
func (p *PrefixTree) Split(prefix Prefix, prefixLimit uint8) {
	switch maskedPrefix, prefixLen := p.findPrefixLength(prefix); {
	case prefixLimit < prefixLen:
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
		if doPropagate {
			p.propagateAllocatedDepth()
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

// Merges the given sub-jet with its pair into a jet (a full node with 2 leafs is converted into a leaf).
// (prefix) - must be zero-branch jet (has the highest denotative bit =0, or prefix[prefixLen]=0)
// (prefixLimit) - number of valuable bits in the given prefix. Will panic when prefixLimit is less than actual prefix length for the given prefix.
//
// O(1) for non-propagating and amortized O(log n) for propagating mode.
func (p *PrefixTree) Merge(prefix Prefix, prefixLimit uint8) {
	switch maskedPrefix, prefixLen := p.findPrefixLength(prefix); {
	case prefixLimit < prefixLen:
		panic("illegal value")
	case prefixLen == 0:
		panic("illegal value")
	default:
		p._merge(maskedPrefix, prefixLen, p.autoPropagate)
	}
}

//func (p *PrefixTree) merge(maskedPrefix uint16, prefixLimit uint8) {
//	switch prefixLen, ok := p.getPrefixLength(maskedPrefix); {
//	case !ok:
//		panic("illegal value")
//	case prefixLen != prefixLimit:
//		panic("illegal value")
//	case prefixLen == 0:
//		panic("illegal value")
//	default:
//		p._merge(maskedPrefix, prefixLen, false)
//	}
//}

func (p *PrefixTree) _merge(maskedPrefix uint16, prefixLen uint8, doPropagate bool) {
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
			p.cleanupReleasedDepth()
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
		if doPropagate {
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
	setDepth := baseDepth - 1
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

func (p *PrefixTree) propagateAllocatedDepth() {
	switch {
	case p.maxDepth == 0:
		panic("illegal state")
	case p.maxDepth <= 2:
		if p.lenNibles[0] != 0 {
			panic("illegal state")
		}
		if p.maxDepth == 2 {
			p.lenNibles[1] = 0
		}
		return
	}
	half := 1 << (p.maxDepth - 2)
	copy(p.lenNibles[half:], p.lenNibles[:half])
}

func (p *PrefixTree) cleanupReleasedDepth() {
	switch p.maxDepth {
	case 1:
		if p.lenNibles[1]&0xEE != 0 {
			panic("illegal state")
		}
		if p.lenNibles[0]&0xEE != 0 {
			panic("illegal state")
		}
		return
	case 0:
		if p.lenNibles[0] != 0 {
			panic("illegal state")
		}
		return
	}

	max := 1 << (p.maxDepth - 1)
	for i := max<<1 - 1; i >= max; i-- {
		p.lenNibles[i] = 0
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

// TODO remove?
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

// Prints a list of jets to StdOut
func (p *PrefixTree) PrintTable() {
	p.printTable(p.getPrefixLength)
}

// Prints a list of jets and propagated jets to StdOut
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
