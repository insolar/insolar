// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package censusimpl

import (
	"testing"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"

	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/network/consensus/common/cryptkit"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
)

func TestNewManyNodePopulation(t *testing.T) {
	svf := cryptkit.NewSignatureVerifierFactoryMock(t)
	sv := cryptkit.NewSignatureVerifierMock(t)
	svf.CreateSignatureVerifierWithPKSMock.Set(func(cryptkit.PublicKeyStore) cryptkit.SignatureVerifier { return sv })
	require.Panics(t, func() { NewManyNodePopulation(nil, 0, nil) })

	sp := profiles.NewStaticProfileMock(t)
	pks := cryptkit.NewPublicKeyStoreMock(t)
	sp.GetPublicKeyStoreMock.Set(func() cryptkit.PublicKeyStore { return pks })
	nodeID := insolar.ShortNodeID(2)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return nodeID })
	sp.GetPrimaryRoleMock.Set(func() member.PrimaryRole { return member.PrimaryRoleNeutral })
	require.Panics(t, func() { NewManyNodePopulation([]profiles.StaticProfile{sp}, 0, nil) })

	require.Panics(t, func() { NewManyNodePopulation([]profiles.StaticProfile{sp}, nodeID+1, svf) })

	mnp := NewManyNodePopulation([]profiles.StaticProfile{sp}, nodeID, svf)
	require.NotNil(t, mnp.local)
}

func TestMNPString(t *testing.T) {
	mnp := ManyNodePopulation{}
	require.NotEmpty(t, mnp.String())

	mnp.isInvalid = true
	require.NotEmpty(t, mnp.String())

	us := updatableSlot{}
	mnp.local = &us
	sp := profiles.NewStaticProfileMock(t)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 1 })
	us.StaticProfile = sp
	require.NotEmpty(t, mnp.String())

	mnp.suspendedCount = 1
	require.NotEmpty(t, mnp.String())

	mnp.mistrustedCount = 1
	require.NotEmpty(t, mnp.String())

	mnp.slots = make([]updatableSlot, 2)
	require.NotEmpty(t, mnp.String())

	mnp.slots[1].index = member.JoinerIndex
	mnp.slots[1].StaticProfile = sp
	require.NotEmpty(t, mnp.String())

	mnp.slots[1].index = 2
	mnp.slots[1].mode = member.ModeEvictedGracefully
	require.NotEmpty(t, mnp.String())

	mnp.slots[1].mode = member.ModeEvictedAsFraud
	require.NotEmpty(t, mnp.String())

	mnp.slots[1].mode = member.ModeFlagValidationWarning
	require.NotEmpty(t, mnp.String())

	mnp.slots[1].mode = member.ModeFlagSuspendedOps
	require.NotEmpty(t, mnp.String())

	mnp.slots[1].mode = member.ModeNormal
	require.NotEmpty(t, mnp.String())

	mnp.slots = make([]updatableSlot, 51)
	require.NotEmpty(t, mnp.String())
}

func TestMNPGetSuspendedCount(t *testing.T) {
	suspendedCount := uint16(1)
	mnp := ManyNodePopulation{suspendedCount: suspendedCount}
	require.Equal(t, int(suspendedCount), mnp.GetSuspendedCount())
}

func TestMNPGetMistrustedCount(t *testing.T) {
	mistrustedCount := uint16(1)
	mnp := ManyNodePopulation{mistrustedCount: mistrustedCount}
	require.Equal(t, int(mistrustedCount), mnp.GetMistrustedCount())
}

func TestMNPGetIdleProfiles(t *testing.T) {
	mnp := ManyNodePopulation{}
	require.Nil(t, mnp.GetIdleProfiles())

	role := roleRecord{}
	roleCount := uint16(1)
	mnp.roles = make([]roleRecord, roleCount)
	mnp.roles[member.PrimaryRoleInactive] = role
	require.Panics(t, func() { mnp.GetIdleProfiles() })

	mnp.roles[member.PrimaryRoleInactive].container = &ManyNodePopulation{slots: make([]updatableSlot, roleCount)}
	require.Nil(t, mnp.GetIdleProfiles())

	mnp.roles[member.PrimaryRoleInactive].roleCount = roleCount
	require.Len(t, mnp.GetIdleProfiles(), int(roleCount))
}

func TestMNPGetIdleCount(t *testing.T) {
	mnp := ManyNodePopulation{}
	require.Zero(t, mnp.GetIdleCount())

	roleCount := uint16(1)
	role := roleRecord{roleCount: roleCount}
	mnp.roles = make([]roleRecord, roleCount)
	mnp.roles[member.PrimaryRoleInactive] = role
	require.Equal(t, int(roleCount), mnp.GetIdleCount())
}

func TestMNPGetIndexedCount(t *testing.T) {
	assignedSlotCount := uint16(1)
	mnp := ManyNodePopulation{assignedSlotCount: assignedSlotCount}
	require.Equal(t, int(assignedSlotCount), mnp.GetIndexedCount())
}

