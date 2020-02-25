// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package jet

import (
	"context"
	"strconv"
	"strings"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bits"
)

//go:generate minimock -i github.com/insolar/insolar/insolar/jet.Accessor -o ./ -s _mock.go -g

// Accessor provides an interface for accessing jet IDs.
type Accessor interface {
	// All returns all jet from jet tree for provided pulse.
	All(ctx context.Context, pulse insolar.PulseNumber) []insolar.JetID
	// ForID finds jet in jet tree for provided pulse and object.
	// Always returns jet id and activity flag for this jet.
	ForID(ctx context.Context, pulse insolar.PulseNumber, recordID insolar.ID) (insolar.JetID, bool)
}

//go:generate minimock -i github.com/insolar/insolar/insolar/jet.Modifier -o ./ -s _mock.go -g

// Modifier provides an interface for modifying jet IDs.
type Modifier interface {
	// Update updates jet tree for specified pulse.
	Update(ctx context.Context, pulse insolar.PulseNumber, actual bool, ids ...insolar.JetID) error
	// Split performs jet split and returns resulting jet ids. Always set Active flag to true for leafs.
	Split(ctx context.Context, pulse insolar.PulseNumber, id insolar.JetID) (insolar.JetID, insolar.JetID, error)
	// Clone copies tree from one pulse to another. Use it to copy the past tree into new pulse.
	Clone(ctx context.Context, from, to insolar.PulseNumber, keepActual bool) error
}

//go:generate minimock -i github.com/insolar/insolar/insolar/jet.Cleaner -o ./ -s _mock.go -g

// Cleaner provides an interface for removing jet.Tree from a storage.
type Cleaner interface {
	// Delete jets for pulse (concurrent safe).
	DeleteForPN(ctx context.Context, pulse insolar.PulseNumber)
}

//go:generate minimock -i github.com/insolar/insolar/insolar/jet.Storage -o ./ -s _mock.go -g

// Storage composes Accessor and Modifier interfaces.
type Storage interface {
	Accessor
	Modifier
}

//go:generate minimock -i github.com/insolar/insolar/insolar/jet.Coordinator -o ./ -s _mock.go -g

// Coordinator provides methods for calculating Jet affinity
// (e.g. to which Jet a message should be sent).
type Coordinator interface {
	// Me returns current node.
	Me() insolar.Reference

	// IsAuthorized checks for role on concrete pulse for the address.
	IsAuthorized(ctx context.Context, role insolar.DynamicRole, obj insolar.ID, pulse insolar.PulseNumber, node insolar.Reference) (bool, error)

	// IsMeAuthorizedNow checks role of the current node in the current pulse for the address. Sugar for IsAuthorized.
	IsMeAuthorizedNow(ctx context.Context, role insolar.DynamicRole, obj insolar.ID) (bool, error)

	// QueryRole returns node refs responsible for role bound operations for given object and pulse.
	QueryRole(ctx context.Context, role insolar.DynamicRole, obj insolar.ID, pulse insolar.PulseNumber) ([]insolar.Reference, error)

	VirtualExecutorForObject(ctx context.Context, objID insolar.ID, pulse insolar.PulseNumber) (*insolar.Reference, error)
	VirtualValidatorsForObject(ctx context.Context, objID insolar.ID, pulse insolar.PulseNumber) ([]insolar.Reference, error)

	LightExecutorForObject(ctx context.Context, objID insolar.ID, pulse insolar.PulseNumber) (*insolar.Reference, error)
	LightValidatorsForObject(ctx context.Context, objID insolar.ID, pulse insolar.PulseNumber) ([]insolar.Reference, error)
	// LightExecutorForJet calculates light material executor for provided jet.
	LightExecutorForJet(ctx context.Context, jetID insolar.ID, pulse insolar.PulseNumber) (*insolar.Reference, error)
	LightValidatorsForJet(ctx context.Context, jetID insolar.ID, pulse insolar.PulseNumber) ([]insolar.Reference, error)

	Heavy(ctx context.Context) (*insolar.Reference, error)

	IsBeyondLimit(ctx context.Context, targetPN insolar.PulseNumber) (bool, error)
	NodeForJet(ctx context.Context, jetID insolar.ID, targetPN insolar.PulseNumber) (*insolar.Reference, error)

	// NodeForObject calculates a node (LME or heavy) for a specific jet for a specific pulseNumber
	NodeForObject(ctx context.Context, objectID insolar.ID, targetPN insolar.PulseNumber) (*insolar.Reference, error)
}

// Parent returns a parent of the jet or jet itself if depth of provided JetID is zero.
func Parent(id insolar.JetID) insolar.JetID {
	depth, prefix := id.Depth(), id.Prefix()
	if depth == 0 {
		return id
	}

	return *insolar.NewJetID(depth-1, bits.ResetBits(prefix, depth-1))
}

// NewIDFromString creates new JetID from string represents binary prefix.
//
// "0"     -> prefix=[0..0], depth=1
// "1"     -> prefix=[1..0], depth=1
// "1010"  -> prefix=[1010..0], depth=4
func NewIDFromString(s string) insolar.JetID {
	id := insolar.NewJetID(uint8(len(s)), parsePrefix(s))
	return *id
}

func parsePrefix(s string) []byte {
	var prefix []byte
	tail := s
	for len(tail) > 0 {
		offset := 8
		if len(tail) < 8 {
			tail += strings.Repeat("0", 8-len(tail))
		}
		parsed, err := strconv.ParseUint(tail[:offset], 2, 8)
		if err != nil {
			panic(err)
		}
		prefix = append(prefix, byte(parsed))
		tail = tail[offset:]
	}
	return prefix
}

// Siblings calculates left and right siblings for provided jet.
func Siblings(id insolar.JetID) (insolar.JetID, insolar.JetID) {
	depth, prefix := id.Depth(), id.Prefix()

	leftPrefix := bits.ResetBits(prefix, depth)
	left := insolar.NewJetID(depth+1, leftPrefix)

	rightPrefix := bits.ResetBits(prefix, depth)
	setBit(rightPrefix, depth)
	right := insolar.NewJetID(depth+1, rightPrefix)

	return *left, *right
}
