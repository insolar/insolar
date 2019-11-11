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

package common

import (
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
)

type ParsedRequestInfo struct {
	requestInfo *payload.RequestInfo

	RequestReference       insolar.Reference
	RequestObjectReference insolar.Reference
	RequestDeduplicated    bool
	Request                record.Request
	Result                 *record.Result
}

func (i *ParsedRequestInfo) fromRequestInfo(info *payload.RequestInfo) error {
	i.requestInfo = info
	i.RequestReference = *insolar.NewReference(info.RequestID)

	if info.Request != nil {
		rec := record.Material{}
		if err := rec.Unmarshal(info.Request); err != nil {
			return errors.Wrap(err, "failed to unmarshal request record")
		}

		virtual := record.Unwrap(&rec.Virtual)
		switch request := virtual.(type) {
		case *record.IncomingRequest:
			i.Request = request
		case *record.OutgoingRequest:
			i.Request = request
		default:
			return errors.Errorf("unexpected type '%T' when unpacking request", virtual)
		}

		i.RequestDeduplicated = true
	}

	if info.Result != nil {
		rec := record.Material{}
		if err := rec.Unmarshal(info.Request); err != nil {
			return errors.Wrap(err, "failed to unmarshal request record")
		}

		virtual := record.Unwrap(&rec.Virtual)
		result, ok := virtual.(*record.Result)
		if !ok {
			return errors.Errorf("unexpected type '%T' when unpacking incoming", virtual)
		}

		i.Result = result
	}

	if i.Request.AffinityRef() != nil {
		i.RequestObjectReference = *i.Request.AffinityRef()
	} else {
		i.RequestObjectReference = *insolar.NewReference(info.ObjectID)
	}

	return nil
}

func (i *ParsedRequestInfo) GetResultBytes() []byte {
	if i.Result != nil {
		return i.Result.Payload
	}
	return nil
}

func NewParsedRequestInfo(request record.Request, rawInfo *payload.RequestInfo) (*ParsedRequestInfo, error) {
	info := &ParsedRequestInfo{Request: request}
	err := info.fromRequestInfo(rawInfo)
	return info, err
}