func TestMNPGetIndexedCapacity(t *testing.T) {
	size := 1
	mnp := ManyNodePopulation{slots: make([]updatableSlot, size)}
	require.Equal(t, size, mnp.GetIndexedCapacity())
}

func TestMNPIsValid(t *testing.T) {
	mnp := ManyNodePopulation{isInvalid: true}
	require.False(t, mnp.IsValid())

	mnp.isInvalid = false
	require.True(t, mnp.IsValid())
}

func TestMNPGetRolePopulation(t *testing.T) {
	mnp := ManyNodePopulation{}
	rolesCount := 2
	mnp.workingRoles = make([]member.PrimaryRole, rolesCount)
	require.Nil(t, mnp.GetRolePopulation(member.PrimaryRoleInactive))

	role := member.PrimaryRoleNeutral
	mnp.workingRoles = nil
	require.Nil(t, mnp.GetRolePopulation(role))

	mnp.workingRoles = make([]member.PrimaryRole, rolesCount)
	mnp.roles = make([]roleRecord, rolesCount)
	require.Nil(t, mnp.GetRolePopulation(role))

	mnp.roles[role].container = &ManyNodePopulation{}
	require.NotNil(t, mnp.GetRolePopulation(role))

	mnp.roles[role].container = nil
	mnp.roles[role].idleCount = 1
	require.NotNil(t, mnp.GetRolePopulation(role))
}

func TestMNPGetWorkingRoles(t *testing.T) {
	mnp := ManyNodePopulation{}
	require.Len(t, mnp.GetWorkingRoles(), 0)

	mnp.workingRoles = make([]member.PrimaryRole, 2)
	roleNumber := 1
	mnp.workingRoles[roleNumber] = member.PrimaryRoleNeutral
	require.Len(t, mnp.GetWorkingRoles(), len(mnp.workingRoles))

	require.Equal(t, mnp.workingRoles[roleNumber], mnp.GetWorkingRoles()[roleNumber])
}

func TestMNPCopyTo(t *testing.T) {
	sp := profiles.NewStaticProfileMock(t)
	nodeID := insolar.ShortNodeID(0)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return nodeID })
	mnp := &ManyNodePopulation{local: &updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}},
		slots: []updatableSlot{{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}}}}
	population := &DynamicPopulation{}
	mnp.copyTo(population)
	require.Equal(t, mnp.local, population.slotByID[nodeID])
}

func TestPanicOnRecoverable(t *testing.T) {
	require.Panics(t, func() { panicOnRecoverable(census.EmptySlot, "") })
}

func TestMakeCopyOfMapAndSeparateEvicts(t *testing.T) {
	mnp := ManyNodePopulation{}
	sp := profiles.NewStaticProfileMock(t)
	nodeID := insolar.ShortNodeID(0)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return *(&nodeID) })
	sp.GetPrimaryRoleMock.Set(func() member.PrimaryRole { return member.PrimaryRoleNeutral })
	spe := profiles.NewStaticProfileExtensionMock(t)
	sp.GetExtensionMock.Set(func() profiles.StaticProfileExtension { return spe })
	local := &updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}}
	require.Panics(t, func() { mnp.makeCopyOfMapAndSeparateEvicts(nil, local, nil) })

	slots := make(map[insolar.ShortNodeID]*updatableSlot)
	slots[nodeID] = &updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp}}
	require.Panics(t, func() { mnp.makeCopyOfMapAndSeparateEvicts(slots, local, nil) })

	delete(slots, nodeID)
	nodeID = 1
	slots[nodeID] = &updatableSlot{NodeProfileSlot: NodeProfileSlot{StaticProfile: sp, mode: member.ModeEvictedGracefully}}
	mnp.slotByID = make(map[insolar.ShortNodeID]*updatableSlot)
	mnp.slotByID[nodeID] = slots[nodeID]
	mnp.assignedSlotCount = 1
	require.Len(t, mnp.makeCopyOfMapAndSeparateEvicts(slots, local, nil), 1)
}

