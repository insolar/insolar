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

package dbsv1

import (
	"github.com/insolar/insolar/ledger-v2/dropbag/dbcommon"
)

const FormatId dbcommon.FileFormat = 1

func OpenReadStorage(sr dbcommon.StorageSeqReader, payloadFactory dbcommon.PayloadFactory,
	config dbcommon.ReadConfig, options dbcommon.FormatOptions,
) (dbcommon.PayloadBuilder, error) {
	pb := payloadFactory.CreatePayloadBuilder(FormatId, sr)
	v1 := StorageFileV1Reader{Config: config, Builder: pb}
	v1.StorageOptions = options
	return pb, v1.Read(sr)
}

func PrepareWriteStorage(sw dbcommon.StorageSeqWriter, options dbcommon.FormatOptions) (dbcommon.PayloadWriter, error) {
	return NewStorageFileV1Writer(sw, StorageFileV1{StorageOptions: options})
}
