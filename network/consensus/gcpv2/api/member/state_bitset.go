package member

import (
	"fmt"
)

type StateBitset []BitsetEntry

func (b StateBitset) Len() int {
	return len(b)
}

type BitsetEntry uint8

const (
	BeHighTrust BitsetEntry = iota
	BeLimitedTrust
	BeBaselineTrust
	BeTimeout
	BeFraud
	maxBitsetEntry
)

const MaxBitsetEntry = int(maxBitsetEntry)

func (s BitsetEntry) IsTrusted() bool { return s < BeBaselineTrust }
func (s BitsetEntry) IsTimeout() bool { return s == BeTimeout }
func (s BitsetEntry) IsFraud() bool   { return s == BeFraud }

func (s BitsetEntry) String() string {
	return FmtBitsetEntry(uint8(s))
}

func FmtBitsetEntry(s uint8) string {
	switch BitsetEntry(s) {
	case BeHighTrust:
		return "H"
	case BeLimitedTrust:
		return "L"
	case BeBaselineTrust:
		return "B"
	case BeTimeout:
		return "Ø"
	case BeFraud:
		return "F"
	default:
		return fmt.Sprintf("?%d?", s)
	}
}