func TestFilterAndFillInSlots(t *testing.T) {
	mnp := ManyNodePopulation{}
	slots := make(map[insolar.ShortNodeID]*updatableSlot, member.MaxNodeIndex+1)
	for i := insolar.ShortNodeID(0); i <= member.MaxNodeIndex; i++ {
		slots[i] = nil
	}
	require.Panics(t, func() { mnp._filterAndFillInSlots(slots, panicOnRecoverable) })

	slots = make(map[insolar.ShortNodeID]*updatableSlot)
	slots[1] = nil
	require.Panics(t, func() { mnp._filterAndFillInSlots(slots, panicOnRecoverable) })

	slots[1] = &updatableSlot{}
	require.Panics(t, func() { mnp._filterAndFillInSlots(slots, panicOnRecoverable) })

	delete(slots, 1)
	sp := profiles.NewStaticProfileMock(t)
	slots[insolar.AbsentShortNodeID] = &updatableSlot{}
	slots[insolar.AbsentShortNodeID].StaticProfile = sp
	require.Panics(t, func() { mnp._filterAndFillInSlots(slots, panicOnRecoverable) })

	delete(slots, 0)
	us := &updatableSlot{}
	slots[1] = us
	us.StaticProfile = sp
	role := member.PrimaryRoleInactive
	sp.GetPrimaryRoleMock.Set(func() member.PrimaryRole { return *(&role) })
	require.Panics(t, func() { mnp._filterAndFillInSlots(slots, panicOnRecoverable) })

	role = member.PrimaryRoleNeutral
	us.index = member.JoinerIndex
	require.Panics(t, func() { mnp._filterAndFillInSlots(slots, panicOnRecoverable) })

	us.index = member.Index(1)
	us.mode = member.ModeEvictedGracefully
	evicts, slotCount := mnp._filterAndFillInSlots(slots, panicOnRecoverable)
	require.Len(t, evicts, 1)

	require.Zero(t, slotCount)

	us.mode = member.ModeRestrictedAnnouncement
	require.Panics(t, func() { mnp._filterAndFillInSlots(slots, panicOnRecoverable) })

	us.index = member.Index(0)
	sp.GetExtensionMock.Set(func() profiles.StaticProfileExtension { return nil })
	require.Panics(t, func() { mnp._filterAndFillInSlots(slots, panicOnRecoverable) })

	spe := profiles.NewStaticProfileExtensionMock(t)
	sp.GetExtensionMock.Set(func() profiles.StaticProfileExtension { return spe })
	slots[2] = us
	require.Panics(t, func() { mnp._filterAndFillInSlots(slots, panicOnRecoverable) })

	us2 := &updatableSlot{}
	slots[2] = us2
	us2.index = member.Index(1)
	us2.mode = member.ModeRestrictedAnnouncement
	us2.StaticProfile = sp
	us.mode = member.ModeEvictedGracefully
	evicts, slotCount = mnp._filterAndFillInSlots(slots, panicOnRecoverable)
	require.Len(t, evicts, 1)

	require.Equal(t, 1, slotCount)
}

func doNothingOnRecoverable(census.RecoverableErrorTypes, string, ...interface{}) {
	// Do nothing.
}

func TestFillInRoleStatsAndMap(t *testing.T) {
	mnp := ManyNodePopulation{}
	localID := insolar.ShortNodeID(0)
	slotCount := 2
	compactIndex := false
	checkUniqueID := false
	fail := panicOnRecoverable
	require.Panics(t, func() {
		mnp._fillInRoleStatsAndMap(localID, member.MaxNodeIndex+1, compactIndex, checkUniqueID, fail)
	})

	require.Panics(t, func() { mnp._fillInRoleStatsAndMap(localID, slotCount, compactIndex, checkUniqueID, fail) })

	mnp = ManyNodePopulation{}
	mnp.slots = make([]updatableSlot, 1)
	mnp.slots[0] = updatableSlot{}
	mnp.slotByID = make(map[insolar.ShortNodeID]*updatableSlot, slotCount)
	mnp.slotByID[localID] = &mnp.slots[0]
	slotCount = 0
	mnp._fillInRoleStatsAndMap(localID, slotCount, compactIndex, checkUniqueID, fail)
	require.False(t, mnp.isInvalid)

	slotCount = 2
	mnp._fillInRoleStatsAndMap(localID, slotCount, compactIndex, checkUniqueID, fail)
	require.False(t, mnp.isInvalid)

	mnp._fillInRoleStatsAndMap(localID, slotCount, !compactIndex, checkUniqueID, fail)
	require.False(t, mnp.isInvalid)

	sp := profiles.NewStaticProfileMock(t)
	mnp.slots[0].StaticProfile = sp
	nodeID := insolar.ShortNodeID(0)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return *(&nodeID) })
	role := member.PrimaryRoleNeutral
	sp.GetPrimaryRoleMock.Set(func() member.PrimaryRole { return *(&role) })
	mnp._fillInRoleStatsAndMap(localID, slotCount, compactIndex, checkUniqueID, fail)
	require.False(t, mnp.isInvalid)

	mnp.slots = append(mnp.slots, updatableSlot{})
	mnp.slots[1].StaticProfile = sp
	mnp.slots[0].StaticProfile = nil
	mnp._fillInRoleStatsAndMap(localID, slotCount, !compactIndex, checkUniqueID, fail)
	require.False(t, mnp.isInvalid)

	require.Panics(t, func() { mnp._fillInRoleStatsAndMap(localID, slotCount, compactIndex, !checkUniqueID, fail) })

	role = member.PrimaryRoleInactive
	require.Panics(t, func() { mnp._fillInRoleStatsAndMap(localID, slotCount, compactIndex, checkUniqueID, fail) })

	role = member.PrimaryRoleNeutral
	mnp.slots[0].power = 1
	mnp.slots[0].mode = member.ModeEvictedGracefully
	mnp._fillInRoleStatsAndMap(localID, slotCount, compactIndex, checkUniqueID, fail)
	require.True(t, mnp.isInvalid)

	mnp.slots[0].mode = member.ModeFlagValidationWarning
	mnp._fillInRoleStatsAndMap(localID, slotCount, compactIndex, checkUniqueID, fail)
	require.True(t, mnp.isInvalid)

	role = member.PrimaryRoleInactive
	mnp.slots[0].mode = member.ModeNormal
	mnp._fillInRoleStatsAndMap(localID, slotCount, compactIndex, checkUniqueID, doNothingOnRecoverable)
	require.True(t, mnp.isInvalid)

	mnp.slots[0].StaticProfile = sp
	sp2 := profiles.NewStaticProfileMock(t)
	mnp.slots[1].StaticProfile = sp2
	sp2.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	role2 := member.PrimaryRoleHeavyMaterial
	sp2.GetPrimaryRoleMock.Set(func() member.PrimaryRole { return *(&role2) })
	mnp.slots[1].power = 1
	mnp._fillInRoleStatsAndMap(localID, slotCount, compactIndex, checkUniqueID, doNothingOnRecoverable)
	require.True(t, mnp.isInvalid)
}

