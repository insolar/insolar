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
	CreatePayloadBuilder(format FileFormat, sr StorageReader) PayloadBuilder
}

type PayloadBuilder interface {
	AddPreamble(bytes []byte, storageStart int64, storageLen int64) error
	AddConclude(bytes []byte, storageStart int64, storageLen int64) error
	AddEntry(entryNo int, bytes []byte, storageStart int64, storageLen int64) error
	Finished() error
	Failed(error) error
}

type FileFormat uint16
type FormatOptions uint64

func OpenStorage(sr StorageReader, config ReadConfig, payloadFactory PayloadFactory) (PayloadBuilder, error) {
	f, opt, err := ReadFormatAndOptions(sr)
	if err != nil {
		return nil, err
	}

	config.StorageOptions = opt
	var pb PayloadBuilder

	switch f {
	case 1:
		pb = payloadFactory.CreatePayloadBuilder(f, sr)
		err = StorageFileV1{Config: config, Builder: pb}.Read(sr)
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

func ReadFormatAndOptions(sr StorageReader) (FileFormat, FormatOptions, error) {
	if v, err := formatField.DecodeFrom(sr); err != nil {
		return 0, 0, err
	} else {
		return FileFormat(v & math.MaxUint16), FormatOptions(v >> 16), nil
	}
}

type ReadConfig struct {
	NonLazyEntries bool
	AlwaysCopy     bool
	StorageOptions FormatOptions
}
