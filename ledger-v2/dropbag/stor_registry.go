//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package dropbag

import (
	"fmt"
	"github.com/insolar/insolar/ledger-v2/dropbag/dbcommon"
	"github.com/insolar/insolar/ledger-v2/dropbag/dbsv1"
)

func OpenReadStorage(sr dbcommon.StorageSeqReader, config dbcommon.ReadConfig, payloadFactory dbcommon.PayloadFactory) (dbcommon.PayloadBuilder, error) {
	f, opt, err := dbcommon.ReadFormatAndOptions(sr)
	if err != nil {
		return nil, err
	}

	var pb dbcommon.PayloadBuilder
	switch f {
	case 1:
		pb, err = dbsv1.OpenReadStorage(sr, payloadFactory, config, opt)
	default:
		return nil, fmt.Errorf("unknown storage format: format%x", f)
	}

	if pb == nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to read storage: format=%x", f)
	}

	if err == nil {
		err = pb.Finished()
	}
	if err != nil {
		err = pb.Failed(err)
	}
	return pb, err
}

func OpenWriteStorage(sw dbcommon.StorageSeqWriter, f dbcommon.FileFormat, options dbcommon.FormatOptions) (dbcommon.PayloadWriter, error) {
	var pw dbcommon.PayloadWriter
	var err error

	switch f {
	case 1:
		pw, err = dbsv1.PrepareWriteStorage(sw, options)
	default:
		return nil, fmt.Errorf("unknown storage format: format%x", f)
	}
	if pw == nil {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to write storage: format=%x", f)
	}

	if err != nil {
		_ = pw.Close()
		return nil, err
	}

	if err := dbcommon.WriteFormatAndOptions(sw, f, options); err != nil {
		_ = pw.Close()
		return nil, err
	}
	return pw, nil
}