func TestAdjustSlotsAndCopyEvicts(t *testing.T) {
	mnp := ManyNodePopulation{}
	localID := insolar.ShortNodeID(1)
	us := &updatableSlot{}
	sp := profiles.NewStaticProfileMock(t)
	nodeID := insolar.ShortNodeID(0)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return *(&nodeID) })
	us.StaticProfile = sp
	evicts := []*updatableSlot{us}
	require.Len(t, mnp._adjustSlotsAndCopyEvicts(localID, evicts), len(evicts))

	require.Nil(t, mnp.local)

	mnp.slots = make([]updatableSlot, 1)
	mnp.slots[0].StaticProfile = sp
	mnp._adjustSlotsAndCopyEvicts(localID, evicts)
	require.Len(t, mnp._adjustSlotsAndCopyEvicts(localID, evicts), len(evicts))

	mnp.slots = append(mnp.slots, updatableSlot{})
	mnp.assignedSlotCount = 1
	nodeID = 1
	mnp._adjustSlotsAndCopyEvicts(localID, evicts)
	require.Len(t, mnp._adjustSlotsAndCopyEvicts(localID, evicts), len(evicts))

	require.Equal(t, us, mnp.local)

	mnp.local = us
	mnp._adjustSlotsAndCopyEvicts(localID, evicts)
	require.Len(t, mnp._adjustSlotsAndCopyEvicts(localID, evicts), len(evicts))

	us2 := &updatableSlot{}
	evicts = append(evicts, us2)
	mnp._adjustSlotsAndCopyEvicts(localID, evicts)
	require.Len(t, mnp._adjustSlotsAndCopyEvicts(localID, evicts), len(evicts))
}

func TestMakeOfProfiles(t *testing.T) {
	sp1 := profiles.NewStaticProfileMock(t)
	nodeID1 := insolar.AbsentShortNodeID
	sp1.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return *(&nodeID1) })
	sp2 := profiles.NewStaticProfileMock(t)
	nodeID2 := insolar.ShortNodeID(2)
	sp2.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return *(&nodeID2) })
	nodes := []profiles.StaticProfile{sp1, sp2}
	mnp := ManyNodePopulation{}
	svf := cryptkit.NewSignatureVerifierFactoryMock(t)
	sv := cryptkit.NewSignatureVerifierMock(t)
	svf.CreateSignatureVerifierWithPKSMock.Set(func(cryptkit.PublicKeyStore) cryptkit.SignatureVerifier { return sv })
	localNodeID := insolar.AbsentShortNodeID
	require.Panics(t, func() { mnp.makeOfProfiles(nodes, localNodeID, svf) })

	require.Panics(t, func() { mnp.makeOfProfiles(nil, localNodeID, svf) })

	nodeID1 = 1
	pks1 := cryptkit.NewPublicKeyStoreMock(t)
	sp1.GetPublicKeyStoreMock.Set(func() cryptkit.PublicKeyStore { return pks1 })
	sp1.GetPrimaryRoleMock.Set(func() member.PrimaryRole { return member.PrimaryRoleHeavyMaterial })
	pks2 := cryptkit.NewPublicKeyStoreMock(t)
	sp2.GetPublicKeyStoreMock.Set(func() cryptkit.PublicKeyStore { return pks2 })
	sp2.GetPrimaryRoleMock.Set(func() member.PrimaryRole { return member.PrimaryRoleHeavyMaterial })
	require.Panics(t, func() { mnp.makeOfProfiles(nodes, localNodeID, svf) })

	us := &updatableSlot{}
	mnp.local = us
	localNodeID = 1
	mnp.makeOfProfiles(nodes, localNodeID, svf)
	require.True(t, mnp.isInvalid)
}

