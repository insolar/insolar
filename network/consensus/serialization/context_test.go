package serialization

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPacketContext_InContext(t *testing.T) {
	ctx := newPacketContext(context.Background(), &Header{})

	require.True(t, ctx.InContext(NoContext))
	require.False(t, ctx.InContext(ContextMembershipAnnouncement))
	require.False(t, ctx.InContext(ContextNeighbourAnnouncement))

	ctx.fieldContext = ContextMembershipAnnouncement
	require.True(t, ctx.InContext(ContextMembershipAnnouncement))

	ctx.fieldContext = ContextNeighbourAnnouncement
	require.True(t, ctx.InContext(ContextNeighbourAnnouncement))
}

func TestPacketContext_SetInContext(t *testing.T) {
	ctx := newPacketContext(context.Background(), &Header{})

	require.True(t, ctx.InContext(NoContext))

	ctx.SetInContext(ContextMembershipAnnouncement)
	require.True(t, ctx.InContext(ContextMembershipAnnouncement))

	ctx.SetInContext(ContextNeighbourAnnouncement)
	require.True(t, ctx.InContext(ContextNeighbourAnnouncement))
}

func TestPacketContext_GetNeighbourNodeID(t *testing.T) {
	ctx := newPacketContext(context.Background(), &Header{})

	ctx.neighbourNodeID = 123
	require.EqualValues(t, 123, ctx.GetNeighbourNodeID())
}

func TestPacketContext_GetNeighbourNodeID_Panics(t *testing.T) {
	ctx := newPacketContext(context.Background(), &Header{})

	require.Panics(t, func() { ctx.GetNeighbourNodeID() })
}

func TestPacketContext_SetNeighbourNodeID(t *testing.T) {
	ctx := newPacketContext(context.Background(), &Header{})

	ctx.SetNeighbourNodeID(123)
	require.EqualValues(t, 123, ctx.GetNeighbourNodeID())
}

func TestPacketContext_GetAnnouncedJoinerNodeID(t *testing.T) {
	ctx := newPacketContext(context.Background(), &Header{})

	require.EqualValues(t, 0, ctx.GetAnnouncedJoinerNodeID())

	ctx.announcedJoinerNodeID = 123
	require.EqualValues(t, 123, ctx.GetAnnouncedJoinerNodeID())
}

func TestPacketContext_SetAnnouncedJoinerNodeID(t *testing.T) {
	ctx := newPacketContext(context.Background(), &Header{})

	require.EqualValues(t, 0, ctx.GetAnnouncedJoinerNodeID())

	ctx.SetAnnouncedJoinerNodeID(123)
	require.EqualValues(t, 123, ctx.GetAnnouncedJoinerNodeID())
}
