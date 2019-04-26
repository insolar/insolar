//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

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

func fuzzer() *fuzz.Fuzzer {
	return fuzz.New().Funcs(FuzzRandomID, FuzzRandomReference).NumElements(50, 100).NilChance(0)
}

func TestMarshalUnmarshalRecord(t *testing.T) {

	t.Run("GenesisRecordTest", func(t *testing.T) {
		f := fuzzer()
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
		f := fuzzer()
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
		f := fuzzer()
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
		f := fuzzer()
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
		f := fuzzer()
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
		f := fuzzer()
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
		f := fuzzer()
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
		f := fuzzer()
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
		f := fuzzer()
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
		f := fuzzer()
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
