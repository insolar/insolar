package phases

// func TestSetReceivedPhase(t *testing.T) {
// 	lp := profiles.NewLocalNodeMock(t)
// 	lp.LocalNodeProfileMock.Set(func() {})
// 	callback := &nodeContext{}
// 	r := NewNodeAppearanceAsSelf(lp, callback)
// 	require.True(t, r.SetReceivedPhase(member.Phase1))
//
// 	require.False(t, r.SetReceivedPhase(member.Phase1))
// }
//
// func TestSetReceivedByPacketType(t *testing.T) {
// 	lp := profiles.NewLocalNodeMock(t)
// 	lp.LocalNodeProfileMock.Set(func() {})
// 	callback := &nodeContext{}
// 	r := NewNodeAppearanceAsSelf(lp, callback)
// 	require.True(t, r.SetReceivedByPacketType(member.PacketPhase1))
//
// 	require.False(t, r.SetReceivedByPacketType(member.PacketPhase1))
//
// 	require.False(t, r.SetReceivedByPacketType(member.MaxPacketType))
// }
//
// func TestSetSentPhase(t *testing.T) {
// 	lp := profiles.NewLocalNodeMock(t)
// 	lp.LocalNodeProfileMock.Set(func() {})
// 	callback := &nodeContext{}
// 	r := NewNodeAppearanceAsSelf(lp, callback)
// 	require.True(t, r.SetSentPhase(member.Phase1))
//
// 	require.False(t, r.SetSentPhase(member.Phase1))
// }
//
// func TestSetSentByPacketType(t *testing.T) {
// 	lp := profiles.NewLocalNodeMock(t)
// 	lp.LocalNodeProfileMock.Set(func() {})
// 	callback := &nodeContext{}
// 	r := NewNodeAppearanceAsSelf(lp, callback)
// 	require.True(t, r.SetSentByPacketType(member.PacketPhase1))
//
// 	require.True(t, r.SetSentByPacketType(member.PacketPhase1))
//
// 	require.False(t, r.SetSentByPacketType(member.MaxPacketType))
// }
//
// func TestSetReceivedWithDupCheck(t *testing.T) {
// 	lp := profiles.NewLocalNodeMock(t)
// 	lp.LocalNodeProfileMock.Set(func() {})
// 	callback := &nodeContext{}
// 	r := NewNodeAppearanceAsSelf(lp, callback)
// 	require.Nil(t, r.SetReceivedWithDupCheck(member.PacketPhase1))
//
// 	require.Equal(t, errors.ErrRepeatedPhasePacket, r.SetReceivedWithDupCheck(member.PacketPhase1))
//
// 	require.Equal(t, errors.ErrRepeatedPhasePacket, r.SetReceivedWithDupCheck(member.MaxPacketType))
// }