func TestMNPFindProfile(t *testing.T) {
	mnp := ManyNodePopulation{}
	sp := profiles.NewStaticProfileMock(t)
	nodeID := insolar.ShortNodeID(1)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return nodeID })
	mnp.slotByID = make(map[insolar.ShortNodeID]*updatableSlot)
	us := &updatableSlot{}
	us.StaticProfile = sp
	mnp.slotByID[nodeID] = us
	nps := mnp.FindProfile(nodeID)
	require.Equal(t, sp, nps.GetStatic())

	nps = mnp.FindProfile(0)
	require.Nil(t, nps)
}

func TestMNPGetProfiles(t *testing.T) {
	mnp := ManyNodePopulation{}
	us := updatableSlot{}
	sp := profiles.NewStaticProfileMock(t)
	us.StaticProfile = sp
	mnp.slots = append(mnp.slots, us)
	an := mnp.GetProfiles()
	require.Len(t, an, 1)
	require.Equal(t, sp, an[0].GetStatic())
}

func TestMNPGetLocalProfile(t *testing.T) {
	mnp := ManyNodePopulation{}
	us := updatableSlot{}
	mnp.local = &us
	require.Equal(t, &us, mnp.GetLocalProfile())
}

func TestSetInvalid(t *testing.T) {
	mnp := ManyNodePopulation{}
	require.False(t, mnp.isInvalid)
	mnp.setInvalid()
	require.True(t, mnp.isInvalid)
}

func TestNewDynamicPopulationCopySelf(t *testing.T) {
	mnp := &ManyNodePopulation{}
	us := updatableSlot{}
	sp := profiles.NewStaticProfileMock(t)
	nodeID := insolar.ShortNodeID(1)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return nodeID })
	us.StaticProfile = sp
	mnp.local = &us
	dp := NewDynamicPopulationCopySelf(mnp)
	require.Equal(t, &us, dp.slotByID[nodeID])
}

func TestMakeCopyOf(t *testing.T) {
	slots := make([]updatableSlot, 1)
	us := &updatableSlot{}
	sp := profiles.NewStaticProfileMock(t)
	nodeID := insolar.ShortNodeID(1)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return nodeID })
	us.StaticProfile = sp
	dp := DynamicPopulation{}
	dp.makeCopyOf(slots, us)
	require.Equal(t, us, dp.slotByID[nodeID])
}

func TestDPFindProfile(t *testing.T) {
	dp := DynamicPopulation{}
	dp.slotByID = make(map[insolar.ShortNodeID]*updatableSlot, 1)
	nodeID := insolar.ShortNodeID(1)
	us := &updatableSlot{}
	dp.slotByID[nodeID] = us
	require.NotNil(t, dp.FindProfile(nodeID))

	require.Panics(t, func() { dp.FindProfile(2) })
}

func TestFindUpdatableProfile(t *testing.T) {
	dp := DynamicPopulation{}
	dp.slotByID = make(map[insolar.ShortNodeID]*updatableSlot, 1)
	nodeID := insolar.ShortNodeID(1)
	us := &updatableSlot{}
	dp.slotByID[nodeID] = us
	require.NotNil(t, dp.FindUpdatableProfile(nodeID))

	require.Nil(t, dp.FindUpdatableProfile(2))
}

func TestGetCount(t *testing.T) {
	dp := DynamicPopulation{}
	dp.slotByID = make(map[insolar.ShortNodeID]*updatableSlot, 2)
	us := &updatableSlot{}
	dp.slotByID[1] = us
	dp.slotByID[2] = us
	require.Equal(t, 2, dp.GetCount())
}

func testLessFunc(c profiles.ActiveNode, o profiles.ActiveNode) bool {
	return c.GetIndex() < o.GetIndex()
}

func TestDPGetProfiles(t *testing.T) {
	dp := DynamicPopulation{}
	dp.slotByID = make(map[insolar.ShortNodeID]*updatableSlot, 1)
	us := &updatableSlot{}
	dp.slotByID[1] = us
	dp.slotByID[2] = us
	require.Panics(t, func() { dp.GetProfiles() })
	us2 := &updatableSlot{}
	us2.index = 1
	dp.slotByID[2] = us2
	require.Len(t, dp.GetProfiles(), 2)
}

func TestDPGetUnorderedProfiles(t *testing.T) {
	dp := DynamicPopulation{}
	dp.slotByID = make(map[insolar.ShortNodeID]*updatableSlot, 2)
	us := &updatableSlot{}
	dp.slotByID[1] = us
	dp.slotByID[2] = us
	require.Len(t, dp.GetUnorderedProfiles(), 2)
}

func TestGetUnorderedSlots(t *testing.T) {
	dp := DynamicPopulation{}
	dp.slotByID = make(map[insolar.ShortNodeID]*updatableSlot, 2)
	us := &updatableSlot{}
	dp.slotByID[1] = us
	dp.slotByID[2] = us
	require.Len(t, dp.getUnorderedSlots(), 2)
}

