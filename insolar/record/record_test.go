// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package record

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/platformpolicy"
)

func FuzzRandomID(t *insolar.ID, _ fuzz.Continue) {
	*t = gen.ID()
}

func FuzzRandomReference(t *insolar.Reference, _ fuzz.Continue) {
	*t = gen.Reference()
}

func fuzzer() *fuzz.Fuzzer {
	return fuzz.New().Funcs(FuzzRandomID, FuzzRandomReference).NumElements(50, 100).NilChance(0)
}

func TestMarshalUnmarshalRecord(t *testing.T) {

	t.Run("GenesisRecordTest", func(t *testing.T) {
		f := fuzzer()
		a := assert.New(t)
		t.Parallel()
		var record Genesis

		for i := 0; i < 1; i++ {
			f.Fuzz(&record)

			bin, err := record.Marshal()
			a.NoError(err)
			for i := 0; i < 2; i++ {
				binNew, err := record.Marshal()
				a.NoError(err)
				a.Equal(bin, binNew)

				var recordNew Genesis
				err = recordNew.Unmarshal(binNew)
				require.NoError(t, err)

				a.Equal(&record, &recordNew)
			}
		}
	})

	t.Run("RequestRecordTest", func(t *testing.T) {
		f := fuzzer()
		a := assert.New(t)
		t.Parallel()
		var record IncomingRequest

		for i := 0; i < 10; i++ {
			f.Fuzz(&record)

			bin, err := record.Marshal()
			a.NoError(err)
			for i := 0; i < 2; i++ {
				binNew, err := record.Marshal()
				a.NoError(err)
				a.Equal(bin, binNew)

				var recordNew IncomingRequest
				err = recordNew.Unmarshal(binNew)
				require.NoError(t, err)

				a.Equal(&record, &recordNew)
			}
		}
	})

	t.Run("ResultRecordTest", func(t *testing.T) {
		f := fuzzer()
		a := assert.New(t)
		t.Parallel()
		var record Result

		for i := 0; i < 10; i++ {
			f.Fuzz(&record)

			bin, err := record.Marshal()
			a.NoError(err)
			for i := 0; i < 2; i++ {
				binNew, err := record.Marshal()
				a.NoError(err)
				a.Equal(bin, binNew)

				var recordNew Result
				err = recordNew.Unmarshal(binNew)
				require.NoError(t, err)

				a.Equal(&record, &recordNew)
			}
		}
	})

	t.Run("CodeRecordTest", func(t *testing.T) {
		f := fuzzer()
		a := assert.New(t)
		t.Parallel()
		var record Code

		for i := 0; i < 10; i++ {
			f.Fuzz(&record)

			bin, err := record.Marshal()
			a.NoError(err)
			for i := 0; i < 2; i++ {
				binNew, err := record.Marshal()
				a.NoError(err)
				a.Equal(bin, binNew)

				var recordNew Code
				err = recordNew.Unmarshal(binNew)
				require.NoError(t, err)

				a.Equal(&record, &recordNew)
			}
		}
	})

	t.Run("ActivateRecordTest", func(t *testing.T) {
		f := fuzzer()
		a := assert.New(t)
		t.Parallel()
		var record Activate

		for i := 0; i < 10; i++ {
			f.Fuzz(&record)

			bin, err := record.Marshal()
			a.NoError(err)
			for i := 0; i < 2; i++ {
				binNew, err := record.Marshal()
				a.NoError(err)
				a.Equal(bin, binNew)

				var recordNew Activate
				err = recordNew.Unmarshal(binNew)
				require.NoError(t, err)

				a.Equal(&record, &recordNew)
			}
		}
	})

	t.Run("AmendRecordTest", func(t *testing.T) {
		f := fuzzer()
		a := assert.New(t)
		t.Parallel()
		var record Amend

		for i := 0; i < 10; i++ {
			f.Fuzz(&record)

			bin, err := record.Marshal()
			a.NoError(err)
			for i := 0; i < 2; i++ {
				binNew, err := record.Marshal()
				a.NoError(err)
				a.Equal(bin, binNew)

				var recordNew Amend
				err = recordNew.Unmarshal(binNew)
				require.NoError(t, err)

				a.Equal(&record, &recordNew)
			}
		}
	})

	t.Run("DeactivateRecordTest", func(t *testing.T) {
		f := fuzzer()
		a := assert.New(t)
		t.Parallel()
		var record Deactivate

		for i := 0; i < 10; i++ {
			f.Fuzz(&record)

			bin, err := record.Marshal()
			a.NoError(err)
			for i := 0; i < 2; i++ {
				binNew, err := record.Marshal()
				a.NoError(err)
				a.Equal(bin, binNew)

				var recordNew Deactivate
				err = recordNew.Unmarshal(binNew)
				require.NoError(t, err)

				a.Equal(&record, &recordNew)
			}
		}
	})
}

func TestRequestInterface_IncomingRequest(t *testing.T) {
	t.Parallel()
	objref := gen.Reference()
	req := &IncomingRequest{
		Caller: gen.Reference(),
		Object: &objref,
		Reason: gen.Reference(),
	}
	iface := Request(req)
	require.Equal(t, false, iface.IsTemporaryUploadCode())
	require.Equal(t, false, iface.IsCreationRequest())
	require.Equal(t, req.Object, iface.AffinityRef())
	require.Equal(t, req.Reason, iface.ReasonRef())
	require.Equal(t, req.Caller, iface.ReasonAffinityRef())
}

func TestRequestInterface_OutgoingRequest(t *testing.T) {
	t.Parallel()
	objref := gen.Reference()
	req := &OutgoingRequest{
		Caller: gen.Reference(),
		Object: &objref,
		Reason: gen.Reference(),
	}
	iface := Request(req)
	require.Equal(t, false, iface.IsTemporaryUploadCode())
	require.Equal(t, false, iface.IsCreationRequest())
	require.Equal(t, &req.Caller, iface.AffinityRef())
	require.Equal(t, req.Reason, iface.ReasonRef())
	require.Equal(t, req.Caller, iface.ReasonAffinityRef())
}

func TestRequestInterface_IncomingRequestSaveAsChild(t *testing.T) {
	t.Parallel()

	objref := gen.Reference()
	req := &IncomingRequest{
		CallType: CTSaveAsChild,
		Object:   &objref,
		Reason:   gen.Reference(),
		APINode:  gen.Reference(),
	}

	iface := Request(req)
	require.Equal(t, false, iface.IsTemporaryUploadCode())
	require.Equal(t, true, iface.IsCreationRequest())
	require.Equal(t, (*insolar.Reference)(nil), iface.AffinityRef())
	require.Equal(t, req.Reason, iface.ReasonRef())
	require.True(t, iface.ReasonAffinityRef().IsEmpty())

	pn := insolar.PulseNumber(256)
	pcs := platformpolicy.NewPlatformCryptographyScheme()

	realAffinityRef := CalculateRequestAffinityRef(iface, pn, pcs)
	assert.NotNil(t, realAffinityRef)
}

func TestRequestInterface_OutgoingRequestSaveAsChild(t *testing.T) {
	t.Parallel()
	objref := gen.Reference()
	req := &OutgoingRequest{
		Caller:   gen.Reference(),
		Object:   &objref,
		Reason:   gen.Reference(),
		CallType: CTSaveAsChild,
	}
	iface := Request(req)
	require.False(t, iface.IsCreationRequest())
}
