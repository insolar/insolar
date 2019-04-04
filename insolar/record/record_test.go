package record

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
)

func FuzzRandomID(t *insolar.ID, _ fuzz.Continue) {
	*t = gen.ID()
}

func FuzzRandomReference(t *insolar.Reference, _ fuzz.Continue) {
	*t = gen.Reference()
}

func TestMarshalUnmarshalRecord(t *testing.T) {
	f := fuzz.New().Funcs(FuzzRandomID, FuzzRandomReference).NumElements(50, 100).NilChance(0)

	t.Run("GenesisRecordTest", func(t *testing.T) {
		a := assert.New(t)
		t.Parallel()
		var record GenesisRecord

		for i := 0; i < 1; i++ {
			f.Fuzz(&record)

			bin, err := MarshalRecord(&record)
			a.NoError(err)
			for i := 0; i < 2; i++ {
				binNew, err := MarshalRecord(&record)
				a.NoError(err)
				a.Equal(bin, binNew)

				recordNew, err := UnmarshalRecord(binNew)

				a.Equal(&record, recordNew)
			}
		}
	})

	t.Run("ChildRecordTest", func(t *testing.T) {
		a := assert.New(t)
		t.Parallel()
		var record ChildRecord

		for i := 0; i < 10; i++ {
			f.Fuzz(&record)

			bin, err := MarshalRecord(&record)
			a.NoError(err)
			for i := 0; i < 100; i++ {
				binNew, err := MarshalRecord(&record)
				a.NoError(err)
				a.Equal(bin, binNew)

				recordNew, err := UnmarshalRecord(binNew)

				a.Equal(&record, recordNew)
			}
		}
	})

	t.Run("JetRecordTest", func(t *testing.T) {
		a := assert.New(t)
		t.Parallel()
		var record JetRecord

		for i := 0; i < 10; i++ {
			f.Fuzz(&record)

			bin, err := MarshalRecord(&record)
			a.NoError(err)
			for i := 0; i < 100; i++ {
				binNew, err := MarshalRecord(&record)
				a.NoError(err)
				a.Equal(bin, binNew)

				recordNew, err := UnmarshalRecord(binNew)

				a.Equal(&record, recordNew)
			}
		}
	})

	t.Run("RequestRecordTest", func(t *testing.T) {
		a := assert.New(t)
		t.Parallel()
		var record RequestRecord

		for i := 0; i < 10; i++ {
			f.Fuzz(&record)

			bin, err := MarshalRecord(&record)
			a.NoError(err)
			for i := 0; i < 100; i++ {
				binNew, err := MarshalRecord(&record)
				a.NoError(err)
				a.Equal(bin, binNew)

				recordNew, err := UnmarshalRecord(binNew)

				a.Equal(&record, recordNew)
			}
		}
	})

	t.Run("ResultRecordTest", func(t *testing.T) {
		a := assert.New(t)
		t.Parallel()
		var record ResultRecord

		for i := 0; i < 10; i++ {
			f.Fuzz(&record)

			bin, err := MarshalRecord(&record)
			a.NoError(err)
			for i := 0; i < 100; i++ {
				binNew, err := MarshalRecord(&record)
				a.NoError(err)
				a.Equal(bin, binNew)

				recordNew, err := UnmarshalRecord(binNew)

				a.Equal(&record, recordNew)
			}
		}
	})

	t.Run("TypeRecordTest", func(t *testing.T) {
		a := assert.New(t)
		t.Parallel()
		var record TypeRecord

		for i := 0; i < 10; i++ {
			f.Fuzz(&record)

			bin, err := MarshalRecord(&record)
			a.NoError(err)
			for i := 0; i < 100; i++ {
				binNew, err := MarshalRecord(&record)
				a.NoError(err)
				a.Equal(bin, binNew)

				recordNew, err := UnmarshalRecord(binNew)

				a.Equal(&record, recordNew)
			}
		}
	})

	t.Run("CodeRecordTest", func(t *testing.T) {
		a := assert.New(t)
		t.Parallel()
		var record CodeRecord

		for i := 0; i < 10; i++ {
			f.Fuzz(&record)

			bin, err := MarshalRecord(&record)
			a.NoError(err)
			for i := 0; i < 100; i++ {
				binNew, err := MarshalRecord(&record)
				a.NoError(err)
				a.Equal(bin, binNew)

				recordNew, err := UnmarshalRecord(binNew)

				a.Equal(&record, recordNew)
			}
		}
	})

	t.Run("ObjectActivateRecordTest", func(t *testing.T) {
		a := assert.New(t)
		t.Parallel()
		var record ObjectActivateRecord

		for i := 0; i < 10; i++ {
			f.Fuzz(&record)

			bin, err := MarshalRecord(&record)
			a.NoError(err)
			for i := 0; i < 100; i++ {
				binNew, err := MarshalRecord(&record)
				a.NoError(err)
				a.Equal(bin, binNew)

				recordNew, err := UnmarshalRecord(binNew)

				a.Equal(&record, recordNew)
			}
		}
	})

	t.Run("ObjectAmendRecordTest", func(t *testing.T) {
		a := assert.New(t)
		t.Parallel()
		var record ObjectAmendRecord

		for i := 0; i < 10; i++ {
			f.Fuzz(&record)

			bin, err := MarshalRecord(&record)
			a.NoError(err)
			for i := 0; i < 100; i++ {
				binNew, err := MarshalRecord(&record)
				a.NoError(err)
				a.Equal(bin, binNew)

				recordNew, err := UnmarshalRecord(binNew)

				a.Equal(&record, recordNew)
			}
		}
	})

	t.Run("ObjectDeactivateRecordTest", func(t *testing.T) {
		a := assert.New(t)
		t.Parallel()
		var record ObjectDeactivateRecord

		for i := 0; i < 10; i++ {
			f.Fuzz(&record)

			bin, err := MarshalRecord(&record)
			a.NoError(err)
			for i := 0; i < 100; i++ {
				binNew, err := MarshalRecord(&record)
				a.NoError(err)
				a.Equal(bin, binNew)

				recordNew, err := UnmarshalRecord(binNew)

				a.Equal(&record, recordNew)
			}
		}
	})
}