func TestDPGetLocalProfile(t *testing.T) {
	dp := DynamicPopulation{}
	us := updatableSlot{}
	dp.local = &us
	require.Equal(t, &us, dp.GetLocalProfile())
}

func TestCopyAndSeparate(t *testing.T) {
	dp := DynamicPopulation{}
	dp.slotByID = make(map[insolar.ShortNodeID]*updatableSlot, 1)
	us := updatableSlot{}
	sp := profiles.NewStaticProfileMock(t)
	nodeID := insolar.ShortNodeID(1)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return nodeID })
	sp.GetPrimaryRoleMock.Set(func() member.PrimaryRole { return member.PrimaryRoleHeavyMaterial })
	spe := profiles.NewStaticProfileExtensionMock(t)
	sp.GetExtensionMock.Set(func() profiles.StaticProfileExtension { return spe })
	us.StaticProfile = sp
	dp.slotByID[nodeID] = &us
	dp.local = &us
	var forceInvalid bool
	r, evPop := dp.CopyAndSeparate(forceInvalid, nil)
	require.NotNil(t, r.local)

	require.Equal(t, evPop.GetCount(), 0)

	r, evPop = dp.CopyAndSeparate(forceInvalid, panicOnRecoverable)
	require.NotNil(t, r.local)

	require.Equal(t, evPop.GetCount(), 0)

	r, evPop = dp.CopyAndSeparate(true, panicOnRecoverable)
	require.NotNil(t, r.local)

	require.Equal(t, evPop.GetCount(), 0)

	delete(dp.slotByID, nodeID)
	require.Panics(t, func() { dp.CopyAndSeparate(true, panicOnRecoverable) })
}

func TestDPAddProfile(t *testing.T) {
	dp := DynamicPopulation{}
	dp.slotByID = make(map[insolar.ShortNodeID]*updatableSlot, 1)
	sp := profiles.NewStaticProfileMock(t)
	nodeID := insolar.ShortNodeID(1)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return nodeID })
	dp.AddProfile(sp)
	require.NotNil(t, dp.slotByID[nodeID])

	require.Panics(t, func() { dp.AddProfile(sp) })
}

func TestDPRemoveProfile(t *testing.T) {
	dp := DynamicPopulation{}
	dp.slotByID = make(map[insolar.ShortNodeID]*updatableSlot, 1)
	sp := profiles.NewStaticProfileMock(t)
	nodeID := insolar.ShortNodeID(1)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return nodeID })
	dp.AddProfile(sp)
	require.NotNil(t, dp.slotByID[nodeID])

	dp.RemoveProfile(nodeID)
	require.Nil(t, dp.slotByID[nodeID])
}

func TestRemoveOthers(t *testing.T) {
	dp := DynamicPopulation{}
	dp.slotByID = make(map[insolar.ShortNodeID]*updatableSlot, 2)
	sp1 := profiles.NewStaticProfileMock(t)
	nodeID1 := insolar.ShortNodeID(1)
	sp1.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return nodeID1 })
	dp.AddProfile(sp1)
	sp2 := profiles.NewStaticProfileMock(t)
	nodeID2 := insolar.ShortNodeID(2)
	sp2.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return nodeID2 })
	dp.AddProfile(sp2)
	dp.local = dp.slotByID[2]
	require.Len(t, dp.slotByID, 2)

	dp.RemoveOthers()
	require.Len(t, dp.slotByID, 1)

	require.Equal(t, dp.local, dp.slotByID[2])
}

func TestLen(t *testing.T) {
	ss := slotSorter{}
	size := 2
	ss.values = make([]*updatableSlot, size)
	require.Equal(t, size, ss.Len())
}

func TestLess(t *testing.T) {
	ss := slotSorter{lessFn: testLessFunc}
	ss.values = make([]*updatableSlot, 2)
	us1 := updatableSlot{}
	us1.index = 2
	ss.values[0] = &us1
	us2 := updatableSlot{}
	us2.index = 1
	ss.values[1] = &us2
	require.False(t, ss.Less(0, 1))
}

func TestSwap(t *testing.T) {
	ss := slotSorter{lessFn: testLessFunc}
	ss.values = make([]*updatableSlot, 2)
	us1 := updatableSlot{}
	ind1 := member.Index(2)
	ind2 := member.Index(1)
	us1.index = ind1
	ss.values[0] = &us1
	us2 := updatableSlot{}
	us2.index = ind2
	ss.values[1] = &us2
	require.Equal(t, ind1, ss.values[0].index)

	require.Equal(t, ind2, ss.values[1].index)

	ss.Swap(0, 1)
	require.Equal(t, ind2, ss.values[0].index)

	require.Equal(t, ind1, ss.values[1].index)
}

