// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package coreapi

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"

	"github.com/stretchr/testify/require"
)

// TODO
/*func TestPickNextJoinCandidate(t *testing.T) {
	require.Equal(t, nil, (&SequentialCandidateFeeder{}).PickNextJoinCandidate())

	s := &SequentialCandidateFeeder{buf: make([]profiles.CandidateProfile, 1)}
	c := profiles.NewCandidateProfileMock(t)
	s.buf[0] = c
	require.Equal(t, c, s.PickNextJoinCandidate())
}*/

func TestRemoveJoinCandidate(t *testing.T) {
	require.False(t, (&SequentialCandidateFeeder{}).RemoveJoinCandidate(false, insolar.ShortNodeID(0)))

	s := &SequentialCandidateFeeder{buf: make([]profiles.CandidateProfile, 1)}
	c := profiles.NewCandidateProfileMock(t)

	s.buf[0] = c
	c.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return insolar.ShortNodeID(1) })
	require.False(t, s.RemoveJoinCandidate(false, insolar.ShortNodeID(2)))

	c.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return insolar.ShortNodeID(1) })
	require.True(t, s.RemoveJoinCandidate(false, insolar.ShortNodeID(1)))

	require.Equal(t, []profiles.CandidateProfile(nil), s.buf)

	s.buf = make([]profiles.CandidateProfile, 2)
	s.buf[0] = c
	c2 := profiles.NewCandidateProfileMock(t)
	s.buf[1] = c2
	require.True(t, s.RemoveJoinCandidate(false, insolar.ShortNodeID(1)))

	require.Equal(t, 1, len(s.buf))

	require.True(t, len(s.buf) > 0 && s.buf[0] == c2)
}

func TestAddJoinCandidatePanicForNil(t *testing.T) {
	s := NewSequentialCandidateFeeder(0)
	require.NotNil(t, s)
	require.Panics(t, func() { s.AddJoinCandidate(nil) })
}

func TestAddJoinCandidate(t *testing.T) {
	s := NewSequentialCandidateFeeder(0)

	f1 := transport.NewFullIntroductionReaderMock(t)
	f2 := transport.NewFullIntroductionReaderMock(t)

	err := s.AddJoinCandidate(f1)
	assert.NoError(t, err)
	require.True(t, len(s.buf) == 1 && s.buf[0] == f1)

	// add second
	err = s.AddJoinCandidate(f2)
	assert.NoError(t, err)
	require.True(t, len(s.buf) == 2 && s.buf[1] == f2)
}

func TestAddJoinCandidateFullQueue(t *testing.T) {
	s := NewSequentialCandidateFeeder(1)

	f1 := transport.NewFullIntroductionReaderMock(t)
	f2 := transport.NewFullIntroductionReaderMock(t)

	err := s.AddJoinCandidate(f1)
	assert.NoError(t, err)
	require.True(t, len(s.buf) == 1 && s.buf[0] == f1)

	// add second
	err = s.AddJoinCandidate(f2)
	assert.Error(t, err)
}
