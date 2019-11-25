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
	"math"
)

type PayloadFactory interface {
	CreatePayloadBuilder(format FileFormat, sr StorageSeqReader) PayloadBuilder
}

type PayloadBuilder interface {
	AddPreamble(bytes []byte, storageOffset int64, localOffset int) error
	AddConclude(bytes []byte, storageOffset int64, localOffset int) error
	AddEntry(entryNo, entryFieldNo int, bytes []byte, storageStart int64, localOffset int) error
	NeedsNextEntry() bool

	Finished() error
	Failed(error) error
}

type FileFormat uint16
type FormatOptions uint64

func OpenStorage(sr StorageSeqReader, config ReadConfig, payloadFactory PayloadFactory) (PayloadBuilder, error) {
	f, opt, err := ReadFormatAndOptions(sr)
	if err != nil {
		return nil, err
	}

	config.StorageOptions = opt
	var pb PayloadBuilder

	switch f {
	case 1:
		pb = payloadFactory.CreatePayloadBuilder(f, sr)
		v1 := StorageFileV1{Config: config, Builder: pb}
		err = v1.Read(sr)
	default:
		return nil, fmt.Errorf("unknown storage format: format%x", f)
	}

	if err == nil {
		err = pb.Finished()
	}
	if err != nil {
		err = pb.Failed(err)
	}
	return pb, err
}

func ReadFormatAndOptions(sr StorageSeqReader) (FileFormat, FormatOptions, error) {
	if v, err := formatField.DecodeFrom(sr); err != nil {
		return 0, 0, err
	} else {
		return FileFormat(v & math.MaxUint16), FormatOptions(v >> 16), nil
	}
}

type ReadConfig struct {
	ReadAllEntries bool
	AlwaysCopy     bool
	StorageOptions FormatOptions
}