func TestPrepare(t *testing.T) {
	rr := roleRecord{}
	require.Panics(t, func() { rr.prepare() })

	rr.container = &ManyNodePopulation{}
	rr.powerPositions = make([]unitizedPowerPosition, 1)
	require.Panics(t, func() { rr.prepare() })

	rr.powerPositions = nil
	rr.rolePower = 0
	rr.prepare()
	require.Len(t, rr.powerPositions, 0)

	rr.rolePower = 1
	rr.container.slots = make([]updatableSlot, 3)
	us1 := updatableSlot{}
	us1.power = 1
	us2 := updatableSlot{}
	us2.power = 2
	us3 := updatableSlot{}
	us3.power = 3
	rr.container.slots[0] = us1
	rr.container.slots[1] = us2
	rr.container.slots[2] = us3
	rr.firstNode = 1
	rr.roleCount = 2
	require.Panics(t, func() { rr.prepare() })

	rr = roleRecord{}
	rr.container = &ManyNodePopulation{}
	rr.rolePower = 1
	rr.container.slots = make([]updatableSlot, 3)
	rr.container.slots[0] = us3
	rr.container.slots[1] = us2
	rr.container.slots[2] = us1
	rr.firstNode = 1
	rr.roleCount = 2
	require.Panics(t, func() { rr.prepare() })

	rr = roleRecord{}
	rr.container = &ManyNodePopulation{}
	rr.rolePower = 79
	size := uint16(24)
	rr.container.slots = make([]updatableSlot, size)
	us1 = updatableSlot{}
	us1.power = 11
	rr.container.slots[0] = us1
	us2 = updatableSlot{}
	us2.power = 10
	rr.container.slots[1] = us2
	us3 = updatableSlot{}
	us3.power = 9
	rr.container.slots[2] = us3
	us4 := updatableSlot{}
	us4.power = 8
	rr.container.slots[3] = us4
	us5 := updatableSlot{}
	us5.power = 7
	rr.container.slots[4] = us5
	us6 := updatableSlot{}
	us6.power = 6
	rr.container.slots[5] = us6
	us7 := updatableSlot{}
	us7.power = 5
	rr.container.slots[6] = us7
	us8 := updatableSlot{}
	us8.power = 4
	rr.container.slots[7] = us8
	us9 := updatableSlot{}
	us9.power = 3
	rr.container.slots[8] = us9
	us10 := updatableSlot{}
	us10.power = 2
	rr.container.slots[9] = us10
	us11 := updatableSlot{}
	us11.power = 1
	for i := 10; i < 24; i++ {
		rr.container.slots[i] = us11
	}
	rr.firstNode = 0
	rr.roleCount = size
	rr.prepare()
	require.Len(t, rr.powerPositions, 11)
}

func TestIsValid(t *testing.T) {
	rr := roleRecord{}
	require.False(t, rr.IsValid())

	rr.container = &ManyNodePopulation{}
	require.True(t, rr.IsValid())

	rr.container.isInvalid = true
	require.False(t, rr.IsValid())
}

func TestGetPrimaryRole(t *testing.T) {
	rr := roleRecord{}
	role := member.PrimaryRoleNeutral
	rr.role = role
	require.Equal(t, role, rr.GetPrimaryRole())
}

func TestGetWorkingPower(t *testing.T) {
	rr := roleRecord{}
	power := uint32(1)
	rr.rolePower = power
	require.Equal(t, power, rr.GetWorkingPower())
}

func TestGetWorkingCount(t *testing.T) {
	rr := roleRecord{}
	roleCount := uint16(1)
	rr.roleCount = roleCount
	require.Equal(t, int(roleCount), rr.GetWorkingCount())
}

func TestGetIdleCount(t *testing.T) {
	rr := roleRecord{}
	idleCount := uint16(1)
	rr.idleCount = idleCount
	require.Equal(t, int(idleCount), rr.GetIdleCount())
}

func TestRRGetProfiles(t *testing.T) {
	rr := roleRecord{}
	require.Panics(t, func() { rr.GetProfiles() })

	rr.container = &ManyNodePopulation{}
	require.Nil(t, rr.GetProfiles())

	rr.container.isInvalid = true
	require.Panics(t, func() { rr.GetProfiles() })

	rr.container.isInvalid = false
	rr.roleCount = 2
	rr.container.slots = make([]updatableSlot, rr.roleCount)
	rr.container.slots[0] = updatableSlot{}
	rr.container.slots[1] = updatableSlot{}
	require.Len(t, rr.GetProfiles(), 2)
}

