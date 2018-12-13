package core_test

import (
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils"
	"github.com/jbenet/go-base58"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//RecordID and RecordRef serialization tests

func TestNewIDFromBase58(t *testing.T) {
	id := testutils.RandomID()
	idStr := base58.Encode(id[:])
	id2, err := core.NewIDFromBase58(idStr)
	require.NoError(t, err)

	assert.Equal(t, id, *id2)
}

func TestRecordID_String(t *testing.T) {
	id := testutils.RandomID()
	idStr := base58.Encode(id[:])

	assert.Equal(t, idStr, id.String())
}

func TestNewRefFromBase58(t *testing.T) {
	recordID := testutils.RandomID()
	domainID := testutils.RandomID()
	refStr := recordID.String() + core.RecordRefIDSeparator + domainID.String()

	expectedRef := core.NewRecordRef(domainID, recordID)
	actualRef, err := core.NewRefFromBase58(refStr)
	require.NoError(t, err)

	assert.Equal(t, expectedRef, actualRef)
}

func TestRecordRef_String(t *testing.T) {
	ref := testutils.RandomRef()
	expectedRefStr := ref.Record().String() + core.RecordRefIDSeparator + ref.Domain().String()

	assert.Equal(t, expectedRefStr, ref.String())
}