func TestGetAssignmentByPower(t *testing.T) {
	rr := roleRecord{}
	metric := uint64(1)
	excludeID := insolar.ShortNodeID(1)
	assigned, excluded := rr.GetAssignmentByPower(metric, excludeID)
	require.Nil(t, assigned)

	require.Nil(t, excluded)

	rr.roleCount = 1
	assigned, excluded = rr.GetAssignmentByPower(metric, excludeID)
	require.Nil(t, assigned)

	require.Nil(t, excluded)

	rr.rolePower = 1
	assigned, excluded = rr.GetAssignmentByPower(metric, excludeID)
	require.Nil(t, assigned)

	require.Nil(t, excluded)

	rr.container = &ManyNodePopulation{}
	rr.container.slots = make([]updatableSlot, rr.roleCount)
	us := updatableSlot{}
	sp := profiles.NewStaticProfileMock(t)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	us.StaticProfile = sp
	rr.container.slots[0] = us

	assigned, excluded = rr.GetAssignmentByPower(metric, excludeID)
	require.NotNil(t, assigned)

	require.Nil(t, excluded)

	excludeID = 0
	assigned, excluded = rr.GetAssignmentByPower(metric, excludeID)
	require.NotNil(t, assigned)

	require.NotNil(t, excluded)

	rr.container.slots = append(rr.container.slots, updatableSlot{})
	rr.roleCount = 2
	rr.powerPositions = make([]unitizedPowerPosition, 3)
	rr.powerPositions[0] = unitizedPowerPosition{powerUnit: 1}
	rr.powerPositions[1] = unitizedPowerPosition{powerUnit: 2}
	require.Panics(t, func() { rr.GetAssignmentByPower(metric, excludeID) })

	excludeID = 1
	assigned, excluded = rr.GetAssignmentByPower(metric, excludeID)
	require.NotNil(t, assigned)

	require.Nil(t, excluded)

	metric = 9
	rr.rolePower = 3
	excludeID = 0
	rr.powerPositions[2] = unitizedPowerPosition{powerUnit: 2, powerStartsAt: 1, indexStartsAt: 1, unitCount: 1}
	sp2 := profiles.NewStaticProfileMock(t)
	sp2.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 1 })
	rr.container.slots[1].StaticProfile = sp2
	assigned, excluded = rr.GetAssignmentByPower(metric, excludeID)
	require.NotNil(t, assigned)

	require.NotNil(t, excluded)
}

func TestGetAssignmentByCount(t *testing.T) {
	rr := roleRecord{}
	metric := uint64(1)
	excludeID := insolar.ShortNodeID(1)
	assigned, excluded := rr.GetAssignmentByCount(metric, excludeID)
	require.Nil(t, assigned)

	require.Nil(t, excluded)

	rr.roleCount = 1
	assigned, excluded = rr.GetAssignmentByCount(metric, excludeID)
	require.Nil(t, assigned)

	require.Nil(t, excluded)

	rr.container = &ManyNodePopulation{}
	rr.container.slots = make([]updatableSlot, rr.roleCount)
	us := updatableSlot{}
	sp := profiles.NewStaticProfileMock(t)
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 0 })
	us.StaticProfile = sp
	rr.container.slots[0] = us

	assigned, excluded = rr.GetAssignmentByCount(metric, excludeID)
	require.NotNil(t, assigned)

	require.Nil(t, excluded)

	excludeID = 0
	assigned, excluded = rr.GetAssignmentByCount(metric, excludeID)
	require.NotNil(t, assigned)

	require.NotNil(t, excluded)

	rr.container.slots = append(rr.container.slots, updatableSlot{})
	sp2 := profiles.NewStaticProfileMock(t)
	sp2.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 1 })
	rr.container.slots[1].StaticProfile = sp2
	rr.roleCount = 2
	rr.powerPositions = make([]unitizedPowerPosition, 3)
	rr.powerPositions[0] = unitizedPowerPosition{powerUnit: 1}
	rr.powerPositions[1] = unitizedPowerPosition{powerUnit: 2}
	assigned, excluded = rr.GetAssignmentByCount(metric, excludeID)
	require.NotNil(t, assigned)

	require.Nil(t, excluded)

	excludeID = 1
	assigned, excluded = rr.GetAssignmentByCount(metric, excludeID)
	require.NotNil(t, assigned)

	require.NotNil(t, excluded)

	metric = 9
	rr.powerPositions[2] = unitizedPowerPosition{}
	sp.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return 1 })
	require.Panics(t, func() { rr.GetAssignmentByCount(metric, excludeID) })
}

func TestGetByIndex(t *testing.T) {
	rr := roleRecord{}
	rr.container = &ManyNodePopulation{}
	rr.container.slots = make([]updatableSlot, 2)
	require.NotNil(t, rr.getByIndex(1))
}

func TestGetIndexByPower(t *testing.T) {
	rr := roleRecord{}
	rr.powerPositions = make([]unitizedPowerPosition, 3)
	rr.powerPositions[0] = unitizedPowerPosition{powerUnit: 1, powerStartsAt: 1}
	rr.powerPositions[1] = unitizedPowerPosition{powerUnit: 2, powerStartsAt: 2}
	require.Equal(t, uint16(0), rr.getIndexByPower(1))

	rr.powerPositions[0].unitCount = 1
	rr.powerPositions[0].indexStartsAt = 2
	require.Equal(t, uint16(2), rr.getIndexByPower(1))
}
